package config_test

import (
	"os"
	"testing"

	"github.com/ephemeralfiles/eph-beta/pkg/config"
	"github.com/stretchr/testify/assert"
)

// createTempConfigFile creates a temporary file with a valid configuration
func createValidConfigFile() (string, error) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("/tmp", "example")
	if err != nil {
		return "", err
	}
	// Write to the file
	_, err = tmpfile.Write([]byte("token: sdf\nendpoint: http://localhost:8080"))
	if err != nil {
		return "", err
	}
	// Close the file
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}
func TestReadyamlConfigFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := createValidConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile) // clean up

	// Test the function
	cfg := config.NewConfigApp()
	err = cfg.LoadConfigFromFile(tmpfile)
	assert.Nil(t, err)
	assert.Equal(t, "sdf", cfg.Token)
	assert.Equal(t, "http://localhost:8080", cfg.Endpoint)
}

func TestReadyamlConfigFileReturnErrorIfFileDoesNotExist(t *testing.T) {
	err := config.NewConfigApp().LoadConfigFromFile("/tmp/test.yml")
	assert.NotNil(t, err)
}

func TestGetConfigFromEnvVar(t *testing.T) {
	os.Setenv("EPHEMERALFILES_TOKEN", "sdf")
	os.Setenv("EPHEMERALFILES_ENDPOINT", "http://localhost:8080")

	// Test the function
	cfg := config.NewConfigApp()
	cfgExpected := config.NewConfigApp()
	cfgExpected.Endpoint = "http://localhost:8080"
	cfgExpected.Token = "sdf"
	cfg.LoadConfigFromEnvVar()
	assert.Equal(t, cfgExpected, cfg)
}

func TestConfigApp_IsConfigValid(t *testing.T) {
	type fields struct {
		Token    string
		Endpoint string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name:   "empty",
			fields: fields{},
			want:   false,
		},
		{
			name:   "empty token",
			fields: fields{Token: "", Endpoint: "http://localhost:8080"},
			want:   false,
		},
		{
			name:   "empty endpoint",
			fields: fields{Token: "sdf", Endpoint: ""},
			want:   false,
		},
		{
			name:   "valid",
			fields: fields{Token: "sdf", Endpoint: "http://localhost:8080"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := config.NewConfigApp()
			c.Token = tt.fields.Token
			c.Endpoint = tt.fields.Endpoint
			if got := c.IsConfigValid(); got != tt.want {
				t.Errorf("ConfigApp.IsConfigValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoadConfiguration(t *testing.T) {
	// Test that env var overrides file
	os.Setenv("EPHEMERALFILES_TOKEN", "envvar")
	os.Setenv("EPHEMERALFILES_ENDPOINT", "http://localhost:8080")

	// Test the function
	cfg := config.NewConfigApp()
	cfg.SetHomedir("/tmp")
	err := cfg.LoadConfiguration()
	assert.Nil(t, err)
	assert.Equal(t, "envvar", cfg.Token)
	assert.Equal(t, "http://localhost:8080", cfg.Endpoint)
}

func TestLoadConfigurationWithNoEnv(t *testing.T) {
	// Create file /tmp/.eph.yml
	tmpfile, err := os.Create("/tmp/.eph.yml")
	if err != nil {
		assert.Nil(t, err)
	}
	// Write to the file
	_, err = tmpfile.Write([]byte("token: sdf\nendpoint: http://localhost:8080"))
	if err != nil {
		assert.Nil(t, err)
	}
	// Close the file
	if err := tmpfile.Close(); err != nil {
		assert.Nil(t, err)
	}

	// Test the function
	cfg := config.NewConfigApp()
	cfg.SetHomedir("/tmp")
	err = cfg.LoadConfiguration()
	assert.Nil(t, err)
	assert.Equal(t, "envvar", cfg.Token)
	assert.Equal(t, "http://localhost:8080", cfg.Endpoint)
}
