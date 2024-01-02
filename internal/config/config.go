package config

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/viper"
)

// Config represents the application configuration.
type Config struct {
	Name                 string     `mapstructure:"name"`
	Prefix               string     `mapstructure:"prefix"`
	Debug                bool       `mapstructure:"debug"`
	MqttConnection       string     `mapstructure:"mqtt_connection"`
	MqttSecureConnection string     `mapstructure:"mqtt_secure_connection"`
	HttpServer           HttpServer `mapstructure:"http_server"`
	Listeners            []string   `mapstructure:"listeners"`
}

type HttpServer struct {
	Enabled bool   `mapstructure:"enabled"`
	Port    string `mapstructure:"port"`
	Host    string `mapstructure:"host"`
}

// LoadConfig loads the configuration from the predefined paths.
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config") // no need to include file extension
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")                      // current working directory
	viper.AddConfigPath("$HOME/.config/tellall/") // home directory
	viper.AddConfigPath("/var/lib/tellall/")      // other directory

	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("error getting hostname: %w", err)
	}

	viper.SetDefault("name", hostname)
	viper.SetDefault("prefix", "tellall")
	viper.SetDefault("debug", false)

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	if c.MqttSecureConnection != "" {
		content, err := ioutil.ReadFile(c.MqttSecureConnection)
		if err != nil {
			return nil, fmt.Errorf("unable to read mqtt secure connection file: %w", err)
		}
		c.MqttConnection = string(content)
	}

	if c.MqttConnection == "" {
		return nil, fmt.Errorf("no mqtt connection defined")
	}

	// trim c.MqttConnection of whitespace
	c.MqttConnection = string(bytes.TrimSpace([]byte(c.MqttConnection)))

	return &c, nil
}
