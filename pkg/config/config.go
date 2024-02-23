package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// configApp is the configuration for the application
type configApp struct {
	Token    string `yaml:"token"`
	Endpoint string `yaml:"endpoint"`
	homedir  string
}

// NewConfigApp creates a new configuration for the application
func NewConfigApp() *configApp {
	cfg := &configApp{
		Token:    "",
		Endpoint: "",
		homedir:  "",
	}
	cfg.initHomedir()
	return cfg
}

// initHomedir initializes the homedir variable
// It tries to get the user's home directory using os.UserHomeDir()
// If it fails, it tries to get the home directory from the HOME environment variable
func (c *configApp) initHomedir() {
	var err error
	c.homedir, err = os.UserHomeDir()
	if err != nil {
		c.homedir = os.Getenv("HOME")
	}
}

// SetHomedir sets the homedir variable
// It is used for testing purposes
func (c *configApp) SetHomedir(homedir string) {
	c.homedir = homedir
}

// LoadConfigFromFile loads the configuration from a file
func (c *configApp) LoadConfigFromFile(filename string) error {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		return fmt.Errorf("error parsing YAML file: %s", err)
	}
	return err
}

// LoadConfigFromEnvVar loads the configuration from the environment variables
func (c *configApp) LoadConfigFromEnvVar() {
	c.Token = os.Getenv("EPHEMERALFILES_TOKEN")
	c.Endpoint = os.Getenv("EPHEMERALFILES_ENDPOINT")
}

// IsConfigValid checks if the configuration is valid
func (c *configApp) IsConfigValid() bool {
	if c.Token == "" || c.Endpoint == "" {
		return false
	}
	return true
}

// LoadConfiguration loads the configuration from the environment variables first
// If the configuration is not valid, it tries to load the configuration from the $HOME/.eph.yml file
func (c *configApp) LoadConfiguration() error {
	c.LoadConfigFromEnvVar()
	if c.IsConfigValid() {
		return nil
	}
	err := c.LoadConfigFromFile(fmt.Sprintf("%s/.eph.yml", c.homedir))
	if err != nil {
		return err
	}
	if c.IsConfigValid() {
		return nil
	}
	return fmt.Errorf("configuration not found")
}
