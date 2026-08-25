package validation

import (
	"flag"
)

// Limits holds configuration limits for Loki.
type Limits struct {
	DiscoverLogLevels bool `yaml:"discover_log_levels" json:"discover_log_levels"`
}

// RegisterFlags adds the flags required to config this to the given FlagSet
func (l *Limits) RegisterFlags(f *flag.FlagSet) {
	l.RegisterFlagsWithPrefix("", f)
}

// RegisterFlagsWithPrefix adds the flags required to config this to the given FlagSet with a prefix
func (l *Limits) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
	f.BoolVar(&l.DiscoverLogLevels, prefix+"discover-log-levels", false, "Flag to enable automatic discovery of log levels.")
}

// Overrides manages limits per tenant.
type Overrides struct {
	defaultLimits Limits
	tenantLimits  map[string]*Limits
}

// NewOverrides creates a new Overrides manager.
func NewOverrides(defaultLimits Limits, tenantLimits map[string]*Limits) *Overrides {
	return &Overrides{
		defaultLimits: defaultLimits,
		tenantLimits:  tenantLimits,
	}
}

// DiscoverLogLevels returns whether log level auto-discovery is enabled for a tenant or globally.
func (o *Overrides) DiscoverLogLevels(userID string) bool {
	if o == nil {
		return false
	}
	if o.tenantLimits != nil {
		if l, ok := o.tenantLimits[userID]; ok && l != nil {
			return l.DiscoverLogLevels
		}
	}
	return o.defaultLimits.DiscoverLogLevels
}
