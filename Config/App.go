package Config

import (
	"errors"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/jasonlvhit/gocron"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"
	"github.com/op/go-logging"
	"github.com/ruslanfedoseenko/dhtcrawler/Models"
)

var clog = logging.MustGetLogger("Config")

type Service interface {
	Start()
}

type App struct {
	Config    *Configuration
	Db        *gorm.DB
	Scheduler *gocron.Scheduler
	Services  []Service
}

func NewApp() *App {
	raven.SetDSN("http://b80ba7ce05cb42e7983d962d95a1a6e5:37cd46ed05b54ce69fca375afd606bff@sentry.btoogle.com/4")
	setupLog()
	clog.Info("Test output")
	config := SetupConfiguration()
	app := App{

		Config:    config,
		Scheduler: gocron.NewScheduler(),
	}
	if config.DbConfig.DbDriver != "" {
		var connectionString string = fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable password=%s", config.DbConfig.Host, config.DbConfig.Port, config.DbConfig.UserName, config.DbConfig.TableName, config.DbConfig.Password)
		log.Println("Using Connection string", connectionString)
		db, err := gorm.Open(config.DbConfig.DbDriver, connectionString)
		if err != nil {
			log.Println("Db open error:", err.Error())
		}
		err = db.DB().Ping()
		if err != nil {
			log.Println("Db open error:", err.Error())
		}
		applySentryPlugin(db)
		db.AutoMigrate(&Models.JwtRefreshToken{}, &Models.User{})
		app.Db = db
	}

	return &app
}

var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{module} %{shortfunc} - %{level:.4s} %{color:reset}: %{message}`,
)

type SentryLogBackend struct {
}

func (b SentryLogBackend) Log(logLvl logging.Level, i int, record *logging.Record) error {


	switch logLvl {
		case logging.ERROR:
		case logging.CRITICAL:
		case logging.NOTICE:
		{
			e := errors.New(record.Formatted(i))

			evntId := raven.CaptureError(e, nil)
			log.Println("captureing error with raven:", e, "event id:", evntId)
			break
		}
	}
	return nil
}


func setupLog() {
	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	leveledSentryBackend := logging.AddModuleLevel(SentryLogBackend{})
	leveledSentryBackend.SetLevel(logging.NOTICE, "")
	leveledStdOutBackend := logging.AddModuleLevel(stdOutBackend)
	leveledStdOutBackend.SetLevel(logging.DEBUG, "")
	logging.SetFormatter(format)
	logging.SetBackend(leveledStdOutBackend, leveledSentryBackend)
}

func reportErrToSentryio(scope *gorm.Scope) {
	if scope.HasError() {
		err := scope.DB().Error
		sql := scope.CombinedConditionSql()
		sentryAdditionalData := map[string]string{
			"sql": sql,
		}
		raven.CaptureError(err, sentryAdditionalData)
	}
}
func applySentryPlugin(db *gorm.DB) {
	db.Callback().Create().After("gorm:create").Register("sentryio_plugin:after_create", reportErrToSentryio)
	db.Callback().Query().After("gorm:query").Register("sentryio_plugin:after_query", reportErrToSentryio)
	db.Callback().Delete().After("gorm:delete").Register("sentryio_plugin:after_delete", reportErrToSentryio)
	db.Callback().Update().After("gorm:update").Register("sentryio_plugin:after_update", reportErrToSentryio)
}

func (app *App) AddService(svc Service) {
	app.Services = append(app.Services, svc)
}

func (app *App) Run() {

	for _, service := range app.Services {
		service.Start()
	}

	go func() {
		log.Println("Starting Scheduler")
		<-app.Scheduler.Start()
	}()
	if app.Db != nil {
		defer app.Db.Close()
	}

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, os.Interrupt)
	signal.Notify(sigs, syscall.SIGTERM)
	signal.Notify(sigs, os.Kill)
	log.Println("pushed global error handler")
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in f", r)
			debug.PrintStack()
			// find out exactly what the error was and set err
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
			if err != nil {
				// sendMeMail(err)
				log.Println("Panic occured", err.Error())
			}
		}
	}()
	go func() {
		log.Println("Waiting for SIGTERM to exit cleanly")
		sig := <-sigs
		log.Println("recived signal", sig)
		done <- true
	}()

	<-done

}
