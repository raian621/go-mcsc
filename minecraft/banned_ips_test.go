package minecraft

import (
	"bytes"
	"testing"

	"github.com/raian621/minecraft-server-controller/api"
)

func TestBanAndPardonIPs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.BanIP(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}
	if err := server.PardonIP("localhost"); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	testCases := []struct {
		name          string
		bannedIPs     api.BannedIPList
		banIPs        []api.BannedIP
		pardonIPs     []string
		wantPardonErr error
		wantBannedIPs api.BannedIPList
	}{
		{
			name: "add ip to ban list",
			banIPs: []api.BannedIP{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Ip:      "127.0.0.69",
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
			wantBannedIPs: api.BannedIPList{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Ip:      "127.0.0.69",
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
		},
		{
			name:          "try to pardon IP with empty list",
			bannedIPs:     make(api.BannedIPList, 0),
			pardonIPs:     []string{"127.0.0.69"},
			wantPardonErr: ErrNotInBannedIPs,
		},
		{
			name: "pardon IP succeeds",
			bannedIPs: api.BannedIPList{
				{
					Created: "2024-05-28T14:37:43-05:00",
					Expires: "forever",
					Ip:      "127.0.0.69",
					Reason:  "erm, idk",
					Source:  "server",
				},
			},
			pardonIPs:     []string{"127.0.0.69"},
			wantBannedIPs: make(api.BannedIPList, 0),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := JavaMinecraftServer{}
			server.bannedIPs = &tc.bannedIPs

			for _, b := range tc.banIPs {
				b := b

				if err := server.BanIP(&b); err != nil {
					t.Errorf("expected no error, got `%v`", err)
				}
			}

			for _, ip := range tc.pardonIPs {
				if err := server.PardonIP(ip); err != tc.wantPardonErr {
					t.Errorf("expected error `%v`, got `%v`", tc.wantPardonErr, err)
				}
			}

			if len(*server.bannedIPs) != len(tc.bannedIPs) {
				t.Fatalf(
					"expected banned IP list to contain %d elements, got %d elements",
					len(*server.bannedIPs),
					len(tc.bannedIPs),
				)
			}

			for i := range tc.banIPs {
				expected := tc.bannedIPs[i]
				got := (*server.bannedIPs)[i]

				if expected.Created != got.Created {
					t.Errorf("expected creation date `%s`, got `%s`", expected.Created, got.Created)
				}
				if expected.Expires != got.Expires {
					t.Errorf("expected expiration date `%s`, got `%s`", expected.Expires, got.Expires)
				}
				if expected.Ip != got.Ip {
					t.Errorf("expected ip `%s`, got `%s`", expected.Ip, got.Ip)
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

func TestLoadBannedIPs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.LoadBannedIPs(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	bannedIPs := api.BannedIPList{
		{
			Created: "now",
			Expires: "forever",
			Ip:      "127.0.0.2",
			Reason:  "was bad i guess",
			Source:  "Brian",
		},
		{
			Created: "now",
			Expires: "forever",
			Ip:      "127.0.0.3",
			Reason:  "very evil",
			Source:  "Jeff",
		},
	}

	server.CreateBannedIPs()
	mockedJSON := bytes.NewBufferString(`[{"created":"now","expires":"forever","ip":"127.0.0.2","reason":"was bad i guess","source":"Brian"},{"created":"now","expires":"forever","ip":"127.0.0.3","reason":"very evil","source":"Jeff"}]` + "\n")
	if err := server.LoadBannedIPs(mockedJSON); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}

	if len(bannedIPs) != len(*server.bannedIPs) {
		t.Fatalf("expected `%d` banned IPs, got `%d`", len(bannedIPs), len(*server.bannedIPs))
	}

	for i, want := range bannedIPs {
		got := (*server.bannedIPs)[i]

		if want.Created != got.Created {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Created, got.Created, i)
		}
		if want.Expires != got.Expires {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Expires, got.Expires, i)
		}
		if want.Ip != got.Ip {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Ip, got.Ip, i)
		}
		if want.Reason != got.Reason {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Reason, got.Reason, i)
		}
		if want.Source != got.Source {
			t.Errorf("expected `%s`, got `%s` at index %d", want.Source, got.Source, i)
		}
	}
}

func TestSaveBannedIPs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.SaveBannedIPs(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	server.bannedIPs = &api.BannedIPList{
		{
			Created: "now",
			Expires: "forever",
			Ip:      "127.0.0.2",
			Reason:  "was bad i guess",
			Source:  "Brian",
		},
		{
			Created: "now",
			Expires: "forever",
			Ip:      "127.0.0.3",
			Reason:  "very evil",
			Source:  "Jeff",
		},
	}

	var out bytes.Buffer
	if err := server.SaveBannedIPs(&out); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}

	expected := `[{"created":"now","expires":"forever","ip":"127.0.0.2","reason":"was bad i guess","source":"Brian"},{"created":"now","expires":"forever","ip":"127.0.0.3","reason":"very evil","source":"Jeff"}]` + "\n"
	if expected != out.String() {
		t.Errorf("expected `%s`, got `%v`", expected, out.String())
	}
}

func TestSetAndGetBannedIPs(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	bannedIPs := api.BannedIPList{}

	if shouldBeNil := server.BannedIPs(); shouldBeNil != nil {
		t.Errorf("banned ips pointer should have been nil")
	}

	server.SetBannedIPs(&bannedIPs)
	if server.bannedIPs != &bannedIPs {
		t.Errorf("expected memory address `%p`, got `%p`", &bannedIPs, server.bannedIPs)
	}

	bannedIPsCpy := server.BannedIPs()
	if bannedIPsCpy == server.bannedIPs {
		t.Errorf("memory addresses should have been different")
	}
}
