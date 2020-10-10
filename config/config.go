package config

import (
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// ConfigurationType Configuration data
type ConfigurationType struct {
	Broker         BrokerConfigType `yaml:"broker"`
	CloudFoundries map[string]CloudFoundryConfigType
}

// BrokerConfigType Store broker config data
type BrokerConfigType struct {
	BasicAuth BasicAuthType `yaml:"basicauth"`
}

// BasicAuthType Store basic auth data
type BasicAuthType struct {
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
}

// CloudFoundryConfigType Defines meta data of a single cloud foundry instance
type CloudFoundryConfigType struct {
	APIURL   string   `yaml:"apiURL"`
	UAAURL   string   `yaml:"uaaURL"`
	UserName string   `yaml:"username"`
	Password string   `yaml:"password"`
	Labels   []string `yaml:"labels"`
}

var (
	mux              sync.Mutex
	lastModifiedHash uint32    = 0
	lastModified     time.Time = time.Time{}
	cachedConfig               = &ConfigurationType{}
	cachedConfigPath string    = ""
)

// NewConfig Read cloud foundry data structure from YAML file.
// The data is cached and file is read only in case of content was modified. Date will
// be returned as a deep copy to avoid synchronization issues.
func NewConfig(configPath string) (*ConfigurationType, error) {

	mux.Lock()
	defer mux.Unlock()

	if hasChanged(configPath) {
		log.Printf("Reading file and update cache %v", configPath)
		dat, err := ioutil.ReadFile(configPath)

		if err != nil {
			log.Printf("ERROR while reading config file: %v", err)
			return nil, err
		}

		file, err := os.Stat(configPath)
		if err != nil {
			log.Printf("Error while reading last modified date: %v", err)
			lastModifiedHash = 0
		}

		lastModified = file.ModTime()
		lastModifiedHash = hash(file.ModTime().String())

		if err := yaml.Unmarshal(dat, &cachedConfig); err != nil {
			log.Printf("Error while parsing YAML file %v: %v", configPath, err)
			return nil, err
		}
		log.Println(cachedConfig)
	} else {
		log.Printf("Using cached data of file %v", configPath)
	}

	copiedConfig, err := deepCopy(cachedConfig)

	return copiedConfig, err
}

func hasChanged(configPath string) bool {
	if configPath != cachedConfigPath {
		cachedConfigPath = configPath
		return true
	}

	file, err := os.Stat(configPath)
	if err != nil {
		log.Printf("Error while reading last modified date: %v", err)
		lastModifiedHash = 0
		return true
	}
	currentHash := hash(file.ModTime().String())

	return currentHash != lastModifiedHash
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func deepCopy(cfg1 *ConfigurationType) (*ConfigurationType, error) {
	data, err := json.Marshal(cfg1)
	if err != nil {
		return nil, err
	}

	cfg2 := ConfigurationType{}
	if err := json.Unmarshal(data, &cfg2); err != nil {
		return nil, err
	}
	return &cfg2, nil
}

// GetLastModifiedHash returns a hash that can be used to build an ETag
func (*ConfigurationType) GetLastModifiedHash() uint32 {
	return lastModifiedHash
}

// GetLastModified returns last modified timestamp for setting Last-Modified header
func (*ConfigurationType) GetLastModified() time.Time {
	return lastModified
}
