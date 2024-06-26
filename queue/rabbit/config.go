package rabbit

import (
	"errors"
	"fmt"
	"github.com/dmitriibb/go-common/utils"
	"gopkg.in/yaml.v2"
	"os"
)

const (
	defaultFileNane string = "rabbitConfig.yaml"
)

var rabbitConfigFileNane = utils.GetEnvProperty(RabbitMqConfigFileEnv, defaultFileNane)

var config RabbitConfig

type RabbitQueueConfig struct {
	Name       string `yaml:"name"`
	Durable    bool   `yaml:"durable"`
	AutoDelete bool   `yaml:"autoDelete"`
	Persistent bool   `yaml:"persistent"`
}

type RabbitConfig struct {
	Uri    string              `yaml:"uri"`
	Queues []RabbitQueueConfig `yaml:"queues"`
}

func (conf RabbitConfig) GetQueueConfig(name string) (RabbitQueueConfig, error) {
	var res RabbitQueueConfig
	var err error
	found := false
	for _, qc := range conf.Queues {
		if qc.Name == name {
			res = qc
			found = true
			break
		}
	}

	if !found {
		err = errors.New(fmt.Sprintf("topic '%v' not found in the config", name))
	}
	return res, err
}

func GetUri() string {
	conf, _ := GetRabbitConfig()
	return conf.Uri
}

func GetQueueConfig(name string) (RabbitQueueConfig, error) {
	conf, _ := GetRabbitConfig()
	return conf.GetQueueConfig(name)
}

func GetRabbitConfig() (RabbitConfig, error) {
	file, err := os.ReadFile(rabbitConfigFileNane)
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(file, &config)
	return config, err
}
