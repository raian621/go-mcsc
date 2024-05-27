package minecraft

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/raian621/minecraft-server-controller/api"
)

// https://minecraft.fandom.com/wiki/Tutorials/Setting_up_a_server
type ServerArgs struct {
	BonusChest    bool   `json:"bonusChest"`
	Demo          bool   `json:"demo"`
	EraseCache    bool   `json:"eraseCache"`
	ForceUpgrade  bool   `json:"forceUpgrade"`
	SafeMode      bool   `json:"safeMode"`
	ServerID      string `json:"serverID"`
	SinglePlayer  string `json:"singlePlayer"`
	Universe      string `json:"universe"`
	World         string `json:"world"`
	Port          uint32 `json:"port"`
	MemoryStartGB uint32 `json:"memoryStartGB"`
	MemoryMaxGB   uint32 `json:"memoryMaxGB"`
}

func NewServerArgs() *api.ServerArguments {
	return &api.ServerArguments{
		BonusChest:    ref(false),
		Demo:          ref(false),
		EraseCache:    ref(false),
		ForceUpgrade:  ref(false),
		SafeMode:      ref(false),
		ServerID:      ref(""),
		SinglePlayer:  ref(""),
		Universe:      ref(""),
		World:         ref(""),
		Port:          ref(0),
		MemoryStartGB: ref(1),
		MemoryMaxGB:   ref(2),
	}
}

func BuildStringArgs(version string, args *api.ServerArguments) []string {
	strArgs := []string{"java"}

	strArgs = append(
		strArgs,
		fmt.Sprintf("-Xms%dG", args.MemoryStartGB),
		fmt.Sprintf("-Xmx%dG", args.MemoryMaxGB),
		"-jar",
		fmt.Sprintf("server-%s.jar", version),
		"--nogui",
	)

	if *args.BonusChest {
		strArgs = append(strArgs, "--bonusChest")
	}
	if *args.Demo {
		strArgs = append(strArgs, "--demo")
	}
	if *args.EraseCache {
		strArgs = append(strArgs, "--eraseCache")
	}
	if *args.ForceUpgrade {
		strArgs = append(strArgs, "--forceUpgrade")
	}
	if *args.SafeMode {
		strArgs = append(strArgs, "--safeMode")
	}
	if len(*args.ServerID) > 0 {
		strArgs = append(strArgs, "--serverId", *args.ServerID)
	}
	if len(*args.SinglePlayer) > 0 {
		strArgs = append(strArgs, "--singleplayer", *args.SinglePlayer)
	}
	if len(*args.Universe) > 0 {
		strArgs = append(strArgs, "--universe", *args.Universe)
	}
	if len(*args.World) > 0 {
		strArgs = append(strArgs, "--world", *args.World)
	}
	if *args.Port > 0 {
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
