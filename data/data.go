package data

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

// CloudFoundriesType Store meta data for multiple cloud foundry instances
type CloudFoundriesType struct {
	CloudFoundries []CloudFoundryType
}

// CloudFoundryType Defines meta data of a single cloud foundry instance
type CloudFoundryType struct {
	APIURL   string   `yaml:"apiURL"`
	User     string   `yaml:"user"`
	Password string   `yaml:"password"`
	Labels   []string `yaml:"labels"`
}

var (
	mux                  sync.Mutex
	lastModifiedHash     uint32    = 0
	lastModified         time.Time = time.Time{}
	cachedCloudFoundries           = &CloudFoundriesType{}
	cachedConfigPath     string    = ""
)

// NewCloudFoundryMetaData Read cloud foundry data structure from YAML file.
// The data is cached and file is read only in case of content was modified. Date will
// be returned as a deep copy to avoid synchronization issues.
func NewCloudFoundryMetaData(configPath string) (*CloudFoundriesType, error) {

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

		if err := yaml.Unmarshal(dat, &cachedCloudFoundries); err != nil {
			log.Printf("Error while parsing YAML file %v: %v", configPath, err)
			return nil, err
		}
		log.Println(cachedCloudFoundries)
	} else {
		log.Printf("Using cached data of file %v", configPath)
	}

	copiedCloudFoundries, err := deepCopy(cachedCloudFoundries)

	return copiedCloudFoundries, err
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

func deepCopy(cf1 *CloudFoundriesType) (*CloudFoundriesType, error) {
	data, err := json.Marshal(cf1)
	if err != nil {
		return nil, err
	}

	cf2 := CloudFoundriesType{}
	if err := json.Unmarshal(data, &cf2); err != nil {
		return nil, err
	}
	return &cf2, nil
}

// GetLastModifiedHash returns a hash that can be used to build an ETag
func (*CloudFoundriesType) GetLastModifiedHash() uint32 {
	return lastModifiedHash
}

// GetLastModified returns last modified timestamp for setting Last-Modified header
func (*CloudFoundriesType) GetLastModified() time.Time {
	return lastModified
}
