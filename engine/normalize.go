package engine

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	defaultApiVersion = "reviewpad.com/v3.x"
	defaultMode       = "silent"
	defaultEdition    = "professional"
)

type (
	validator   func(val *ReviewpadFile) error
	modificator func(file *ReviewpadFile) (*ReviewpadFile, error)
)

type NormalizeRule struct {
	Validators   []validator
	Modificators []modificator
}

func defaultEditionNormalizer() *NormalizeRule {
	normalizedRule := NewNormalizeRule()
	normalizedRule.WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Edition == "" {
			file.Edition = defaultEdition
			return file, nil
		}

		file.Edition = strings.ToLower(strings.TrimSpace(file.Edition))
		return file, nil
	})
	return normalizedRule
}

func defaultModeNormalizer() *NormalizeRule {
	defaultModeNormalizer := NewNormalizeRule()
	defaultModeNormalizer.WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Mode == "" {
			file.Mode = defaultMode
			return file, nil
		}

		file.Mode = strings.ToLower(strings.TrimSpace(file.Mode))
		return file, nil
	})
	return defaultModeNormalizer
}

func defaultVersionNormalizer() *NormalizeRule {
	apiVersReg := regexp.MustCompile(`^reviewpad\.com\/v\d+\.(\d+\.\d+|x)$`)
	defaultVersionNormalizer := NewNormalizeRule()
	defaultVersionNormalizer.WithModificators(func(file *ReviewpadFile) (*ReviewpadFile, error) {
		if file.Version == "" {
			file.Version = defaultApiVersion
			return file, nil
		}

		file.Version = strings.ToLower(strings.TrimSpace(file.Version))
		return file, nil
	})
	defaultVersionNormalizer.WithValidators(func(file *ReviewpadFile) error {
		if apiVersReg.MatchString(file.Version) {
			return nil
		}
		return fmt.Errorf("incorrect api-version: %s", file.Version)
	})
	return defaultVersionNormalizer
}

func normalize(f *ReviewpadFile, customRules ...*NormalizeRule) (*ReviewpadFile, error) {
	var err error
	var errStrings []string
	rules := []*NormalizeRule{
		defaultEditionNormalizer(),
		defaultModeNormalizer(),
		defaultVersionNormalizer(),
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
		for _, modificator := range n.Modificators {
			file, err = modificator(file)
			if err != nil {
				return nil, err
			}
		}
	}

	if n.Validators != nil {
		for _, validator := range n.Validators {
			if err := validator(file); err != nil {
				return file, fmt.Errorf("normalize.Do validation: %w", err)
			}
		}
		return file, nil
	}

	return file, nil
}

func (n *NormalizeRule) WithValidators(v ...validator) {
	if n.Validators == nil {
		n.Validators = make([]validator, 0, len(v))
	}
	n.Validators = append(n.Validators, v...)
}

func (n *NormalizeRule) WithModificators(m ...modificator) {
	if n.Modificators == nil {
		n.Modificators = make([]modificator, 0, len(m))
	}
	n.Modificators = append(n.Modificators, m...)
}
