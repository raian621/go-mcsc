package minecraft

import (
	"bytes"
	"testing"

	"github.com/raian621/go-mcsc/api"
)

func TestConfig(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	if server.Config() != nil {
		t.Errorf("expected nil config")
	}

	server.CreateConfig()
	if server.Config() == nil {
		t.Errorf("expected non nil config")
	}
}

func TestLoadConfig(t *testing.T) {
	t.Parallel()

	mockedJSON := `{"version":"1.20.6"}` + "\n"
	server := &JavaMinecraftServer{}

	if err := server.LoadConfig(bytes.NewBufferString(mockedJSON)); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}
	server.CreateConfig()
	if err := server.LoadConfig(bytes.NewBufferString(mockedJSON)); err != nil {
		t.Errorf("expected nil error, got `%v`", err)
	}
	if server.config.Version != "1.20.6" {
		t.Errorf("expected version `1.20.6`, got `%s`", server.config.Version)
	}
}

func TestSaveConfig(t *testing.T) {
	t.Parallel()

	server := &JavaMinecraftServer{}

	if err := server.SaveConfig(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.SetConfig(&api.MinecraftServerConfig{
		Version: "1.20.6",
	})

	var out bytes.Buffer
	if err := server.SaveConfig(&out); err != nil {
		t.Errorf("expected nil error, got `%v`", err)
	}

	expected := `{"version":"1.20.6"}` + "\n"
	if expected != out.String() {
		t.Errorf("expected version `%s`, got `%s`", expected, out.String())
	}
}

func TestSetConfig(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	config := NewServerConfig()
	server.SetConfig(config)
	serverConfig := server.Config()

	if config == serverConfig {
		t.Error("expected original config and returned config to have different memory addresses")
	}
}
