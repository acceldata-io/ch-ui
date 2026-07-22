package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLicenseFromEnv(t *testing.T) {
	dir := t.TempDir()
	licPath := filepath.Join(dir, "license.json")
	if err := os.WriteFile(licPath, []byte("  {\"from\":\"file\"}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	t.Run("unset", func(t *testing.T) {
		t.Setenv("CHUI_LICENSE_FILE", "")
		t.Setenv("CHUI_LICENSE", "")
		if lic, _ := LicenseFromEnv(); lic != "" {
			t.Fatalf("expected empty, got %q", lic)
		}
	})

	t.Run("inline env", func(t *testing.T) {
		t.Setenv("CHUI_LICENSE_FILE", "")
		t.Setenv("CHUI_LICENSE", ` {"from":"env"} `)
		lic, src := LicenseFromEnv()
		if lic != `{"from":"env"}` || src != "CHUI_LICENSE" {
			t.Fatalf("got %q from %q", lic, src)
		}
	})

	t.Run("file wins over inline", func(t *testing.T) {
		t.Setenv("CHUI_LICENSE_FILE", licPath)
		t.Setenv("CHUI_LICENSE", `{"from":"env"}`)
		lic, src := LicenseFromEnv()
		if lic != `{"from":"file"}` {
			t.Fatalf("got %q from %q", lic, src)
		}
	})

	t.Run("unreadable file falls back to inline", func(t *testing.T) {
		t.Setenv("CHUI_LICENSE_FILE", filepath.Join(dir, "missing.json"))
		t.Setenv("CHUI_LICENSE", `{"from":"env"}`)
		lic, _ := LicenseFromEnv()
		if lic != `{"from":"env"}` {
			t.Fatalf("got %q", lic)
		}
	})

	t.Run("empty file falls back to inline", func(t *testing.T) {
		empty := filepath.Join(dir, "empty.json")
		if err := os.WriteFile(empty, []byte("  \n"), 0o600); err != nil {
			t.Fatal(err)
		}
		t.Setenv("CHUI_LICENSE_FILE", empty)
		t.Setenv("CHUI_LICENSE", `{"from":"env"}`)
		lic, _ := LicenseFromEnv()
		if lic != `{"from":"env"}` {
			t.Fatalf("got %q", lic)
		}
	})
}
