package main

import (
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type config struct {
	CUDAVisibleDevices  string  `yaml:"cuda_visible_devices"`
	ShaperModelPath     string  `yaml:"shaper_model_path"`
	RecognizerModelPath string  `yaml:"recognizer_model_path"`
	Jittering           int     `yaml:"jittering"`
	RabbitMQURI         string  `yaml:"rabbitmq_uri"`
	RabbitMQExchange    string  `yaml:"rabbitmq_exchange"`
	PassageID           string  `yaml:"passage_id"`
	ConfidenceLimit     float64 `yaml:"confidence_limit"`
}

func (c config) Validate() error {
	if c.ShaperModelPath == "" {
		return errors.New("shaper_model_path is empty")
	}
	if c.RecognizerModelPath == "" {
		return errors.New("recognizer_model_path is empty")
	}
	if c.RabbitMQURI == "" {
		return errors.New("rabbitmq_uri is empty")
	}
	if c.RabbitMQExchange == "" {
		return errors.New("rabbitmq_exchange is empty")
	}
	if c.PassageID == "" {
		return errors.New("passage_id is empty")
	}
	if c.ConfidenceLimit < 0 {
		return errors.New("confidence_limit is invalid")
	}
	return nil
}

func loadConfig(configPath string) (config, error) {
	configYAML, err := ioutil.ReadFile(configPath)
	if err != nil {
		return config{}, fmt.Errorf("read config %s file: %w", configPath, err)
	}

	var c config

	err = yaml.Unmarshal(configYAML, &c)
	if err != nil {
		return config{}, fmt.Errorf("YAML unmarshal config: %w", err)
	}

	return c, nil
}
