package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("ERDTREE_URL", "http://erdtree:8001")
	t.Setenv("GLADOS_URL", "http://glados:8002")
	erdtreeURL := os.Getenv("ERDTREE_URL")
	gladosURL := os.Getenv("GLADOS_URL")

	conf, err := LoadConfig()

	if err != nil {
		t.Fatalf("expected error to be nil, received: %v", err)
	}

	if conf.ErdtreeURL != erdtreeURL {
		t.Errorf("expected: %s, received: %s", erdtreeURL, conf.ErdtreeURL)
	}

	if conf.GladosURL != gladosURL {
		t.Errorf("expected: %s, received: %s", gladosURL, conf.GladosURL)
	}
}
