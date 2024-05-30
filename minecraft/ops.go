package minecraft

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/raian621/go-mcsc/api"
)

var ErrNotInOps = errors.New("player not in server operator list")

// Deop implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Deop(p *api.PlayerInfo) error {
	m.Lock()
	defer m.Unlock()

	if m.ops == nil {
		return ErrNilConfig
	}

	idx := -1
	for i, op := range *m.ops {
		if op.Name == *p.Name && op.Uuid == *p.Uuid {
			idx = i
			break
		}
	}

	if idx == -1 {
		return ErrNotInOps
	}
	*m.ops = append((*m.ops)[:idx], (*m.ops)[idx+1:]...)

	return nil
}

// Op implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Op(op *api.ServerOperator) error {
	m.Lock()
	defer m.Unlock()

	if m.ops == nil {
		return ErrNilConfig
	}

	*m.ops = append(*m.ops, *op)

	return nil
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
	m.Lock()
	defer m.Unlock()

	m.ops = ref(make(api.ServerOperatorList, 0))
}

func (m *JavaMinecraftServer) LoadOperators(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	if m.ops == nil {
		return ErrNilConfig
	}

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
