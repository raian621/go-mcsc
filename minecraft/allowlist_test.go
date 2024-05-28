package minecraft

import (
	"bytes"
	"testing"

	"github.com/raian621/minecraft-server-controller/api"
)

func TestAllowlistAddAndDelete(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	server.CreateAllowlist()

	testCases := []struct {
		name          string
		addPlayers    []api.PlayerInfo
		removePlayers []api.PlayerInfo
		wantPlayers   api.Allowlist
		wantDeleteErr error
	}{
		{
			name:       "delete from empty allowlist",
			addPlayers: make([]api.PlayerInfo, 0),
			removePlayers: []api.PlayerInfo{
				{
					Name: ref("player"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
			},
			wantPlayers:   make(api.Allowlist, 0),
			wantDeleteErr: ErrPlayerNotInAllowlist,
		},
		{
			name:       "delete from empty allowlist",
			addPlayers: make([]api.PlayerInfo, 0),
			removePlayers: []api.PlayerInfo{
				{
					Name: ref("player"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
			},
			wantPlayers:   make(api.Allowlist, 0),
			wantDeleteErr: ErrPlayerNotInAllowlist,
		},
		{
			name:       "delete from empty allowlist",
			addPlayers: make([]api.PlayerInfo, 0),
			removePlayers: []api.PlayerInfo{
				{
					Name: ref("player"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
			},
			wantPlayers:   make(api.Allowlist, 0),
			wantDeleteErr: ErrPlayerNotInAllowlist,
		},
		{
			name: "add to allowlist",
			addPlayers: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
				{
					Name: ref("player2"),
					Uuid: ref("d0d1ead8-28d2-4197-bf38-65bcf7319365"),
				},
				{
					Name: ref("player3"),
					Uuid: ref("8be57f30-efa4-49bc-a341-96d4fc4d7ec7"),
				},
			},
			removePlayers: make([]api.PlayerInfo, 0),
			wantPlayers: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
				{
					Name: ref("player2"),
					Uuid: ref("d0d1ead8-28d2-4197-bf38-65bcf7319365"),
				},
				{
					Name: ref("player3"),
					Uuid: ref("8be57f30-efa4-49bc-a341-96d4fc4d7ec7"),
				},
			},
		},
		{
			name:       "delete from allowlist",
			addPlayers: make([]api.PlayerInfo, 0),
			removePlayers: []api.PlayerInfo{
				{
					Name: ref("player2"),
					Uuid: ref("d0d1ead8-28d2-4197-bf38-65bcf7319365"),
				},
				{
					Name: ref("player3"),
					Uuid: ref("8be57f30-efa4-49bc-a341-96d4fc4d7ec7"),
				},
			},
			wantPlayers: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
			},
		},
		{
			name:       "try to delete player from allowlist that isn't in allowlist",
			addPlayers: make([]api.PlayerInfo, 0),
			removePlayers: []api.PlayerInfo{
				{
					Name: ref("player2"),
					Uuid: ref("d0d1ead8-28d2-4197-bf38-65bcf7319365"),
				},
			},
			wantPlayers: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref("052acc86-065d-49f6-b518-c508a7cf55ae"),
				},
			},
			wantDeleteErr: ErrPlayerNotInAllowlist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, player := range tc.addPlayers {
				if err := server.AllowPlayer(player); err != nil {
					t.Fatalf("expected no error, got `%v`", err)
				}
			}

			for _, player := range tc.removePlayers {
				if err := server.DisallowPlayer(player); err != tc.wantDeleteErr {
					t.Fatalf("expected error `%v`, got `%v`", tc.wantDeleteErr, err)
				}
			}

			players := server.Allowlist()
			if players == nil {
				t.Fatal("allowlist not initialized")
			}
			if len(tc.wantPlayers) != len(*players) {
				t.Fatalf("expected allowlist `%v`, got `%v`", tc.wantPlayers, *players)
			}
			for i := range tc.wantPlayers {
				wantPlayer := tc.wantPlayers[i]
				gotPlayer := (*players)[i]

				if *wantPlayer.Name != *gotPlayer.Name || *wantPlayer.Uuid != *gotPlayer.Uuid {
					t.Fatalf("expected allowlist `%v`, got `%v`", tc.wantPlayers, *players)
				}
			}
		})
	}
}

func TestAllowlistLoad(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	mockAllowlistJSON := bytes.NewBuffer([]byte(
		`[
	{
		"name": "player1",
		"uuid": "f3058911-2a7b-4928-b1f7-5255d586c227"
	},
	{
		"name": "player2",
		"uuid": "f3058911-2a7b-4928-b1f7-5255d586c228"
	}
]`,
	))

	if err := server.LoadAllowlist(mockAllowlistJSON); err != ErrNilConfig {
		t.Fatalf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.CreateAllowlist()

	if err := server.LoadAllowlist(mockAllowlistJSON); err != nil {
		t.Fatalf("expected no error, got `%v`", err)
	}

	expectedAllowlist := api.Allowlist{
		{
			Name: ref("player1"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c227"),
		},
		{
			Name: ref("player2"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c228"),
		},
	}

	if len(expectedAllowlist) != len(*server.allowlist) {
		t.Fatalf(
			"wanted an allowlist of length %d, got length %d",
			len(expectedAllowlist), len(*server.allowlist),
		)
	}

	equal := true
	for i := range expectedAllowlist {
		want := expectedAllowlist[i]
		got := (*server.allowlist)[i]

		if *want.Name != *got.Name || *want.Uuid != *got.Uuid {
			equal = false
			t.Errorf(
				"wanted { *Name: %s, *Uuid: %s }, got { *Name: %s, *Uuid: %s } at index %d",
				*want.Name,
				*want.Uuid,
				*got.Name,
				*got.Uuid,
				i,
			)
		}
	}

	if !equal {
		t.FailNow()
	}
}

func TestAllowlistSave(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	var out bytes.Buffer

	if err := server.SaveAllowlist(&out); err != ErrNilConfig {
		t.Fatalf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.CreateAllowlist()
	server.allowlist = &api.Allowlist{
		{
			Name: ref("player1"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c227"),
		},
		{
			Name: ref("player2"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c228"),
		},
	}

	if err := server.SaveAllowlist(&out); err != nil {
		t.Fatalf("expected no error, got `%v`", err)
	}

	expected := `[{"name":"player1","uuid":"f3058911-2a7b-4928-b1f7-5255d586c227"},{"name":"player2","uuid":"f3058911-2a7b-4928-b1f7-5255d586c228"}]` + "\n"
	if expected != out.String() {
		t.Fatalf("expected `%s` allowlist, got `%s`", expected, out.String())
	}
}

func TestSetAllowlist(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	allowlist := api.Allowlist{
		{
			Name: ref("player1"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c227"),
		},
		{
			Name: ref("player2"),
			Uuid: ref("f3058911-2a7b-4928-b1f7-5255d586c228"),
		},
	}

	server.SetAllowlist(&allowlist)

	if &allowlist != server.allowlist {
		t.Fatalf("expected `%p` memory address, got `%p`", &allowlist, server.allowlist)
	}
}
