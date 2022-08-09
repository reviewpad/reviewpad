package engine

import (
	"fmt"
	"regexp"
	"strings"
)

// normalize function normalizes *ReviewpadFile with default values, validates and modifies property values.
// If no customNormalizeRules specified, normalize uses default normalization rules
// specified in defaultNormalizers
// general usage: err := normalize(reviewpadFile)
// normalize is not critical to errors. So if the error occurs while normalizing some property,
// it just adds the error to the error list, and assigns the default value to the property.
// It's up to you how to react on errors, but the idea is to normalize *ReviewpadFile with default values,
// and to log errors if needed
func normalize(f *ReviewpadFile, customNormalizeRules ...map[propertyKey]*NormalizeRule) error {
	rules := defaultNormalizers
	if len(customNormalizeRules) > 0 {
		rules = customNormalizeRules[0]
	}

	// create a slice for the errors
	var errStrings []string

	for key, rule := range rules {
		var err error
		switch key {
		case _version:
			f.Version, err = rule.Do(f.Version)
		case _edition:
			f.Edition, err = rule.Do(f.Edition)
		case _mode:
			f.Mode, err = rule.Do(f.Mode)
		}
		if err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}

	if errStrings != nil {
		return fmt.Errorf(strings.Join(errStrings, "\n"))
	}
	return nil
}

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

// defaultNormalizers contains all we need to properly normalize property values by default
var defaultNormalizers = map[propertyKey]*NormalizeRule{
	_version: NewNormalizeRule(defaultApiVersion).
		WithValidators(func(val string) error {
			apiVersReg := regexp.MustCompile(`^reviewpad\.com/v[0-3]\.[\dx]$`)
			if apiVersReg.MatchString(val) {
				return nil
			}
			return fmt.Errorf("incorrect api-version: %s", val)
		}),
	_edition: NewNormalizeRule(defaultEdition, allowedEditions...),
	_mode:    NewNormalizeRule(defaultMode, allowedModes...),
}

// validator & modificator functions signature
type (
	validator   func(val string) error
	modificator func(val string) string
)

// NormalizeRule normalizes property values
type NormalizeRule struct {
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

// NewNormalizeRule creates NormalizeRule with default value as first argument
// and with the list of allowed values if provided.
// Two modificators to trim spaces and to lower case are added by default.
func NewNormalizeRule(defaultVal string, allowed ...string) *NormalizeRule {
	n := &NormalizeRule{Default: defaultVal, Allowed: allowed}
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
func (n *NormalizeRule) Do(val string) (string, error) {
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

// WithValidators adds validator functions to the NormalizeRule
func (n *NormalizeRule) WithValidators(v ...validator) *NormalizeRule {
	if n.Validators == nil {
		n.Validators = make([]validator, 0, len(v))
	}
	n.Validators = append(n.Validators, v...)
	return n
}

// WithModificators adds modificator functions to the NormalizeRule
func (n *NormalizeRule) WithModificators(m ...modificator) *NormalizeRule {
	if n.Modificators == nil {
		n.Modificators = make([]modificator, 0, len(m))
	}
	n.Modificators = append(n.Modificators, m...)
	return n
}
