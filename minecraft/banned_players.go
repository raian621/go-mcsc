package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/raian621/go-mcsc/api"
)

var ErrNotInBannedPlayers = errors.New("player not in banned players list")

// BanPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) BanPlayer(p *api.BannedPlayer) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedPlayers == nil {
		return ErrNilConfig
	}

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/ban %s", *p.Name)); err != nil {
			return err
		}
	}

	*m.bannedPlayers = append(*m.bannedPlayers, *p)

	return nil
}

// PardonPlayer implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) PardonPlayer(p *api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedPlayers == nil {
		return ErrNilConfig
	}

	idx := -1
	for i, b := range *m.bannedPlayers {
		if *b.Name == *p.Name && b.Uuid == *p.Uuid {
			idx = i
			break
		}
	}
	if idx == -1 {
		return ErrNotInBannedPlayers
	}

	*m.bannedPlayers = append((*m.bannedPlayers)[:idx], (*m.bannedPlayers)[idx+1:]...)

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/pardon %s", *p.Name)); err != nil {
			return err
		}
	}

	return nil
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

	if m.bannedPlayers == nil {
		return ErrNilConfig
	}

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
