package sensu

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	RabbitMQ RabbitMQConfig
}

type RabbitMQConfig struct {
	Host     string    `json:"host"`
	Port     int       `json:"port"`
	User     string    `json:"user"`
	Password string    `json:"password"`
	VHost    string    `json:"vhost"`
	SSL      SSLConfig `json:"ssl"`
}

type SSLConfig struct {
	CertChainFile  string `json:"cert_chain_file"`
	PrivateKeyFile string `json:"private_key_file"`
}

type MyError struct {
	Message string
}

func (e MyError) Error() string {
	return fmt.Sprintf("%s", e.Message)
}

const DEFAULT_SENSU_CONFIG_DIR string = "/etc/sensu"

const DEFAULT_SENSU_CONFIG_FILE string = "config.json"

func LoadConfig() (config *Config, err error) {
	config, err = LoadConfigDir(DEFAULT_SENSU_CONFIG_DIR)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func LoadConfigDir(config_dir string) (config *Config, err error) {

	var config_file string = config_dir + "/" + DEFAULT_SENSU_CONFIG_FILE

	// Ensure config file exists
	if _, err := os.Stat(config_file); os.IsNotExist(err) {
		return nil, MyError{"File does not exist: " + config_file}
	}

	// Load main config file
	data, err := ioutil.ReadFile(config_file)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	// Load additional config in conf.d
	files, _ := ioutil.ReadDir(config_dir + "/conf.d/")
	for _, f := range files {
		if !f.IsDir() && strings.HasSuffix(f.Name(), ".json") {
			data, err = ioutil.ReadFile(config_dir + "/conf.d/" + f.Name())
			if err != nil {
				return nil, err
			}

			err = json.Unmarshal(data, &config)
			if err != nil {
				return nil, err
			}
		}
	}

	return config, nil
}
