package minecraft

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// set up server
	// CreateServerFolder("../test-server-data")

	// testServer = NewJavaMinecraftServer(&MinecraftServerConfigFilepaths{
	// 	Allowlist:          "../test-server-data/allowlist.json",
	// 	Args:               "../test-server-data/args.json",
	// 	BannedIPs:          "../test-server-data/banned-ips.json",
	// 	BannedPlayers:      "../test-server-data/banned-players.json",
	// 	Config:             "../test-server-data/config.json",
	// 	Ops:                "../test-server-data/ops.json",
	// 	Properties:         "../test-server-data/properties.json",
	// 	PropertiesTemplate: "../templates/server.properties.tmpl",
	// 	Versions:           "data/server-download-links.json",
	// }).(*JavaMinecraftServer)

	// testServer.LoadConfigs()

	os.Exit(m.Run())
}
