package minecraft

import (
	"encoding/json"
	"io"

	"github.com/raian621/minecraft-server-controller/api"
)

func NewServerProperties() *api.ServerProperties {
	return &api.ServerProperties{
		AcceptTransfers:                ref(false),
		AllowFlight:                    ref(false),
		AllowNether:                    ref(true),
		BroadcastConsoleToOps:          ref(true),
		BroadcastRCONToOps:             ref(true),
		Difficulty:                     ref[api.ServerPropertiesDifficulty]("easy"),
		EnableCommandBlock:             ref(false),
		EnableJMXMonitoring:            ref(false),
		EnableQuery:                    ref(false),
		EnableRCON:                     ref(false),
		EnableStatus:                   ref(true),
		EnforceSecureProfile:           ref(true),
		EnforceWhitelist:               ref(false),
		EntityBroadcastRangePercentage: ref(100),
		ForceGamemode:                  ref(false),
		FunctionPermissionLevel:        ref(2),
		Gamemode:                       ref[api.ServerPropertiesGamemode]("survival"),
		GenerateStructures:             ref(true),
		GeneratorSettings:              ref("{}"),
		Hardcore:                       ref(false),
		HideOnlinePlayers:              ref(false),
		InitialDisabledPacks:           ref(""),
		InitialEnabledPacks:            ref("vanilla"),
		LevelName:                      ref("world"),
		LevelSeed:                      ref(""),
		LevelType:                      ref("minecraft\\:normal"),
		LogIPs:                         ref(true),
		MaxChainedNeighborUpdates:      ref(1000000),
		MaxPlayers:                     ref(20),
		MaxTickTime:                    ref(60000),
		MaxWorldSize:                   ref(29999984),
		MOTD:                           ref("A Minecraft Server"),
		NetworkCompressionThreshold:    ref(256),
		OnlineMode:                     ref(true),
		OpPermissionLevel:              ref(4),
		PlayerIdleTimeout:              ref(0),
		PreventProxyConnections:        ref(false),
		PreviewsChat:                   ref(false),
		PVP:                            ref(true),
		QueryPort:                      ref(25565),
		RateLimit:                      ref(0),
		RCONPassword:                   ref(""),
		RCONPort:                       ref(25575),
		RegionFileCompression:          ref("deflate"),
		ResourcePack:                   ref(""),
		RequireResourcePack:            ref(false),
		ResourcePackID:                 ref(""),
		ResourcePackPrompt:             ref(""),
		ResourcePackSHA1:               ref(""),
		ServerIP:                       ref(""),
		ServerPort:                     ref(25565),
		SimulationDistance:             ref(10),
		SnooperEnabled:                 ref(true),
		SpawnAnimals:                   ref(true),
		SpawnMonsters:                  ref(true),
		SpawnNPCs:                      ref(true),
		SpawnProtection:                ref(16),
		SyncChunkWrites:                ref(true),
		TextFilteringConfig:            ref(""),
		UseNativeTransport:             ref(true),
		ViewDistance:                   ref(10),
		Whitelist:                      ref(false),
	}
}

func (m *JavaMinecraftServer) Properties() *api.ServerProperties {
	m.Lock()
	defer m.Unlock()

	if m.properties == nil {
		return nil
	}

	propsCpy := *m.properties

	return &propsCpy
}

// UpdateProperties implements api.MinecraftServerInterface.
func (m *JavaMinecraftServer) UpdateProperties(props api.ServerProperties) error {
	panic("unimplemented")
}

func (m *JavaMinecraftServer) CreateProperties() {
	m.Lock()
	defer m.Unlock()

	m.properties = NewServerProperties()
}

func (m *JavaMinecraftServer) LoadProperties(file io.Reader) error {
	m.Lock()
	defer m.Unlock()

	return json.NewDecoder(file).Decode(m.properties)
}

func (m *JavaMinecraftServer) SaveProperties(file io.Writer) error {
	m.Lock()
	defer m.Unlock()

	if m.properties == nil {
		return ErrNilConfig
	}

	return json.NewEncoder(file).Encode(m.properties)
}

func (m *JavaMinecraftServer) SetProperties(properties *api.ServerProperties) {
	m.Lock()
	defer m.Unlock()

	m.properties = properties
}
