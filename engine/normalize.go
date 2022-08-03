package engine

import (
	"fmt"
	"regexp"
	"strings"
)

// defaults
const (
	_defaultApiVersion = "reviewpad.com/v3.x"
	_defaultEdition    = "professional"
	_defaultMode       = "silent"
)

// allowed values
var (
	_allowedEditions = []string{"professional", "team"} //[]string{_defaultEdition, "team"}
	_allowedModes    = []string{"silent", "verbose"}
)

var normalizers = map[string]*Normalize{
	_defaultApiVersion: NewNormalize(_defaultApiVersion).
		WithValidators(func(val string) error {
			apiVersReg := regexp.MustCompile(`^reviewpad\.com/v[0-3]\.[\dx]$`)
			if apiVersReg.MatchString(val) {
				return nil
			}
			return fmt.Errorf("incorrect api-version: %s", val)
		}),
	_defaultEdition: NewNormalize(_defaultEdition, _allowedEditions...),
	_defaultMode:    NewNormalize(_defaultMode, _allowedModes...),
}

type (
	validator   func(val string) error
	modificator func(val string) string
)

// Normalize normalizes property values
type Normalize struct {
	// Default contains default property value
	Default string

	// Allowed contains allowed property values
	Allowed []string

	// Validators slice of functions to control property value
	// could be checks, regexp etc
	Validators []validator

	// Modificators modify property value before validate and other controls
	Modificators []modificator
}

func NewNormalize(defaultVal string, allowed ...string) *Normalize {
	n := &Normalize{Default: defaultVal, Allowed: allowed}
	n.WithModificators(strings.TrimSpace, strings.ToLower)
	return n
}

func (n *Normalize) Do(val string) (string, error) {
	if n.Modificators != nil {
		for _, modiF := range n.Modificators {
			val = modiF(val)
		}
	}

	if val == "" {
		return n.Default, nil
	}

	if n.Validators != nil {
		for _, valFunc := range n.Validators {
			if err := valFunc(val); err != nil {
				return n.Default, fmt.Errorf("normalize.Do validation: %w", err)
			}
		}
		return val, nil
	}

	if n.Allowed != nil {
		for _, a := range n.Allowed {
			if val == a {
				return val, nil
			}
		}
		return n.Default, fmt.Errorf("normalize.Do not in allowed list: %s", val)
	}

	return val, nil
}

func (n *Normalize) WithValidators(v ...validator) *Normalize {
	if n.Validators == nil {
		n.Validators = make([]validator, 0, len(v))
	}
	n.Validators = append(n.Validators, v...)
	return n
}

func (n *Normalize) WithModificators(m ...modificator) *Normalize {
	if n.Modificators == nil {
		n.Modificators = make([]modificator, 0, len(m))
	}
	n.Modificators = append(n.Modificators, m...)
	return n
}
