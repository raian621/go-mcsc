package minecraft

import (
	"bytes"
	"os"
	"reflect"
	"testing"
)

func TestGetVersion(t *testing.T) {
	t.Parallel()

	var expected = []string{
		"1.1", "1.10", "1.10.1", "1.10.2", "1.11", "1.11.1", "1.11.2", "1.12",
		"1.12.1", "1.12.2", "1.13", "1.13.1", "1.13.2", "1.14", "1.14.1", "1.14.2",
		"1.14.3", "1.14.4", "1.15", "1.15.1", "1.15.2", "1.16", "1.16.1", "1.16.2",
		"1.16.3", "1.16.4", "1.16.5", "1.17", "1.17.1", "1.18", "1.18.1", "1.18.2",
		"1.19", "1.19.1", "1.19.2", "1.19.3", "1.19.4", "1.2.5", "1.20", "1.20.1",
		"1.20.2", "1.20.3", "1.20.4", "1.20.5", "1.20.6", "1.3", "1.3.1", "1.3.2",
		"1.4", "1.4.1", "1.4.2", "1.4.3", "1.4.4", "1.4.5", "1.4.6", "1.4.7", "1.5",
		"1.5.1", "1.5.2", "1.6", "1.6.1", "1.6.2", "1.6.3", "1.6.4", "1.7", "1.7.1",
		"1.7.10", "1.7.2", "1.7.3", "1.7.4", "1.7.5", "1.7.6", "1.7.7", "1.7.8",
		"1.7.9", "1.8", "1.8.1", "1.8.2", "1.8.3", "1.8.4", "1.8.5", "1.8.6", "1.8.7",
		"1.8.8", "1.8.9", "1.9", "1.9.1", "1.9.2", "1.9.3", "1.9.4",
	}

	server := JavaMinecraftServer{}

	file, err := os.Open("../data/server-download-links.json")
	if err != nil {
		t.Fatal(err)
	}
	defer closeFile(file)

	if err := server.LoadVersions(file); err != nil {
		t.Errorf("expected no error, got `%v`", err)
	}

	versions := server.Versions()
	if !reflect.DeepEqual(expected, *versions) {
		t.Errorf("expected `%v`, got `%v`", expected, *versions)
	}
}

func TestSetVersion(t *testing.T) {
	t.Parallel()

	server := JavaMinecraftServer{}

	server.CreateConfig()

	err := server.LoadVersions(bytes.NewBuffer([]byte(
		"{\"1.20.6\":{\"link\":\"\",\"sum\":\"\"}}",
	)))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// supported version:
	if err := server.SetVersion("1.20.6"); err != nil {
		t.Fatalf("expected no error, got `%v`", err)
	}

	// check that the version was actually set:
	if server.config.Version != "1.20.6" {
		t.Fatalf("expected no `1.20.6`, got `%s`", server.config.Version)
	}

	// unsupported version:
	if err := server.SetVersion("1.20.7"); err != ErrVersionUnsupported {
		t.Fatalf("expected error `%v`, got `%v`", ErrVersionUnsupported, err)
	}
}
