# Solution for Issue #1

## 🛠️ Proposed Solution (by Aditya Waghamare)

### Analysis
The `limits_config` parser in Loki (`pkg/util/validation/limits.go` and related validation/overrides packages) fails with an unmarshal error because `discover_log_levels` is missing from the `Limits` and `plainLimits` struct definitions, YAML tags, CLI flag registrations, and tenant override helper methods.

### Fix
Add `DiscoverLogLevels` field to struct definitions with correct YAML/JSON tags, CLI flag registration, and override evaluator method.

### Implementation
```go
// In pkg/util/validation/limits.go (or corresponding validation package)

type Limits struct {
    // ... existing fields ...
    DiscoverLogLevels bool `yaml:"discover_log_levels,omitempty" json:"discover_log_levels,omitempty"`
}

type plainLimits struct {
    // ... existing fields ...
    DiscoverLogLevels *bool `yaml:"discover_log_levels,omitempty" json:"discover_log_levels,omitempty"`
}

func (l *Limits) RegisterFlagsWithPrefix(prefix string, f *flag.FlagSet) {
    // ... existing flags ...
    f.BoolVar(&l.DiscoverLogLevels, prefix+"discover-log-levels", false, "Enable log level discovery per tenant/globally.")
}

func (l *Limits) Validate(v *validation.Validator) error {
    // ...
    return nil
}

// In pkg/util/validation/overrides.go
func (o *Overrides) DiscoverLogLevels(userID string) bool {
    if o == nil || o.Overrides == nil {
        return o.defaults.DiscoverLogLevels
    }
    val := o.FindOverride(userID, func(limits *Limits) interface{} {
        return limits.DiscoverLogLevels
    })
    if val != nil {
        return val.(bool)
    }
    return o.defaults.DiscoverLogLevels
}
```

### Testing
Verify configuration unmarshaling and startup with YAML:
```yaml
limits_config:
  discover_log_levels: true
```
And run unit tests:
```bash
go test -v ./pkg/util/validation/...
```

---
*Submitted by Aditya Waghamare*
💰 **Payout Address (Base L2 / EVM):** `0xb61dBcdBc3407F71EaCb64D4CBFAcf9FFfe2415C`