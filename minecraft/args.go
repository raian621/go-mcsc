package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/raian621/go-mcsc/api"
)

func NewServerArgs() *api.ServerArguments {
	return &api.ServerArguments{
		MemoryStartGB: ref(1),
		MemoryMaxGB:   ref(2),
	}
}

func BuildStringArgs(version string, args *api.ServerArguments) []string {
	strArgs := []string{"java"}

	strArgs = append(
		strArgs,
		fmt.Sprintf("-Xms%dG", *args.MemoryStartGB),
		fmt.Sprintf("-Xmx%dG", *args.MemoryMaxGB),
		"-jar",
		fmt.Sprintf("server-%s.jar", version),
		"--nogui",
	)

	if args.BonusChest != nil && *args.BonusChest {
		strArgs = append(strArgs, "--bonusChest")
	}
	if args.Demo != nil && *args.Demo {
		strArgs = append(strArgs, "--demo")
	}
	if args.EraseCache != nil && *args.EraseCache {
		strArgs = append(strArgs, "--eraseCache")
	}
	if args.ForceUpgrade != nil && *args.ForceUpgrade {
		strArgs = append(strArgs, "--forceUpgrade")
	}
	if args.SafeMode != nil && *args.SafeMode {
		strArgs = append(strArgs, "--safeMode")
	}
	if args.ServerID != nil && len(*args.ServerID) > 0 {
		strArgs = append(strArgs, "--serverId", *args.ServerID)
	}
	if args.SinglePlayer != nil && len(*args.SinglePlayer) > 0 {
		strArgs = append(strArgs, "--singleplayer", *args.SinglePlayer)
	}
	if args.Universe != nil && len(*args.Universe) > 0 {
		strArgs = append(strArgs, "--universe", *args.Universe)
	}
	if args.World != nil && len(*args.World) > 0 {
		strArgs = append(strArgs, "--world", *args.World)
	}
	if args.Port != nil && *args.Port > 0 {
		strArgs = append(strArgs, "--port", strconv.FormatInt(int64(*args.Port), 10))
	}

	return strArgs
}

func (m *JavaMinecraftServer) Args() *api.ServerArguments {
	m.Lock()
	defer m.Unlock()

	if m.args == nil {
		return nil
	}

	argsCpy := *m.args

	return &argsCpy
}

func (m *JavaMinecraftServer) CreateArgs() {
	m.Lock()
	defer m.Unlock()

	m.args = NewServerArgs()
}

func (m *JavaMinecraftServer) LoadArgs(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	if m.args == nil {
		return ErrNilConfig
	}

	return json.NewDecoder(file).Decode(m.args)
}

func (m *JavaMinecraftServer) SaveArgs(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.args == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.args)
}

func (m *JavaMinecraftServer) SetArgs(args *api.ServerArguments) {
	m.Lock()
	defer m.Unlock()

	m.args = args
}
