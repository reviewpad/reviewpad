package engine

import (
	"fmt"
	"regexp"
	"strings"
)

// default properties values
const (
	defaultApiVersion = "reviewpad.com/v3.x"
	defaultMode       = "silent"
	defaultEdition    = "professional"
)

var (
	apiVersReg               = regexp.MustCompile(`^reviewpad\.com/v[0-3]\.[\dx]$`)
	defaultEditionNormalizer = NewNormalizeRule().WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Edition == "" {
			file.Edition = defaultEdition
			return file, nil
		}

		file.Edition = strings.ToLower(strings.TrimSpace(file.Edition))
		return file, nil
	})
	defaultModeNormalizer = NewNormalizeRule().WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Mode == "" {
			file.Mode = defaultMode
			return file, nil
		}

		file.Mode = strings.ToLower(strings.TrimSpace(file.Mode))
		return file, nil
	})
	defaultVersionNormalizer = NewNormalizeRule().WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Version == "" {
			file.Version = defaultApiVersion
			return file, nil
		}

		file.Version = strings.ToLower(strings.TrimSpace(file.Version))
		return file, nil
	}).WithValidators(func(file *ReviewpadFile) error {
		if apiVersReg.MatchString(file.Version) {
			return nil
		}
		return fmt.Errorf("incorrect api-version: %s", file.Version)
	})
)

// validator & modificator functions signature
type (
	validator   func(val *ReviewpadFile) error
	modificator func(file *ReviewpadFile) (*ReviewpadFile, error)
)

// NormalizeRule normalizes property values
type NormalizeRule struct {
	// Validators slice of functions to control property value
	// could be checks, regexp etc
	Validators []validator

	// Modificators modify reviewpad file before validate and other controls
	Modificators []modificator
}

// normalize function normalizes *ReviewpadFile with default values, validates and modifies property values.
func normalize(f *ReviewpadFile, customRules ...*NormalizeRule) (*ReviewpadFile, error) {
	var err error
	var errStrings []string
	rules := []*NormalizeRule{
		defaultEditionNormalizer,
		defaultModeNormalizer,
		defaultVersionNormalizer,
	}
	rules = append(rules, customRules...)

	for _, rule := range rules {
		f, err = rule.Do(f)
		if err != nil {
			errStrings = append(errStrings, err.Error())
		}
	}

	if errStrings != nil {
		return nil, fmt.Errorf(strings.Join(errStrings, "\n"))
	}

	return f, nil
}

func NewNormalizeRule() *NormalizeRule {
	return &NormalizeRule{}
}

func (n *NormalizeRule) Do(file *ReviewpadFile) (*ReviewpadFile, error) {
	var err error

	if n.Modificators != nil {
		for _, modiF := range n.Modificators {
			file, err = modiF(file)
			if err != nil {
				return nil, err
			}
		}
	}

	if n.Validators != nil {
		for _, valFunc := range n.Validators {
			if err := valFunc(file); err != nil {
				return file, fmt.Errorf("normalize.Do validation: %w", err)
			}
		}
		return file, nil
	}

	return file, nil
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
