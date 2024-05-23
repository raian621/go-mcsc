package minecraft

import (
	"fmt"
	"strconv"
	"sync"
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

func NewServerArgs() *ServerArgs {
	return &ServerArgs{
		BonusChest:    false,
		Demo:          false,
		EraseCache:    false,
		ForceUpgrade:  false,
		SafeMode:      false,
		ServerID:      "",
		SinglePlayer:  "",
		Universe:      "",
		World:         "",
		Port:          0,
		MemoryStartGB: 1,
		MemoryMaxGB:   2,
	}
}

type ServerArgsMonitor struct {
	args  *ServerArgs
	mutex sync.RWMutex
}

func NewServerArgsMonitor() *ServerArgsMonitor {
	return &ServerArgsMonitor{
		args: NewServerArgs(),
	}
}

func (sam *ServerArgsMonitor) Lock() {
	sam.mutex.Lock()
}

func (sam *ServerArgsMonitor) Unlock() {
	sam.mutex.Unlock()
}

func (sam *ServerArgsMonitor) RLock() {
	sam.mutex.RLock()
}

func (sam *ServerArgsMonitor) RUnlock() {
	sam.mutex.RUnlock()
}

func (sam *ServerArgsMonitor) LockAndGetData() any {
	sam.Lock()
	return sam.args
}

func (sam *ServerArgsMonitor) GetData() any {
	return sam.args
}

func (sam *ServerArgsMonitor) Load(filepath string) error {
	sam.Lock()
	defer sam.Unlock()
	return loadServerConfigObject(sam, func(obj ServerConfigObject) {
		argsMonitor, ok := obj.(*ServerArgsMonitor)
		if ok {
			argsMonitor.args = NewServerArgs()
		}
	}, filepath)
}

func (sam *ServerArgsMonitor) Save(filepath string) error {
	sam.mutex.Lock()
	defer sam.mutex.Unlock()
	return saveServerConfigObject(sam, filepath)
}

func BuildStringArgs(args *ServerArgs) []string {
	strArgs := []string{"java"}
	version := serverConfig.LockAndGetData().(*ServerConfig).Version
	serverConfig.Unlock()

	strArgs = append(
		strArgs,
		fmt.Sprintf("-Xms%dG", args.MemoryStartGB),
		fmt.Sprintf("-Xmx%dG", args.MemoryMaxGB),
		"-jar",
		fmt.Sprintf("server-%s.jar", version),
		"--nogui",
	)

	if args.BonusChest {
		strArgs = append(strArgs, "--bonusChest")
	}
	if args.Demo {
		strArgs = append(strArgs, "--demo")
	}
	if args.EraseCache {
		strArgs = append(strArgs, "--eraseCache")
	}
	if args.ForceUpgrade {
		strArgs = append(strArgs, "--forceUpgrade")
	}
	if args.SafeMode {
		strArgs = append(strArgs, "--safeMode")
	}
	if len(args.ServerID) > 0 {
		strArgs = append(strArgs, "--serverId", args.ServerID)
	}
	if len(args.SinglePlayer) > 0 {
		strArgs = append(strArgs, "--singleplayer", args.SinglePlayer)
	}
	if len(args.Universe) > 0 {
		strArgs = append(strArgs, "--universe", args.Universe)
	}
	if len(args.World) > 0 {
		strArgs = append(strArgs, "--world", args.World)
	}
	if args.Port > 0 {
		strArgs = append(strArgs, "--port", strconv.FormatInt(int64(args.Port), 10))
	}

	return strArgs
}
