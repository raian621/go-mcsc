package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/raian621/go-mcsc/api"
)

var (
	ErrPlayerNotInAllowlist = errors.New("player is not in server allowlist")
)

// Adds a player to the Minecraft server's allowlist.
//
// If the Minecraft server process is not currently running, the player is only
// added to the in-memory allowlist and will not be added to the Minecraft
// server's allowlist file until the JavaMinecraftServer's SaveAllowList method
// is called.
//
// If the Minecraft server process is currently running, the command `/whitelist
// <playername>` is sent to the Minecraft server console and the player will be
// added to the Minecraft server's allowlist and the in-memory allowlist upon
// success.
func (m *JavaMinecraftServer) AllowPlayer(p *api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/whitelist add %s", *p.Name)); err != nil {
			return err
		}
	}

	*m.allowlist = append(*m.allowlist, *p)

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

// Removes a player from the Minecraft server's allowlist.
//
// If the Minecraft server process is not currently running, the player is only
// removed from the in-memory allowlist and will not be removed from the
// Minecraft server's allowlist file until the JavaMinecraftServer's
// SaveAllowList method is called.
//
// If the Minecraft server process is currently running, the command `/whitelist
// remove <playername>` is sent to the Minecraft server console and the player will be
// added to the Minecraft server's allowlist and the in-memory allowlist upon
// success.func (m *JavaMinecraftServer) DisallowPlayer(p api.PlayerInfo) error {
func (m *JavaMinecraftServer) DisallowPlayer(p *api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.allowlist == nil {
		return ErrNilConfig
	}

	if len(*m.allowlist) == 0 {
		return ErrPlayerNotInAllowlist
	}

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/whitelist remove %s", *p.Name)); err != nil {
			return err
		}
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

	if m.console != nil {
		playerlist := make(map[string]bool, 0)

		for _, player := range *m.allowlist {
			playerlist[*player.Name] = false
		}
		for _, player := range *a {
			playerlist[*player.Name] = true
		}

		for name, allowed := range playerlist {
			if allowed {
				if err := m.console.SendCommand(fmt.Sprintf("/whitelist add %s", name)); err != nil {
					log.Println(err)
					return
				}
			} else {
				if err := m.console.SendCommand(fmt.Sprintf("/whitelist remove %s", name)); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	m.allowlist = a
}
