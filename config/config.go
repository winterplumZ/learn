package config

import (
	"github.com/Unknwon/goconfig"
	"github.com/coreos/pkg/capnslog"
	"reflect"
)

var (
	log  = capnslog.NewPackageLogger(reflect.TypeOf(struct{}{}).PkgPath(), "config.config")
	conf *Config
)

type ServerConfig struct {
	Host string
	Port string
}

/*
type KafkaConfig struct {
	Broker string
	Topic  string
}
*/

type DBconfig struct {
	User   string
	Pwd    string
	Ip     string
	Port   string
	DbName string
}

type Config struct {
	Server ServerConfig
	DB     DBconfig
	//	Kafka  KafkaConfig
}

func LoadConfig() *Config {
	if conf != nil {
		return conf
	}
	cfg, err := goconfig.LoadConfigFile("config.ini")
	if err != nil {
		log.Debug("load config file error", err.Error())
		return nil
	}

	udphost, _ := cfg.GetValue("server", "host")
	udpport, _ := cfg.GetValue("server", "port")

	//	kbroker, _ := cfg.GetValue("kafka", "broker")
	//	ktopic, _ := cfg.GetValue("kafka", "topic")

	dbuser, _ := cfg.GetValue("mysql", "user")
	dbpassword, _ := cfg.GetValue("mysql", "password")
	dbip, _ := cfg.GetValue("mysql", "ip")
	dbport, _ := cfg.GetValue("mysql", "port")
	dbname, _ := cfg.GetValue("mysql", "dbname")
	//	conf = &Config{Server: ServerConfig{Host: udphost, Port: udpport}, Kafka: KafkaConfig{Broker: kbroker, Topic: ktopic}}
	conf = &Config{
		Server: ServerConfig{Host: udphost, Port: udpport},
		DB:     DBconfig{User: dbuser, Pwd: dbpassword, Ip: dbip, Port: dbport, DbName: dbname},
	}
	return conf
}
