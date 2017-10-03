package main

import (
	"log"
	"os"
	"testing"
)

func TestLoadCfgDir(t *testing.T) {
	c, err := loadCfgDir("./conf.d")

	if err != nil {
		t.Errorf("Expected to load files from configdir! got: %#v", err)
	}
	if len(c) != 4 {
		t.Errorf("Expected 4 checks loaded! got: %d", len(c))
	}
}

func TestLoadCfgDirFail(t *testing.T) {

	_, err := loadCfgDir("./non_existant_conf.d")

	if err == nil {
		t.Error("Expected Error to load non existant directory")
	}
}

func TestLoadCfgDirNotReadableFile(t *testing.T) {
	if err := os.Chmod("conf.d/two_sched.cfg", 0100); err != nil {
		log.Fatal(err)
	}
	_, err := loadCfgDir("./conf.d")
	if err == nil {
		t.Error("Expected Error to load non readable file in directory")
	}
	if err := os.Chmod("conf.d/two_sched.cfg", 0664); err != nil {
		log.Fatal(err)
	}
}

func TestLoadCfg(t *testing.T) {
	c, err := loadCfg("./conf.d/one_sched.cfg")

	if err != nil {
		t.Errorf("Expected no error loading cfg file! got: %s", err)
	}
	if len(c) != 3 {
		t.Errorf("Expected 3 checks loaded! got: %d", len(c))
	}
	if c[1].Command != "echo parameter eins" {
		t.Errorf("Expect second command to be 'command mit parametern', got %s", c[1].Command)
	}
}

func TestLoadCfgFail(t *testing.T) {
	_, err := loadCfg("./non_existant_conf.d")

	if err == nil {
		t.Error("Expected : Error to load non existant directory")
	}
}

func TestCheckNames(t *testing.T) {
	c, _ := loadCfg("./conf.d/one_sched.cfg")

	if c[0].Name != "one_sched" {
		t.Errorf("Expected first checks Name to be one_sched! got: %#v", c[0].Name)
	}
}

func TestCheckType(t *testing.T) {
	c, _ := loadCfg("./conf.d/one_sched.cfg")

	if c[0].Comment.Type != "" {
		t.Errorf("Expected first checks Type to be empty! got: %#v", c[2].Comment.Type)
	}
	if c[1].Comment.Type != "gauge" {
		t.Errorf("Expected second checks Type to be gauge! got: %#v", c[2].Comment.Type)
	}
}

func TestParseLineWithIntervalError(t *testing.T) {
	c := NewCheck()
	_, err := c.parseLine([]string{"#", "INTERVAL", "b"})
	if err == nil {
		t.Errorf("Expected Error with wrong INTERVAL settings! got: %s", err)
	}
}
