package config

import (
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// AuthTypeBasic basic authentification
	AuthTypeBasic string = "basic"
)

// Configuration struct for server configuration
type Configuration struct {
	Server struct {
		AuthType  string `yaml:"authtype"`
		BasicAuth struct {
			UserName string `yaml:"username"`
			Password string `yaml:"password"`
		} `yaml:"basicauth"`
	} `yaml:"server"`
	CloudFoundries map[string]struct {
		APIURL   string   `yaml:"apiURL"`
		UAAURL   string   `yaml:"uaaURL"`
		UserName string   `yaml:"username"`
		Password string   `yaml:"password"`
		Labels   []string `yaml:"labels"`
	} `yaml:"cloudfoundries"`
}

var (
	cfg                        = &Configuration{}
	lastModifiedHash uint32    = 0
	lastModified     time.Time = time.Time{}
)

// Read cloud foundry data structure from YAML file.
func Read(configPath string) error {

	log.Printf("Reading file %v", configPath)
	dat, err := ioutil.ReadFile(configPath)

	if err != nil {
		log.Printf("ERROR while reading config file: %v", err)
		return err
	}

	file, err := os.Stat(configPath)
	if err != nil {
		log.Printf("Error while reading last modified date: %v", err)
		return err
	}

	lastModified = file.ModTime()
	lastModifiedHash = hash(file.ModTime().String())

	if err := yaml.Unmarshal(dat, &cfg); err != nil {
		log.Printf("Error while parsing YAML file %v: %v", configPath, err)
		return err
	}
	log.Println(cfg)

	return nil
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

// GetLastModifiedHash returns a hash that can be used to build an ETag
func GetLastModifiedHash() uint32 {
	return lastModifiedHash
}

// GetLastModified returns last modified timestamp for setting Last-Modified header
func GetLastModified() time.Time {
	return lastModified
}

// Get returns configuration object
func Get() Configuration {
	return *cfg
}
