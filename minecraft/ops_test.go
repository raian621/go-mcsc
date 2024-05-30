package minecraft

import (
	"bytes"
	"testing"

	"github.com/google/uuid"
	"github.com/raian621/go-mcsc/api"
)

func TestOperatorsGetAndSet(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	operators := api.ServerOperatorList{}

	// should be uninitialized
	if server.Ops() != nil {
		t.Error("expected server operators list pointer to be nil")
	}

	server.SetOperators(&operators)
	serverOperators := server.Ops()

	if &operators == serverOperators {
		t.Error("memory addresses for original operators and serverOperators should have been different")
	}
}

func TestOpAndDeop(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}
	if err := server.Op(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}
	if err := server.Deop(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	testCases := []struct {
		name        string
		ops         api.ServerOperatorList
		op          []api.ServerOperator
		deop        []api.PlayerInfo
		wantDeopErr error
		wantOps     api.ServerOperatorList
	}{
		{
			name: "add server operator",
			ops:  make(api.ServerOperatorList, 0),
			op: []api.ServerOperator{
				{
					BypassesPlayerLimit: true,
					Level:               4,
					Name:                "player1",
					Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
				},
			},
			wantOps: api.ServerOperatorList{
				{
					BypassesPlayerLimit: true,
					Level:               4,
					Name:                "player1",
					Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
				},
			},
		},
		{
			name: "remove server operator",
			ops: api.ServerOperatorList{
				{
					BypassesPlayerLimit: true,
					Level:               4,
					Name:                "player1",
					Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
				},
			},
			deop: []api.PlayerInfo{
				{
					Name: ref("player1"),
					Uuid: ref(uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5")),
				},
			},
			wantOps: api.ServerOperatorList{},
		},
		{
			name: "remove nonexistent server operator",
			ops: api.ServerOperatorList{
				{
					BypassesPlayerLimit: true,
					Level:               4,
					Name:                "player1",
					Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
				},
			},
			deop: []api.PlayerInfo{
				{
					Name: ref("player2"),
					Uuid: ref(uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5")),
				},
			},
			wantOps: api.ServerOperatorList{
				{
					BypassesPlayerLimit: true,
					Level:               4,
					Name:                "player1",
					Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
				},
			},
			wantDeopErr: ErrNotInOps,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			server := JavaMinecraftServer{}
			server.ops = &tc.ops

			for _, op := range tc.op {
				op := op
				if err := server.Op(&op); err != nil {
					t.Errorf("expected no error, got `%v`", err)
				}
			}

			for _, p := range tc.deop {
				p := p
				if err := server.Deop(&p); err != tc.wantDeopErr {
					t.Errorf("expected error `%v`, got `%v`", tc.wantDeopErr, err)
				}
			}

			if len(*server.ops) != len(tc.wantOps) {
				t.Errorf("expected operator list to be length `%d`, got `%d`", len(tc.wantOps), len(*server.ops))
			}
		})
	}
}

func TestSaveAndLoadOps(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	if err := server.LoadOperators(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}
	if err := server.SaveOperators(nil); err != ErrNilConfig {
		t.Errorf("expected error `%v`, got `%v`", ErrNilConfig, err)
	}

	ops := api.ServerOperatorList{
		{
			BypassesPlayerLimit: true,
			Level:               4,
			Name:                "player1",
			Uuid:                uuid.MustParse("1b5c7df5-69c5-44d2-beab-8191f593e2e3"),
		},
		{
			BypassesPlayerLimit: true,
			Level:               4,
			Name:                "player2",
			Uuid:                uuid.MustParse("7b5c7df5-69c5-44d2-beab-8191f593e2e5"),
		},
	}

	server.ops = &ops
	var out bytes.Buffer
	expected := `[{"bypassesPlayerLimit":true,"level":4,"name":"player1","uuid":"1b5c7df5-69c5-44d2-beab-8191f593e2e3"},{"bypassesPlayerLimit":true,"level":4,"name":"player2","uuid":"7b5c7df5-69c5-44d2-beab-8191f593e2e5"}]` + "\n"
	if err := server.SaveOperators(&out); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}
	if expected != out.String() {
		t.Errorf("expected saved JSON to be `%s`, got `%s`", expected, out.String())
	}

	server.CreateOperators()

	if err := server.LoadOperators(&out); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}
	if len(*server.ops) != len(ops) {
		t.Errorf("expected operator list to be length `%d`, got `%d`", len(ops), len(*server.ops))
	}
}
