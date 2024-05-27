package minecraft

import (
	"encoding/json"
	"errors"
	"io"
	"slices"
)

var ErrVersionUnsupported = errors.New("unsupported server version passed")

type VersionInfo struct {
	Link string `json:"link"`
	Sum  string `json:"sum"` // md5 sum
}

type VersionMap map[string]VersionInfo

func (v VersionMap) Versions() []string {
	versions := make([]string, 0, len(v))
	for version := range v {
		versions = append(versions, version)
	}

	slices.Sort(versions)
	return versions
}

func (v *VersionMap) Load(file io.Reader) error {
	return json.NewDecoder(file).Decode(v)
}

// SetVersion implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) SetVersion(version string) error {
	if m.versions == nil {
		return ErrNilConfig
	}

	supported := false
	for _, supportedVersion := range *m.versions {
		if version == supportedVersion {
			supported = true
			break
		}
	}

	if !supported {
		return ErrVersionUnsupported
	}

	m.config.Version = version

	return nil
}

// Versions implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Versions() *[]string {
	m.Lock()
	defer m.Unlock()

	if m.versions == nil {
		return nil
	}

	versionsCpy := make([]string, len(*m.versions))
	copy(versionsCpy, *m.versions)

	return &versionsCpy
}

func (m *JavaMinecraftServer) CreateVersions() {
	m.Lock()
	defer m.Unlock()

	m.versions = ref(make([]string, 0))
}

func (m *JavaMinecraftServer) LoadVersions(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	m.versions = new([]string)

	versionsMap := make(VersionMap, 0)
	if err := versionsMap.Load(file); err != nil {
		return err
	}

	// TODO: do this in place
	*m.versions = versionsMap.Versions()

	return nil
}
