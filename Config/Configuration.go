package Config

import (
	"encoding/json"
	"flag"
	"log"
	"os"
)

type DbConfig struct {
	DbDriver  string
	TableName string
	UserName  string
	Password  string
}
type DhtConfig struct {
	StartPort int
	Workers   int
}
type HttpConfig struct {
	StaticDataFolder string
}

type ScrapeConfig struct {
	Trackers      []string
	WorkerThreads int
	ScrapeTimeout int
}

type Configuration struct {
	DbConfig     DbConfig
	DhtConfig    DhtConfig
	HttpConfig   HttpConfig
	ScrapeConfig ScrapeConfig
	ItemsPerPage int
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

	return &configuration
}
