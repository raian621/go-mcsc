package minecraft

import (
	"errors"
	"io"
	"os"
	"path"
	"sync"

	"github.com/raian621/go-mcsc/api"
)

var (
	ErrFilepathsNotProvided error = errors.New("filepaths for configuration files not provided")
	ErrNilConfig            error = errors.New("config object was not initialized")
)

type MinecraftServer struct {
	allowlist     *api.Allowlist
	args          *api.ServerArguments
	bannedIPs     *api.BannedIPList
	bannedPlayers *api.BannedPlayerList
	config        *api.MinecraftServerConfig
	filepaths     *MinecraftServerConfigFilepaths
	ops           *api.ServerOperatorList
	properties    *api.ServerProperties
	console       *Console
	// process      *exec.Cmd
	versions *[]string

	mutex sync.Mutex
}

type MinecraftServerConfigFilepaths struct {
	Allowlist          string
	Args               string
	BannedPlayers      string
	BannedIPs          string
	Config             string
	Ops                string
	Properties         string
	PropertiesTemplate string
	Versions           string
}

type JavaMinecraftServer MinecraftServer

func (m *JavaMinecraftServer) Lock()    { m.mutex.Lock() }
func (m *JavaMinecraftServer) Unlock()  { m.mutex.Unlock() }
func (m *JavaMinecraftServer) TryLock() { m.mutex.TryLock() }

func (m *JavaMinecraftServer) LoadConfigs() error {
	if m.filepaths == nil {
		return ErrFilepathsNotProvided
	}

	loadData := []struct {
		LoadFn   func(file io.Reader) error
		SaveFn   func(file io.Writer) error
		CreateFn func()
		Filepath string
	}{
		{
			LoadFn:   m.LoadAllowlist,
			CreateFn: m.CreateAllowlist,
			SaveFn:   m.SaveAllowlist,
			Filepath: m.filepaths.Allowlist,
		},
		{
			LoadFn:   m.LoadArgs,
			CreateFn: m.CreateArgs,
			SaveFn:   m.SaveArgs,
			Filepath: m.filepaths.Args,
		},
		{
			LoadFn:   m.LoadBannedIPs,
			CreateFn: m.CreateBannedIPs,
			SaveFn:   m.SaveBannedIPs,
			Filepath: m.filepaths.Allowlist,
		},
		{
			LoadFn:   m.LoadBannedPlayers,
			CreateFn: m.CreateBannedPlayers,
			SaveFn:   m.SaveBannedPlayers,
			Filepath: m.filepaths.BannedPlayers,
		},
		{
			LoadFn:   m.LoadConfig,
			CreateFn: m.CreateConfig,
			SaveFn:   m.SaveConfig,
			Filepath: m.filepaths.Config,
		},
		{
			LoadFn:   m.LoadOperators,
			CreateFn: m.CreateOperators,
			SaveFn:   m.SaveOperators,
			Filepath: m.filepaths.Ops,
		},
		{
			LoadFn:   m.LoadProperties,
			CreateFn: m.CreateProperties,
			SaveFn:   m.SaveProperties,
			Filepath: m.filepaths.Properties,
		},
		{
			LoadFn:   m.LoadVersions,
			CreateFn: m.CreateVersions,
			SaveFn:   func(file io.Writer) error { return nil },
			Filepath: m.filepaths.Versions,
		},
	}

	for _, ld := range loadData {
		if err := loadJSON(ld.LoadFn, ld.SaveFn, ld.CreateFn, ld.Filepath); err != nil {
			return err
		}
	}

	return saveServerPropertiesTemplate(
		m.properties,
		m.filepaths.PropertiesTemplate,
		path.Join(path.Dir(m.filepaths.Properties), "server.properties"),
	)
}

// Restart implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Restart() error {
	panic("unimplemented")
}

// Start implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Start() error {
	m.Lock()
	defer m.Unlock()
	return nil
	// version := serverConfig.LockAndGetData().(*ServerConfig).Version
	// serverConfig.Unlock()
	// if err := downloadServerVersion("server-data", version); err != nil {
	// 	return err
	// }

	// if serverCmd != nil && serverCmd.Process.Signal(syscall.Signal(0)) == nil {
	// 	return nil
	// }

	// if err := stopExistingMinecraftServers(); err != nil {
	// 	return err
	// }

	// args := serverArgs.LockAndGetData().(*ServerArgs)
	// argsStr := BuildStringArgs(args)
	// log.Println("args", argsStr)
	// serverArgs.Unlock()
	// serverCmd = exec.Command(argsStr[0], argsStr[1:]...)

	// go func() {
	// 	log.Println("starting minecraft server process...")

	// 	oldDir, err := os.Getwd()
	// 	if err != nil {
	// 		log.Println(err)
	// 		return
	// 	}
	// 	if err := os.Chdir("server-data"); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	ServerConsole = new(ConsoleMonitor)

	// 	ServerConsole.Stdin, err = serverCmd.StdinPipe()
	// 	if err != nil {
	// 		log.Fatalln("error creating stdin pipe for minecraft server:", err)
	// 	}
	// 	ServerConsole.Stdout, err = serverCmd.StdoutPipe()
	// 	if err != nil {
	// 		log.Fatalln("error creating stdout pipe for minecraft server:", err)
	// 	}
	// 	ServerConsole.Stderr, err = serverCmd.StderrPipe()
	// 	if err != nil {
	// 		log.Fatalln("error creating stderr pipe for minecraft server:", err)
	// 	}

	// 	if err := serverCmd.Start(); err != nil {
	// 		log.Println("error starting server:", err)
	// 		serverCmd = nil
	// 		return
	// 	}
	// 	if err := os.Chdir(oldDir); err != nil {
	// 		log.Println(err)
	// 		return
	// 	}

	// 	log.Println("minecraft server process started")
	// }()

	// return nil
}

// Start implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) Stop() error {
	panic("unimplemented")
}

func NewJavaMinecraftServer(filepaths *MinecraftServerConfigFilepaths) api.MinecraftServerInterface {
	jms := JavaMinecraftServer{filepaths: filepaths}
	return &jms
}

var _ api.MinecraftServerInterface = (*JavaMinecraftServer)(nil)

// func StartMinecraftServer() error {
// 	version := serverConfig.LockAndGetData().(*ServerConfig).Version
// 	serverConfig.Unlock()
// 	if err := downloadServerVersion("server-data", version); err != nil {
// 		return err
// 	}

// 	if serverCmd != nil && serverCmd.Process.Signal(syscall.Signal(0)) == nil {
// 		return nil
// 	}

// 	if err := stopExistingMinecraftServers(); err != nil {
// 		return err
// 	}

// 	args := serverArgs.LockAndGetData().(*ServerArgs)
// 	argsStr := BuildStringArgs(args)
// 	log.Println("args", argsStr)
// 	serverArgs.Unlock()
// 	serverCmd = exec.Command(argsStr[0], argsStr[1:]...)

// 	go func() {
// 		log.Println("starting minecraft server process...")

// 		oldDir, err := os.Getwd()
// 		if err != nil {
// 			log.Println(err)
// 			return
// 		}
// 		if err := os.Chdir("server-data"); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		ServerConsole = new(ConsoleMonitor)

// 		ServerConsole.Stdin, err = serverCmd.StdinPipe()
// 		if err != nil {
// 			log.Fatalln("error creating stdin pipe for minecraft server:", err)
// 		}
// 		ServerConsole.Stdout, err = serverCmd.StdoutPipe()
// 		if err != nil {
// 			log.Fatalln("error creating stdout pipe for minecraft server:", err)
// 		}
// 		ServerConsole.Stderr, err = serverCmd.StderrPipe()
// 		if err != nil {
// 			log.Fatalln("error creating stderr pipe for minecraft server:", err)
// 		}

// 		if err := serverCmd.Start(); err != nil {
// 			log.Println("error starting server:", err)
// 			serverCmd = nil
// 			return
// 		}
// 		if err := os.Chdir(oldDir); err != nil {
// 			log.Println(err)
// 			return
// 		}

// 		log.Println("minecraft server process started")
// 	}()

// 	return nil
// }

// func StopMinecraftServer() {
// 	if serverCmd != nil && serverCmd.Process.Signal(syscall.Signal(0)) == nil {
// 		if ServerConsole != nil {
// 			if !ServerConsole.TryLock() {

// 				return
// 			}

// 			if _, err := fmt.Fprintln(ServerConsole.Stdin, "/stop"); err != nil {
// 				log.Println("error sending stop command to server:", err)
// 			}

// 			ServerConsole.Stdin.Close()
// 			ServerConsole.Stdout.Close()
// 			ServerConsole.Stderr.Close()

// 			ServerConsole.Unlock()
// 			ServerConsole = nil
// 		}

// 		serverCmd = nil
// 	}
// }

// func stopExistingMinecraftServers() error {
// 	cmd := exec.Command("ps", "-ao", "args,pid")
// 	outBytes, err := cmd.Output()
// 	if err != nil {
// 		return err
// 	}
// 	// if there's no output
// 	if len(outBytes) == 0 {
// 		return nil
// 	}

// 	outStr := string(outBytes)
// 	commandData := strings.Split(outStr, "\n")[1:]
// 	pattern, err := regexp.Compile(`^java .*-jar server.jar.*`)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	for _, cd := range commandData {
// 		if pattern.MatchString(cd) {
// 			split := strings.Split(cd, " ")
// 			pidStr := split[len(split)-1]
// 			pid, err := strconv.ParseInt(pidStr, 10, 64)

// 			if err != nil {
// 				log.Println(err)
// 				continue
// 			}
// 			log.Println("killing minecraft server process with PID", pid)
// 			process, err := os.FindProcess(int(pid))
// 			if err != nil {
// 				log.Println(err)
// 				continue
// 			}

// 			if err := process.Kill(); err != nil {
// 				log.Println(err)
// 			}
// 		}
// 	}

// 	return nil
// }

func CreateServerFolder(filepath string) error {
	if err := os.Mkdir(filepath, os.ModePerm); !errors.Is(err, os.ErrExist) && err != nil {
		return err
	} else if errors.Is(err, os.ErrExist) {
		return nil
	}

	if err := os.WriteFile(path.Join(filepath, "eula.txt"), []byte("eula=TRUE\n"), os.ModePerm); err != nil {
		return err
	}

	return nil
}
