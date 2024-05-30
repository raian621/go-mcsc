package minecraft

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/raian621/minecraft-server-controller/api"
)

func TestBanAndPardonPlayers(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.BanPlayer(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}
	if err := server.PardonPlayer(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	testCases := []struct {
		name              string
		bannedPlayers     api.BannedPlayerList
		banPlayers        []api.BannedPlayer
		pardonPlayers     []api.PlayerInfo
		wantPardonErr     error
		wantBannedPlayers api.BannedPlayerList
	}{
		{
			name: "add player to ban list",
			banPlayers: []api.BannedPlayer{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Name:    ref("player1"),
					Uuid:    uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
			wantBannedPlayers: api.BannedPlayerList{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Name:    ref("player1"),
					Uuid:    uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
			bannedPlayers: api.BannedPlayerList{},
		},
		{
			name:          "try to pardon player from empty list",
			bannedPlayers: make(api.BannedPlayerList, 0),
			pardonPlayers: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref(uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5")),
				},
			},
			wantPardonErr: ErrNotInBannedPlayers,
		},
		{
			name: "remove player from banned list",
			bannedPlayers: api.BannedPlayerList{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Name:    ref("player1"),
					Uuid:    uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
			pardonPlayers: []api.PlayerInfo{{
				Name: ref("player1"),
				Uuid: ref(uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5")),
			}},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := JavaMinecraftServer{}
			server.bannedPlayers = &tc.bannedPlayers

			for _, b := range tc.banPlayers {
				b := b

				if err := server.BanPlayer(&b); err != nil {
					t.Errorf("expected no error, got `%v`", err)
				}
			}

			for _, p := range tc.pardonPlayers {
				if err := server.PardonPlayer(&p); err != tc.wantPardonErr {
					t.Errorf("expected error `%v`, got `%v`", tc.wantPardonErr, err)
				}
			}

			if len(*server.bannedPlayers) != len(tc.bannedPlayers) {
				t.Fatalf(
					"expected banned IP list to contain %d elements, got %d elements",
					len(*server.bannedPlayers),
					len(tc.bannedPlayers),
				)
			}

			for i := range tc.bannedPlayers {
				expected := tc.bannedPlayers[i]
				got := (*server.bannedPlayers)[i]

				if expected.Created != got.Created {
					t.Errorf("expected creation date `%s`, got `%s`", expected.Created, got.Created)
				}
				if expected.Expires != got.Expires {
					t.Errorf("expected expiration date `%s`, got `%s`", expected.Expires, got.Expires)
				}
				if *expected.Name != *got.Name {
					t.Errorf("expected ip `%s`, got `%s`", *expected.Name, *got.Name)
				}
				if expected.Reason != got.Reason {
					t.Errorf("expected Reason `%s`, got `%s`", expected.Reason, got.Reason)
				}
				if expected.Source != got.Source {
					t.Errorf("expected Source `%s`, got `%s`", expected.Source, got.Source)
				}
			}
		})
	}
}

func TestLoadBannedPlayers(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.LoadBannedPlayers(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	bannedPlayers := api.BannedPlayerList{
		{
			Created: "2024-05-28T14:37:43-05:00",
			Expires: "forever",
			Name:    ref("player1"),
			Uuid:    uuid.MustParse("f3058911-2a7b-4928-b1f7-5255d586c228"),
			Reason:  "erm, idk",
			Source:  "server",
		},
		{
			Created: "2024-05-28T14:37:43-05:00",
			Expires: "forever",
			Name:    ref("player2"),
			Uuid:    uuid.MustParse("f3058911-2a7b-4928-b1f7-5255d586c227"),
			Reason:  "very evil",
			Source:  "Jeff",
		},
	}

	server.CreateBannedPlayers()
	mockedJSON := bytes.NewBufferString(`[{"created":"2024-05-28T14:37:43-05:00","expires":"forever","name":"player1","reason":"erm, idk","source":"server","uuid":"7b5c7df5-69c5-44d2-beab-8191f593e2e5"},{"created":"2024-05-28T14:37:43-05:00","expires":"forever","name":"player2","reason":"very evil","source":"Jeff","uuid":"f3058911-2a7b-4928-b1f7-5255d586c227"}]` + "\n")
	if err := server.LoadBannedPlayers(mockedJSON); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}

	if len(bannedPlayers) != len(*server.bannedPlayers) {
		t.Fatalf("expected `%d` banned players, got `%d`", len(bannedPlayers), len(*server.bannedPlayers))
	}

	for i, want := range bannedPlayers {
		got := (*server.bannedPlayers)[i]

		if want.Created != got.Created {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Created, got.Created, i)
		}
		if want.Expires != got.Expires {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Expires, got.Expires, i)
		}
		if *want.Name != *got.Name {
			t.Errorf("expected `%s`, got `%s` at index %d", *want.Name, *got.Name, i)
		}
		if want.Reason != got.Reason {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Reason, got.Reason, i)
		}
		if want.Source != got.Source {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Source, got.Source, i)
		}
	}
}

func TestSaveBannedPlayers(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.SaveBannedPlayers(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.bannedPlayers = &api.BannedPlayerList{
		{
			Created: "2024-05-28T14:37:43-05:00",
			Expires: "forever",
			Name:    ref("player1"),
			Uuid:    uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
			Reason:  "erm, idk",
			Source:  "server",
		},
		{
			Created: "2024-05-28T14:37:43-05:00",
			Expires: "forever",
			Name:    ref("player2"),
			Uuid:    uuid.MustParse("f3058911-2a7b-4928-b1f7-5255d586c227"),
			Reason:  "very evil",
			Source:  "Jeff",
		},
	}

	var out bytes.Buffer
	if err := server.SaveBannedPlayers(&out); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}

	expected := `[{"created":"2024-05-28T14:37:43-05:00","expires":"forever","name":"player1","reason":"erm, idk","source":"server","uuid":"7b5c7df5-69c5-44d2-beab-8191f593e2e5"},{"created":"2024-05-28T14:37:43-05:00","expires":"forever","name":"player2","reason":"very evil","source":"Jeff","uuid":"f3058911-2a7b-4928-b1f7-5255d586c227"}]` + "\n"
	if expected != out.String() {
		t.Errorf("expected `%s`, got `%v`", expected, out.String())
	}
}

func TestSetAndGetBannedPlayers(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	bannedPlayers := api.BannedPlayerList{}

	if shouldBeNil := server.BannedPlayers(); shouldBeNil != nil {
		t.Errorf("banned ips pointer should have been nil")
	}

	server.SetBannedPlayers(&bannedPlayers)
	if server.bannedPlayers != &bannedPlayers {
		t.Errorf("expected memory address `%p`, got `%p`", &bannedPlayers, server.bannedPlayers)
	}

	bannedPlayersCpy := server.BannedPlayers()
	if bannedPlayersCpy == server.bannedPlayers {
		t.Errorf("memory addresses should have been different")
	}
}
