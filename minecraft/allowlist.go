package minecraft

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

var (
	ErrPlayerNotInAllowlist = errors.New("player is not in server allowlist")
)

// AllowPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) AllowPlayer(p api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	*m.allowlist = append(*m.allowlist, p)

	return nil
}

// Allowlist implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Allowlist() *api.Allowlist {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return nil
	}

	allowlistCpy := *m.allowlist

	return &allowlistCpy
}

// DisallowPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) DisallowPlayer(p api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	if len(*m.allowlist) == 0 {
		return ErrPlayerNotInAllowlist
	}

	idx := -1
	for i, player := range *m.allowlist {
		if *player.Name == *p.Name && *player.Uuid == *p.Uuid {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ErrPlayerNotInAllowlist
	}

	*m.allowlist = append((*m.allowlist)[:idx], (*m.allowlist)[idx+1:]...)

	return nil
}

func (m *JavaMinecraftServer) CreateAllowlist() {
	m.Lock()
	defer m.Unlock()

	m.allowlist = ref(make(api.Allowlist, 0))
}

func (m *JavaMinecraftServer) LoadAllowlist(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	return json.NewDecoder(file).Decode(m.allowlist)
}

func (m *JavaMinecraftServer) SaveAllowlist(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.allowlist)
}

func (m *JavaMinecraftServer) SetAllowlist(a *api.Allowlist) {
	m.Lock()
	defer m.Unlock()

	m.allowlist = a
}
