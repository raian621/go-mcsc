package minecraft

import (
	"encoding/json"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

type ServerConfig struct {
	Version string `json:"version"`
}

func NewServerConfig() *api.MinecraftServerConfig {
	return &api.MinecraftServerConfig{
		Version: "1.20.6",
	}
}

func (m *JavaMinecraftServer) Config() *api.MinecraftServerConfig {
	m.Lock()
	defer m.Unlock()

	if m.config == nil {
		return nil
	}

	configCpy := *m.config

	return &configCpy
}

func (m *JavaMinecraftServer) CreateConfig() {
	m.config = NewServerConfig()
}

func (m *JavaMinecraftServer) LoadConfig(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	return json.NewDecoder(file).Decode(m.config)
}

func (m *JavaMinecraftServer) SaveConfig(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.config == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.config)
}
