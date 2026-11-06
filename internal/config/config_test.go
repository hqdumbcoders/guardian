package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	t.Setenv("ERDTREE_URL", "http://erdtree:8001")
	t.Setenv("GLADOS_URL", "http://glados:8002")

	conf, err := LoadConfig()

	if conf.ErdtreeURL != os.Getenv("ERDTREE_URL") {
		t.Errorf("expected: %s, received: %s", os.Getenv("ERDTREE_URL"), conf.ErdtreeURL)
		if err != "ERDTREE_URL is empty" {
			t.Errorf("expected: %s, received: %s", "ERDTREE_URL is empty", err)
		}
	}

	if conf.GladosURL != os.Getenv("GLADOS_URL") {
		t.Errorf("expected: %s, received: %s", os.Getenv("GLADOS_URL"), conf.GladosURL)
		if err != "GLADOS_URL is empty" {
			t.Errorf("expected: %s, received: %s", "GLADOS_URL is empty", err)
		}
	}

	if conf.ErdtreeURL == os.Getenv("ERDTREE_URL") && conf.GladosURL == os.Getenv("GLADOS_URL") {
		if err != nil {
			t.Errorf("expected: %s, received: %s", nil, err)
		}
	}

}
