package minecraft

import (
	"encoding/json"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

// BanPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) BanPlayer(p *api.BannedPlayer) error {
	panic("unimplemented")
}

// UpdateBannedPlayers implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) UpdateBannedPlayers(bp *api.BannedPlayerList) error {
	panic("unimplemented")
}

// PardonPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) PardonPlayer(p *api.PlayerInfo) error {
	panic("unimplemented")
}

func (m *JavaMinecraftServer) BannedPlayers() *api.BannedPlayerList {
	m.Lock()
	defer m.Unlock()

	if m.bannedPlayers == nil {
		return nil
	}

	bannedPlayersCpy := make(api.BannedPlayerList, len(*m.bannedPlayers))
	copy(*m.bannedPlayers, bannedPlayersCpy)

	return &bannedPlayersCpy
}

func (m *JavaMinecraftServer) CreateBannedPlayers() {
	m.bannedPlayers = ref(make(api.BannedPlayerList, 0))
}

func (m *JavaMinecraftServer) LoadBannedPlayers(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	return json.NewDecoder(file).Decode(m.bannedPlayers)
}

func (m *JavaMinecraftServer) SaveBannedPlayers(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedPlayers == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.bannedPlayers)
}

func (m *JavaMinecraftServer) SetBannedPlayers(b *api.BannedPlayerList) {
	m.Lock()
	defer m.Unlock()

	m.bannedPlayers = b
}
