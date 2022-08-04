package engine

import (
	"fmt"
	"regexp"
	"strings"
)

// default properties values
const (
	defaultApiVersion = "reviewpad.com/v3.x"
	defaultEdition    = "professional"
	defaultMode       = "silent"
)

// allowed properties values
var (
	allowedEditions = []string{"professional", "team"} //[]string{defaultEdition, "team"}
	allowedModes    = []string{"silent", "verbose"}
)

type propertyKey string

const (
	_version propertyKey = "Version"
	_edition propertyKey = "Edition"
	_mode    propertyKey = "Mode"
)

// normalizers contains all we need to properly normalize property values
var normalizers = map[propertyKey]*Normalize{
	_version: NewNormalize(defaultApiVersion).
		WithValidators(func(val string) error {
			apiVersReg := regexp.MustCompile(`^reviewpad\.com/v[0-3]\.[\dx]$`)
			if apiVersReg.MatchString(val) {
				return nil
			}
			return fmt.Errorf("incorrect api-version: %s", val)
		}),
	_edition: NewNormalize(defaultEdition, allowedEditions...),
	_mode:    NewNormalize(defaultMode, allowedModes...),
}

// validator & modificator functions signature
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

// NewNormalize creates Normalize with default value as first argument
// and with the list of allowed values if provided.
// Two modificators to trim spaces and to lower case are added by default.
func NewNormalize(defaultVal string, allowed ...string) *Normalize {
	n := &Normalize{Default: defaultVal, Allowed: allowed}
	n.WithModificators(strings.TrimSpace, strings.ToLower)
	return n
}

// Do is the main method for property value normalization.
// 1) if Modificators available - modifies value.
// 2) if empty value - return default value
// 3) if Validators available - try to validate value:
//    - if validation fails returns default value and validation error
//    - if validation passes returns value and nil error
// 4) if Allowed list of values provided, checks if value in this list.
//    - if ok return value and nil error
//    - if not ok return default value and not allowed error
// 5) returns value (if: not empty, nil Validators, nil Allowed). Returned value is modified with Modificators.
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

// WithValidators adds validator functions to the Normalize
func (n *Normalize) WithValidators(v ...validator) *Normalize {
	if n.Validators == nil {
		n.Validators = make([]validator, 0, len(v))
	}
	n.Validators = append(n.Validators, v...)
	return n
}

// WithModificators adds modificator functions to the Normalize
func (n *Normalize) WithModificators(m ...modificator) *Normalize {
	if n.Modificators == nil {
		n.Modificators = make([]modificator, 0, len(m))
	}
	n.Modificators = append(n.Modificators, m...)
	return n
}
