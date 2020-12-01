package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/bennyharvey/soma/entity"
	"gopkg.in/yaml.v2"
)

type passageOpenerConfigRaw struct {
	Type          entity.PassageType `yaml:"type"`
	Address       string             `yaml:"address"`
	Direction     entity.Direction   `yaml:"direction"`
	WaitAfterOpen string             `yaml:"wait_after_open"`
}

type passageOpenerConfig struct {
	passageOpenerConfigRaw
	WaitAfterOpen time.Duration
}

func (sc *passageOpenerConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cRaw passageOpenerConfigRaw

	err := unmarshal(&cRaw)
	if err != nil {
		return fmt.Errorf("YAML unmarshal: %w", err)
	}

	sc.passageOpenerConfigRaw = cRaw

	sc.WaitAfterOpen, err = time.ParseDuration(cRaw.WaitAfterOpen)
	if err != nil {
		return fmt.Errorf("wait_after_open parse: %w", err)
	}

	return nil
}

func (c passageOpenerConfig) Validate() error {
	switch c.Type {
	case entity.Z5R, entity.Sigur, entity.Dummy, entity.Beward:
	default:
		return errors.New("type is unknown")
	}
	if c.Address == "" {
		return errors.New("address is empty")
	}
	if !(c.Direction == entity.In || c.Direction == entity.Out) {
		return errors.New("direction is invalid")
	}
	if c.WaitAfterOpen < 0 {
		return errors.New("wait_after_open is invalid")
	}
	return nil
}

type webServerConfig struct {
	BindAddr       string            `yaml:"bind_addr"`
	JWTSigningKey  string            `yaml:"jwt_signing_key"`
	TLSCrtFilePath string            `yaml:"tls_crt_file_path"`
	TLSKeyFilePath string            `yaml:"tls_key_file_path"`
	Debug          bool              `yaml:"debug"`
	PassageNames   map[string]string `yaml:"passage_names"`
}

func (c webServerConfig) Validate() error {
	if c.BindAddr == "" {
		return errors.New("bind_addr is empty")
	}
	if c.JWTSigningKey == "" {
		return errors.New("jwt_signing_key is empty")
	}
	if c.TLSCrtFilePath == "" {
		return errors.New("tls_crt_file_path is empty")
	}
	if c.TLSKeyFilePath == "" {
		return errors.New("tls_key_file_path is empty")
	}
	return nil
}

type config struct {
	CUDAVisibleDevices       string                         `yaml:"cuda_visible_devices"`
	DetectorModelPath        string                         `yaml:"detector_model_path"`
	ShaperModelPath          string                         `yaml:"shaper_model_path"`
	RecognizerModelPath      string                         `yaml:"recognizer_model_path"`
	Jittering                int                            `yaml:"jittering"`
	RabbitMQURI              string                         `yaml:"rabbitmq_uri"`
	RabbitMQExchange         string                         `yaml:"rabbitmq_exchange"`
	PostgresURI              string                         `yaml:"postgres_uri"`
	DetectConfidenceLimit    float64                        `yaml:"detect_confidence_limit"`
	DescriptorsMatchDistance float64                        `yaml:"descriptors_match_distance"`
	PassageOpeners           map[string]passageOpenerConfig `yaml:"passage_openers"`
	PhotoStoragePath         string                         `yaml:"photo_storage_path"`
	WebServer                webServerConfig                `yaml:"web_server"`
}

func (c config) Validate() error {
	if c.DetectorModelPath == "" {
		return errors.New("detector_model_path is empty")
	}
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
	if c.PostgresURI == "" {
		return errors.New("postgres_uri is empty")
	}
	if c.DetectConfidenceLimit < 0 {
		return errors.New("detect_confidence_limit is invalid")
	}
	if c.DescriptorsMatchDistance < 0 {
		return errors.New("descriptors_match_distance is invalid")
	}
	for passageID, po := range c.PassageOpeners {
		err := po.Validate()
		if err != nil {
			return fmt.Errorf("passage_opener %s: %w", passageID, err)
		}
	}
	err := c.WebServer.Validate()
	if err != nil {
		return fmt.Errorf("web_server: %w", err)
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
