package minecraft

import (
	"encoding/json"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

// Deop implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Deop(p api.PlayerInfo) error {
	panic("unimplemented")
}

// Op implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Op(op api.ServerOperator) error {
	panic("unimplemented")
}

// Ops implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Ops() *api.ServerOperatorList {
	m.Lock()
	defer m.Unlock()

	if m.ops == nil {
		return nil
	}

	serverOperatorsCpy := make(api.ServerOperatorList, len(*m.ops))
	copy(*m.ops, serverOperatorsCpy)

	return &serverOperatorsCpy
}

func (m *JavaMinecraftServer) CreateOperators() {
	m.ops = ref(make(api.ServerOperatorList, 0))
}

func (m *JavaMinecraftServer) LoadOperators(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	return json.NewDecoder(file).Decode(m.ops)
}

func (m *JavaMinecraftServer) SaveOperators(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.ops == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.ops)
}

func (m *JavaMinecraftServer) SetOperators(ops *api.ServerOperatorList) {
	m.Lock()
	defer m.Unlock()

	m.ops = ops
}
