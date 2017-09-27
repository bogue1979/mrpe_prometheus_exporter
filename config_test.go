package main

import (
	"local/mrpe_node_exporter_bridge/tests"
	"testing"
)

// the go way
func TestLoadCfg(t *testing.T) {
	c, err := loadCfg("./conf.d/one_sched.cfg")
	if err != nil {
		t.Errorf("Expected no error loading cfg file! got: %s", err)
	}
	if len(c) != 3 {
		t.Errorf("Expected 3 checks loaded! got: %d", len(c))
	}
}

// simple spec
func TestCheckNames(t *testing.T) {
	spec := tests.Spec(t)
	c, _ := loadCfg("./conf.d/one_sched.cfg")
	spec.Expect(c[0].Name).ToEqual("one_sched")
	spec.Expect(c[1].Name).ToEqual("parametercheck")
	spec.Expect(c[2].Name).ToEqual("pingcheck")
}
