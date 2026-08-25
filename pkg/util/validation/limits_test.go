package validation

import (
	"flag"
	"testing"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LimitsConfig Limits `yaml:"limits_config"`
}

func TestLimits_UnmarshalYAML(t *testing.T) {
	yamlDataTrue := `
limits_config:
  discover_log_levels: true
`
	var cfgTrue Config
	err := yaml.Unmarshal([]byte(yamlDataTrue), &cfgTrue)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling yaml: %v", err)
	}
	if !cfgTrue.LimitsConfig.DiscoverLogLevels {
		t.Fatalf("expected DiscoverLogLevels to be true, got false")
	}

	yamlDataFalse := `
limits_config:
  discover_log_levels: false
`
	var cfgFalse Config
	err = yaml.Unmarshal([]byte(yamlDataFalse), &cfgFalse)
	if err != nil {
		t.Fatalf("unexpected error unmarshaling yaml: %v", err)
	}
	if cfgFalse.LimitsConfig.DiscoverLogLevels {
		t.Fatalf("expected DiscoverLogLevels to be false, got true")
	}
}

func TestLimits_RegisterFlags(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	var limits Limits
	limits.RegisterFlagsWithPrefix("limits.", fs)

	if fs.Lookup("limits.discover-log-levels") == nil {
		t.Fatal("expected flag limits.discover-log-levels to be registered")
	}
}

func TestOverrides_DiscoverLogLevels(t *testing.T) {
	defaultLimits := Limits{DiscoverLogLevels: false}
	tenantLimits := map[string]*Limits{
		"tenant-a": {DiscoverLogLevels: true},
		"tenant-b": {DiscoverLogLevels: false},
	}

	overrides := NewOverrides(defaultLimits, tenantLimits)

	if overrides.DiscoverLogLevels("tenant-default") != false {
		t.Errorf("expected tenant-default to have DiscoverLogLevels = false")
	}
	if overrides.DiscoverLogLevels("tenant-a") != true {
		t.Errorf("expected tenant-a to have DiscoverLogLevels = true")
	}
	if overrides.DiscoverLogLevels("tenant-b") != false {
		t.Errorf("expected tenant-b to have DiscoverLogLevels = false")
	}
}
