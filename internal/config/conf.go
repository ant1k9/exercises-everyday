package config

import (
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// Config is a configuration from yaml
type Config struct {
	Database struct {
		Name     string `yaml:"name"`
		Username string `yaml:"username"`
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		Extra    string `yaml:"extra"`
	} `yaml:"database"`

	Server struct {
		Host string `yaml:"host"`
		Port string `yaml:"port"`
	} `yaml:"server"`
}

var (
	Conf Config
)

func init() {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &Conf)
	if err != nil {
		log.Fatal(err)
	}
	if port := os.Getenv("PORT"); port != "" {
		Conf.Server.Port = port
	}
	if database := os.Getenv("DATABASE_URL"); database != "" {
		parseDatabaseConfig(database)
	}
}

func parseDatabaseConfig(database string) {
	r := regexp.MustCompile(`postgres://(\w+):(\w+)@([\w-.]+):(\d+)/([\w_]+)`)
	parts := r.FindStringSubmatch(database)
	if len(parts) > 0 {
		Conf.Database.Username = parts[1]
		Conf.Database.Password = parts[2]
		Conf.Database.Host = parts[3]
		Conf.Database.Port = parts[4]
		Conf.Database.Name = parts[5]
	}
}
