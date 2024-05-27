package minecraft

import (
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
