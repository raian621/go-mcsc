package minecraft

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/raian621/minecraft-server-controller/api"
)

func TestArgs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if server.Args() != nil {
		t.Fatal("expected args to be nil")
	}

	server.CreateArgs()
	args := server.Args()
	if args == nil {
		t.Fatal("expected args not to be nil")
	}
	if args == server.args {
		t.Fatal("returned args and server arg memory addresses should be different")
	}
}

func TestBuildArgs(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		version  string
		args     api.ServerArguments
		wantArgs []string
	}{
		{
			name:    "no extra args",
			version: "1.20.6",
			args:    *NewServerArgs(),
			wantArgs: []string{
				"java",
				"-Xms1G",
				"-Xmx2G",
				"-jar",
				"server-1.20.6.jar",
				"--nogui",
			},
		},
		{
			name:    "all extra args",
			version: "1.20.6",
			args: api.ServerArguments{
				BonusChest:    ref(true),
				Demo:          ref(true),
				EraseCache:    ref(true),
				ForceUpgrade:  ref(true),
				SafeMode:      ref(true),
				ServerID:      ref("server_id"),
				SinglePlayer:  ref("single_and_pringle"),
				Universe:      ref("world"),
				World:         ref("a_world"),
				Port:          ref(1000),
				MemoryStartGB: ref(1),
				MemoryMaxGB:   ref(2),
			},
			wantArgs: []string{
				"java",
				"-Xms1G",
				"-Xmx2G",
				"-jar",
				"server-1.20.6.jar",
				"--nogui",
				"--bonusChest",
				"--demo",
				"--eraseCache",
				"--forceUpgrade",
				"--safeMode",
				"--serverId",
				"server_id",
				"--singleplayer",
				"single_and_pringle",
				"--universe",
				"world",
				"--world",
				"a_world",
				"--port",
				"1000",
			},
		},
		{
			name:    "all args off",
			version: "1.20.6",
			args: api.ServerArguments{
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
			},
			wantArgs: []string{
				"java",
				"-Xms1G",
				"-Xmx2G",
				"-jar",
				"server-1.20.6.jar",
				"--nogui",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotArgs := BuildStringArgs(tc.version, &tc.args)
			if !reflect.DeepEqual(tc.wantArgs, gotArgs) {
				t.Fatalf("expected `%v` args, got `%v`", tc.wantArgs, gotArgs)
			}
		})
	}
}

func TestLoadArgs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	mockedJSON := bytes.NewBuffer([]byte(`{"memoryMaxGB":2,"memoryStartGB":1}`))
	if err := server.LoadArgs(mockedJSON); err != ErrNilConfig {
		t.Fatalf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.CreateArgs()
	if err := server.LoadArgs(mockedJSON); err != nil {
		t.Fatalf("expected no error, got `%v`", err)
	}

	// TODO: check server.args is correct
	if server.args.BonusChest != nil {
		t.Error("expected nil pointer for BonusChest")
	}
	if server.args.Demo != nil {
		t.Error("expected nil pointer for Demo")
	}
	if server.args.EraseCache != nil {
		t.Error("expected nil pointer for EraseCache")
	}
	if server.args.ForceUpgrade != nil {
		t.Error("expected nil pointer for ForceUpgrade")
	}
	if server.args.MemoryMaxGB == nil {
		t.Error("expected non-nil pointer for MemoryMaxGB")
	} else if *server.args.MemoryMaxGB != 2 {
		t.Errorf("expected MemoryMaxGB = `2`, got `%d`", *server.args.MemoryMaxGB)
	}
	if server.args.MemoryStartGB == nil {
		t.Error("expected non-nil pointer for MemoryStartGB")
	} else if *server.args.MemoryStartGB != 1 {
		t.Errorf("expected MemoryStartGB = `1`, got `%d`", *server.args.MemoryStartGB)
	}
	if server.args.Port != nil {
		t.Error("expected nil pointer for Port")
	}
	if server.args.SafeMode != nil {
		t.Error("expected nil pointer for SafeMode")
	}
	if server.args.ServerID != nil {
		t.Error("expected nil pointer for ServerID")
	}
	if server.args.SinglePlayer != nil {
		t.Error("expected nil pointer for SinglePlayer")
	}
	if server.args.Universe != nil {
		t.Error("expected nil pointer for Universe")
	}
	if server.args.World != nil {
		t.Error("expected nil pointer for World")
	}
}

func TestSaveArgs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	var out bytes.Buffer
	if err := server.SaveArgs(&out); err != ErrNilConfig {
		t.Fatalf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	out.Reset()
	server.CreateArgs()
	if err := server.SaveArgs(&out); err != nil {
		t.Fatalf("expected no error, got `%v`", err)
	}

	expected := `{"memoryMaxGB":2,"memoryStartGB":1}` + "\n"
	if expected != out.String() {
		t.Fatalf("expected `%s` output, got `%s`", expected, out.String())
	}
}

func TestSetArgs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	args := NewServerArgs()

	server.SetArgs(args)

	if args != server.args {
		t.Fatal("args and server args should have the same memory address")
	}
}
