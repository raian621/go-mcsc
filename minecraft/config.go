package minecraft

import (
	"sync"
)

type ServerConfig struct {
	ArgsFilepath       string `json:"argsFilepath"`
	PropertiesFilepath string `json:"propertiesFilepath"`
	OpsFilepath        string `json:"opsFilepath"`
	Version            string `json:"version"`
}

func NewServerConfig() *ServerConfig {
	return &ServerConfig{
		ArgsFilepath:       "server-data/args.json",
		PropertiesFilepath: "server-data/properties.json",
		OpsFilepath:        "server-data/ops.json",
		Version:            "1.20.6",
	}
}

type ServerConfigMonitor struct {
	config *ServerConfig
	mutex  sync.RWMutex
}

func NewServerConfigMonitor() *ServerConfigMonitor {
	return &ServerConfigMonitor{
		config: NewServerConfig(),
	}
}

func (scm *ServerConfigMonitor) Lock() {
	scm.mutex.Lock()
}

func (scm *ServerConfigMonitor) Unlock() {
	scm.mutex.Unlock()
}

func (scm *ServerConfigMonitor) RLock() {
	scm.mutex.RLock()
}

func (scm *ServerConfigMonitor) RUnlock() {
	scm.mutex.RUnlock()
}

func (scm *ServerConfigMonitor) LockAndGetData() any {
	scm.Lock()
	return scm.config
}

func (scm *ServerConfigMonitor) GetData() any {
	return scm.config
}

func (scm *ServerConfigMonitor) Load(filepath string) error {
	scm.Lock()
	defer scm.Unlock()
	return loadServerConfigObject(scm, func(obj ServerConfigObject) {
		configMonitor, ok := obj.(*ServerConfigMonitor)
		if ok {
			configMonitor.config = NewServerConfig()
		}
	}, filepath)
}

func (scm *ServerConfigMonitor) Save(filepath string) error {
	scm.Lock()
	defer scm.Unlock()
	return saveServerConfigObject(scm, filepath)
}
