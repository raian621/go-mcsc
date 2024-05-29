package minecraft

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"

	"github.com/raian621/minecraft-server-controller/api"
)

var ErrNotInBannedIPs = errors.New("IP was not in ban list")

// BanIP implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) BanIP(ip *api.BannedIP) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedIPs == nil {
		return ErrNilConfig
	}

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/ban-ip %s", ip.Ip)); err != nil {
			return nil
		}
	}

	*m.bannedIPs = append(*m.bannedIPs, *ip)

	return nil
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
	m.Lock()
	defer m.Unlock()

	if m.bannedIPs == nil {
		return ErrNilConfig
	}

	idx := -1
	for i, b := range *m.bannedIPs {
		if b.Ip == ip {
			idx = i
		}
	}

	if idx == -1 {
		return ErrNotInBannedIPs
	}

	if m.console != nil {
		if err := m.console.SendCommand(fmt.Sprintf("/pardon-ip %s", ip)); err != nil {
			return err
		}
	} else {
		*m.bannedIPs = append((*m.bannedIPs)[:idx], (*m.bannedIPs)[:idx+1]...)
	}

	return nil
}

func (m *JavaMinecraftServer) CreateBannedIPs() {
	m.Lock()
	defer m.Unlock()

	m.bannedIPs = ref(make(api.BannedIPList, 0))
}

func (m *JavaMinecraftServer) LoadBannedIPs(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	if m.bannedIPs == nil {
		return ErrNilConfig
	}

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

	if m.console != nil {
		bannedIPs := make(map[api.BannedIP]bool, 0)

		for i := range *m.bannedIPs {
			bannedIPs[(*m.bannedIPs)[i]] = false
		}
		for i := range *b {
			bannedIPs[(*b)[i]] = true
		}

		for bannedIP, banned := range bannedIPs {
			if banned {
				if err := m.console.SendCommand(fmt.Sprintf("/ban-ip %s", bannedIP.Ip)); err != nil {
					log.Println(err)
					return
				}
			} else {
				if err := m.console.SendCommand(fmt.Sprintf("/pardon-ip %s", bannedIP.Ip)); err != nil {
					log.Println(err)
					return
				}
			}
		}
	}

	m.bannedIPs = b
}
