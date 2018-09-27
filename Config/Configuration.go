package Config

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"os"
	"strings"
	"github.com/op/go-logging"
)

type DbConfig struct {
	DbDriver  string
	Host      string
	Port      uint32
	TableName string
	UserName  string
	Password  string
}
type DhtConfig struct {
	StartPort int
	Workers   int
}
type HttpConfig struct {
	StaticDataFolder      string
	JwtAuthPrivateKeyPath string
	JwtAuthPublicKeyPath  string
}
type TagProducersConfig struct {
	TagProducers map[string][]string
}
type ScrapeConfig struct {
	Trackers      []string
	WorkerThreads int
	ScrapeTimeout int
}
type RpcMode string

const (
	SERVER RpcMode = "SERVER"
	CLIENT         = "CLIENT"
)

func (o *RpcMode) UnmarshalText(b []byte) (e error) {
	str := strings.Trim(string(b), `"`)

	switch str {
	case string(SERVER), string(CLIENT):
		*o = RpcMode(str)

	default:
		e = errors.New("Unknown RpcMode specified")
	}

	return e
}

type RpcConfig struct {
	Mode RpcMode
	Host string
	Port int
}
type TmdbApi struct {
	ApiKey string
}
type Configuration struct {
	DbConfig           DbConfig
	DhtConfig          DhtConfig
	HttpConfig         HttpConfig
	ScrapeConfig       ScrapeConfig
	RpcConfig          RpcConfig
	TagProducersConfig TagProducersConfig
	TmdbApi            TmdbApi
	ItemsPerPage       uint64
}

func SetupConfiguration() *Configuration {
	configFileName := flag.String("config", "config.json", "Path to config file")
	flag.Parse()
	file, err := os.Open(*configFileName)
	if err != nil {
		log.Panicln("Error opening Config:", err.Error())
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Panicln("Error reading Config:", err)
	}
	log.Println("Level values", logging.CRITICAL, logging.ERROR, logging.WARNING, logging.NOTICE, logging.INFO, logging.DEBUG)
	log.Println("Configuration ", configuration)
	return &configuration
}
