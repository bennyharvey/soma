package main

import (
	"errors"
	"fmt"
	"image"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

type streamConfig struct {
	streamConfigRaw
	ClosedDuration time.Duration
}

type streamConfigRaw struct {
	URI            string `yaml:"uri"`
	PassageID      string `yaml:"passage_id"`
	ClosedDuration string `yaml:"closed_duration"`
	FrameRate      int    `yaml:"frame_rate"`
}

func (sc *streamConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cRaw streamConfigRaw

	err := unmarshal(&cRaw)
	if err != nil {
		return fmt.Errorf("YAML unmarshal: %w", err)
	}

	sc.streamConfigRaw = cRaw

	sc.ClosedDuration, err = time.ParseDuration(cRaw.ClosedDuration)
	if err != nil {
		return fmt.Errorf("closed_duration parse: %w", err)
	}

	return nil
}

func (sc streamConfig) Validate() error {
	if sc.URI == "" {
		return errors.New("uri is empty")
	}
	if sc.PassageID == "" {
		return errors.New("passage_id is empty")
	}
	if sc.ClosedDuration < time.Second {
		return errors.New("closed_duration is too small")
	}
	return nil
}

type configRaw struct {
	CUDAVisibleDevices   string         `yaml:"cuda_visible_devices"`
	DetectorModelPath    string         `yaml:"detector_model_path"`
	RabbitMQURI          string         `yaml:"rabbitmq_uri"`
	RabbitMQExchange     string         `yaml:"rabbitmq_exchange"`
	Streams              []streamConfig `yaml:"streams"`
	StreamsResolution    string         `yaml:"streams_resolution"`
	FaceDetectorWaitTime string         `yaml:"face_detector_wait_time"`
}

type config struct {
	configRaw
	StreamsResolution    image.Point
	FaceDetectorWaitTime time.Duration
}

func (c *config) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var cRaw configRaw

	err := unmarshal(&cRaw)
	if err != nil {
		return fmt.Errorf("YAML unmarshal: %w", err)
	}

	c.configRaw = cRaw

	streamsResolutionStrs := strings.Split(cRaw.StreamsResolution, "x")
	if len(streamsResolutionStrs) != 2 {
		return fmt.Errorf("streams_resolution parse: %w", err)
	}

	c.StreamsResolution.X, err = strconv.Atoi(streamsResolutionStrs[0])
	if err != nil {
		return fmt.Errorf("streams_resolutions X parse: %w", err)
	}

	c.StreamsResolution.Y, err = strconv.Atoi(streamsResolutionStrs[1])
	if err != nil {
		return fmt.Errorf("streams_resolutions Y parse: %w", err)
	}

	c.FaceDetectorWaitTime, err = time.ParseDuration(cRaw.FaceDetectorWaitTime)
	if err != nil {
		return fmt.Errorf("closed_duration parse: %w", err)
	}

	return nil
}

func (c config) Validate() error {
	if c.DetectorModelPath == "" {
		return errors.New("detector_model_path is empty")
	}
	if c.RabbitMQURI == "" {
		return errors.New("rabbitmq_uri is empty")
	}
	if c.RabbitMQExchange == "" {
		return errors.New("rabbitmq_exchange is empty")
	}
	if len(c.Streams) == 0 {
		return errors.New("streams is empty")
	}
	for i, s := range c.Streams {
		err := s.Validate()
		if err != nil {
			return fmt.Errorf("stream #%d: %v", i, err)
		}
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
