package minecraft

import (
	"encoding/json"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

// BanIP implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) BanIP(ip api.BannedIP) error {
	panic("unimplemented")
}

// BannedIPs implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) BannedIPs() *api.BannedIPList {
	m.Lock()
	defer m.Unlock()

	if m.bannedIPs == nil {
		return nil
	}

	bannedIPsCpy := make(api.BannedIPList, len(*m.bannedIPs))
	copy(*m.bannedIPs, bannedIPsCpy)

	return &bannedIPsCpy
}

// PardonIP implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) PardonIP(ip string) error {
	panic("unimplemented")
}

func (m *JavaMinecraftServer) CreateBannedIPs() {
	m.bannedIPs = ref(make(api.BannedIPList, 0))
}

func (m *JavaMinecraftServer) LoadBannedIPs(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	return json.NewDecoder(file).Decode(m.bannedIPs)
}

func (m *JavaMinecraftServer) SaveBannedIPs(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedIPs == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.bannedIPs)
}

func (m *JavaMinecraftServer) SetBannedIPs(b *api.BannedIPList) {
	m.Lock()
	defer m.Unlock()

	m.bannedIPs = b
}
