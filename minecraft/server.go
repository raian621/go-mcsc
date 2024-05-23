package minecraft

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

var (
	serverCmd        *exec.Cmd
	serverConfig     *ServerConfigMonitor
	serverArgs       *ServerArgsMonitor
	serverProperties *ServerPropertiesMonitor
	ServerConsole    *ConsoleMonitor
)

func init() {
	serverConfig = &ServerConfigMonitor{config: new(ServerConfig)}
	serverArgs = &ServerArgsMonitor{args: new(ServerArgs)}
	serverProperties = &ServerPropertiesMonitor{properties: new(ServerProperties)}
}

func StartMinecraftServer() error {
	version := serverConfig.LockAndGetData().(*ServerConfig).Version
	serverConfig.Unlock()
	if err := downloadServerVersion("server-data", version); err != nil {
		return err
	}

	if serverCmd != nil && serverCmd.Process.Signal(syscall.Signal(0)) == nil {
		return nil
	}

	if err := stopExistingMinecraftServers(); err != nil {
		return err
	}

	args := serverArgs.LockAndGetData().(*ServerArgs)
	argsStr := BuildStringArgs(args)
	log.Println("args", argsStr)
	serverArgs.Unlock()
	serverCmd = exec.Command(argsStr[0], argsStr[1:]...)

	go func() {
		log.Println("starting minecraft server process...")

		oldDir, err := os.Getwd()
		if err != nil {
			log.Println(err)
			return
		}
		if err := os.Chdir("server-data"); err != nil {
			log.Println(err)
			return
		}

		ServerConsole = new(ConsoleMonitor)

		ServerConsole.Stdin, err = serverCmd.StdinPipe()
		if err != nil {
			log.Fatalln("error creating stdin pipe for minecraft server:", err)
		}
		ServerConsole.Stdout, err = serverCmd.StdoutPipe()
		if err != nil {
			log.Fatalln("error creating stdout pipe for minecraft server:", err)
		}
		ServerConsole.Stderr, err = serverCmd.StderrPipe()
		if err != nil {
			log.Fatalln("error creating stderr pipe for minecraft server:", err)
		}

		if err := serverCmd.Start(); err != nil {
			log.Println("error starting server:", err)
			serverCmd = nil
			return
		}
		if err := os.Chdir(oldDir); err != nil {
			log.Println(err)
			return
		}

		log.Println("minecraft server process started")
		// 	<-serverStop

		// 	if err := serverCmd.Process.Kill(); err != nil {
		// 		log.Println("error killing server process:", err)
		// 	}

		// 	serverCmd = nil
		// 	log.Println("minecraft server process killed")
	}()

	return nil
}

func StopMinecraftServer() {
	if serverCmd != nil && serverCmd.Process.Signal(syscall.Signal(0)) == nil {
		if ServerConsole != nil {
			if !ServerConsole.TryLock() {

				return
			}

			if _, err := fmt.Fprintln(ServerConsole.Stdin, "/stop"); err != nil {
				log.Println("error sending stop command to server:", err)
			}

			ServerConsole.Stdin.Close()
			ServerConsole.Stdout.Close()
			ServerConsole.Stderr.Close()

			ServerConsole.Unlock()
			ServerConsole = nil
		}

		serverCmd = nil
	}
}

func stopExistingMinecraftServers() error {
	cmd := exec.Command("ps", "-ao", "args,pid")
	outBytes, err := cmd.Output()
	if err != nil {
		return err
	}
	// if there's no output
	if len(outBytes) == 0 {
		return nil
	}

	outStr := string(outBytes)
	commandData := strings.Split(outStr, "\n")[1:]
	pattern, err := regexp.Compile(`^java .*-jar server.jar.*`)
	if err != nil {
		log.Fatalln(err)
	}
	for _, cd := range commandData {
		if pattern.MatchString(cd) {
			split := strings.Split(cd, " ")
			pidStr := split[len(split)-1]
			pid, err := strconv.ParseInt(pidStr, 10, 64)

			if err != nil {
				log.Println(err)
				continue
			}
			log.Println("killing minecraft server process with PID", pid)
			process, err := os.FindProcess(int(pid))
			if err != nil {
				log.Println(err)
				continue
			}

			if err := process.Kill(); err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}

func LoadConfig(filepath string) error {
	log.Println("loading main config...")
	if err := serverConfig.Load(filepath); err != nil {
		return err
	}
	log.Println("loading server args...")
	if err := serverArgs.Load(serverConfig.config.ArgsFilepath); err != nil {
		return err
	}
	log.Println("loading server properties...")
	if err := serverProperties.Load(serverConfig.config.PropertiesFilepath); err != nil {
		return err
	}

	return nil
}

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
