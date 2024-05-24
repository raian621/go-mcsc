package minecraft

import (
	"io"
	"os"
	"path"
	"sync"
	"text/template"
)

// https://minecraft.fandom.com/wiki/Server.properties
type ServerProperties struct {
	AcceptTransfers                bool   `json:"acceptTransfers"`
	AllowFlight                    bool   `json:"allowFlight"`
	AllowNether                    bool   `json:"allowNether"`
	BroadcastConsoleToOps          bool   `json:"broadcastConsoleToOps"`
	BroadcastRCONToOps             bool   `json:"broadcastRCONToOps"`
	Difficulty                     string `json:"difficulty"`
	EnableCommandBlock             bool   `json:"enableCommandBlock"`
	EnableJMXMonitoring            bool   `json:"enableJMXMonitoring"`
	EnableQuery                    bool   `json:"enableQuery"`
	EnableRCON                     bool   `json:"enableRCON"`
	EnableStatus                   bool   `json:"enableStatus"`
	EnforceSecureProfile           bool   `json:"enforceSecureProfile"`
	EnforceWhitelist               bool   `json:"enforceWhitelist"`
	EntityBroadcastRangePercentage uint16 `json:"entityBroadcastRangePercentage"`
	ForceGamemode                  bool   `json:"forceGamemode"`
	FunctionPermissionLevel        uint8  `json:"functionPermissionLevel"`
	Gamemode                       string `json:"gamemode"`
	GenerateStructures             bool   `json:"generateStructures"`
	GeneratorSettings              string `json:"generatorSettings"`
	Hardcore                       bool   `json:"hardcore"`
	HideOnlinePlayers              bool   `json:"hideOnlinePlayers"`
	// Comma-separated list of datapacks to not be auto-enabled on world creation.
	InitialDisabledPacks string `json:"initialDisabledPacks"`
	// Comma-separated list of datapacks to be enabled during world creation. Feature packs need to be explicitly enabled.
	InitialEnabledPacks         string `json:"initialEnabledPacks"`
	LevelName                   string `json:"levelName"`
	LevelSeed                   string `json:"levelSeed"`
	LevelType                   string `json:"levelType"`
	LogIPs                      bool   `json:"logIPs"`
	MaxChainedNeighborUpdates   uint32 `json:"maxChainedNeighborUpdates"`
	MaxPlayers                  uint32 `json:"maxPlayers"`
	MaxTickTime                 int64  `json:"maxTickTime"`
	MaxWorldSize                uint32 `json:"maxWorldSize"`
	MOTD                        string `json:"MOTD"`
	NetworkCompressionThreshold int16  `json:"networkCompressionThreshold"`
	OnlineMode                  bool   `json:"onlineMode"`
	OpPermissionLevel           uint8  `json:"opPermissionLevel"`
	PlayerIdleTimeout           uint32 `json:"playerIdleTimeout"`
	PreventProxyConnections     bool   `json:"preventProxyConnections"`
	PreviewsChat                bool   `json:"previewsChat"`
	PVP                         bool   `json:"PVP"`
	QueryPort                   uint16 `json:"queryPort"`
	RateLimit                   uint32 `json:"rateLimit"`
	RCONPassword                string `json:"RCONPassword"`
	RCONPort                    uint16 `json:"RCONPort"`
	RegionFileCompression       string `json:"regionFileCompression"`
	RequireResourcePack         bool   `json:"requireResourcePack"`
	ResourcePack                string `json:"resourcePack"`
	ResourcePackID              string `json:"resourcePackID"`
	ResourcePackPrompt          string `json:"resourcePackPrompt"`
	ResourcePackSHA1            string `json:"resourcePackSHA1"`
	ServerIP                    string `json:"serverIP"`
	ServerPort                  uint16 `json:"serverPort"`
	SimulationDistance          uint8  `json:"simulationDistance"`
	SnooperEnabled              bool   `json:"snooperEnabled"`
	SpawnAnimals                bool   `json:"spawnAnimals"`
	SpawnMonsters               bool   `json:"spawnMonsters"`
	SpawnNPCs                   bool   `json:"spawnNPCs"`
	SpawnProtection             uint8  `json:"spawnProtection"`
	SyncChunkWrites             bool   `json:"syncChunkWrites"`
	TextFilteringConfig         string `json:"textFilteringConfig"`
	UseNativeTransport          bool   `json:"useNativeTransport"`
	ViewDistance                uint8  `json:"viewDistance"`
	Whitelist                   bool   `json:"whitelist"`
}

type ServerPropertiesMonitor struct {
	properties *ServerProperties
	mutex      sync.RWMutex
}

func NewServerProperties() *ServerProperties {
	return &ServerProperties{
		AcceptTransfers:                false,
		AllowFlight:                    false,
		AllowNether:                    true,
		BroadcastConsoleToOps:          true,
		BroadcastRCONToOps:             true,
		Difficulty:                     "easy",
		EnableCommandBlock:             false,
		EnableJMXMonitoring:            false,
		EnableQuery:                    false,
		EnableRCON:                     false,
		EnableStatus:                   true,
		EnforceSecureProfile:           true,
		EnforceWhitelist:               false,
		EntityBroadcastRangePercentage: 100,
		ForceGamemode:                  false,
		FunctionPermissionLevel:        2,
		Gamemode:                       "survival",
		GenerateStructures:             true,
		GeneratorSettings:              "{}",
		Hardcore:                       false,
		HideOnlinePlayers:              false,
		InitialDisabledPacks:           "",
		InitialEnabledPacks:            "vanilla",
		LevelName:                      "world",
		LevelSeed:                      "",
		LevelType:                      "minecraft\\:normal",
		LogIPs:                         true,
		MaxChainedNeighborUpdates:      1000000,
		MaxPlayers:                     20,
		MaxTickTime:                    60000,
		MaxWorldSize:                   29999984,
		MOTD:                           "A Minecraft Server",
		NetworkCompressionThreshold:    256,
		OnlineMode:                     true,
		OpPermissionLevel:              4,
		PlayerIdleTimeout:              0,
		PreventProxyConnections:        false,
		PreviewsChat:                   false,
		PVP:                            true,
		QueryPort:                      25565,
		RateLimit:                      0,
		RCONPassword:                   "",
		RCONPort:                       25575,
		RegionFileCompression:          "deflate",
		ResourcePack:                   "",
		RequireResourcePack:            false,
		ResourcePackID:                 "",
		ResourcePackPrompt:             "",
		ResourcePackSHA1:               "",
		ServerIP:                       "",
		ServerPort:                     25565,
		SimulationDistance:             10,
		SnooperEnabled:                 true,
		SpawnAnimals:                   true,
		SpawnMonsters:                  true,
		SpawnNPCs:                      true,
		SpawnProtection:                16,
		SyncChunkWrites:                true,
		TextFilteringConfig:            "",
		UseNativeTransport:             true,
		ViewDistance:                   10,
		Whitelist:                      false,
	}
}

func NewServerPropertiesMonitor() *ServerPropertiesMonitor {
	return &ServerPropertiesMonitor{
		properties: NewServerProperties(),
	}
}

func (spm *ServerPropertiesMonitor) Lock() {
	spm.mutex.Lock()
}

func (spm *ServerPropertiesMonitor) Unlock() {
	spm.mutex.Unlock()
}

func (spm *ServerPropertiesMonitor) RLock() {
	spm.mutex.RLock()
}

func (spm *ServerPropertiesMonitor) RUnlock() {
	spm.mutex.RUnlock()
}

func (spm *ServerPropertiesMonitor) LockAndGetData() any {
	spm.Lock()
	return spm.properties
}

func (spm *ServerPropertiesMonitor) GetData() any {
	return spm.properties
}

func (spm *ServerPropertiesMonitor) SaveToServerFile(
	templateFile string,
	file io.Writer,
) error {
	tmpl, err := template.New(path.Base(templateFile)).ParseFiles(templateFile)
	if err != nil {
		return err
	}

	spm.mutex.Lock()
	defer spm.mutex.Unlock()

	return tmpl.Execute(file, spm.properties)
}

func (spm *ServerPropertiesMonitor) Load(filepath string) error {
	spm.Lock()
	defer spm.Unlock()
	return loadServerConfigObject(spm, func(obj ServerConfigObject) {
		if propertiesMonitor, ok := obj.(*ServerPropertiesMonitor); ok {
			propertiesMonitor.properties = NewServerProperties()
		}
	}, filepath)
}

func (spm *ServerPropertiesMonitor) Save(filepath string) error {
	spm.Lock()
	defer spm.Unlock()

	tmpl, err := template.New("server.properties.tmpl").ParseFiles("minecraft/templates/server.properties.tmpl")
	if err != nil {
		return err
	}

	file, err := os.Create("server-data/server.properties")
	if err != nil {
		return nil
	}
	if err := tmpl.Execute(file, spm.properties); err != nil {
		return err
	}

	return saveServerConfigObject(spm, filepath)
}
