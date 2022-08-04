package engine

import (
	"gopkg.in/yaml.v3"
)

// UnmarshalYAML our custom implementation of unmarshaller yaml interface.
// On this stage we can inject normalizers for property values without affecting existing codebase.
func (r *ReviewpadFile) UnmarshalYAML(value *yaml.Node) error {
	var interim struct {
		Version      string    `yaml:"api-version"`
		Edition      string    `yaml:"edition"`
		Mode         string    `yaml:"mode"`
		IgnoreErrors bool      `yaml:"ignore-errors"`
		Imports      yaml.Node `yaml:"imports"`
		Groups       yaml.Node `yaml:"groups"`
		Rules        yaml.Node `yaml:"rules"`
		Labels       yaml.Node `yaml:"labels"`    //map[string]PadLabel
		Workflows    yaml.Node `yaml:"workflows"` //[]PadWorkflow
	}
	if err := value.Decode(&interim); err != nil {
		return err
	}

	// we need only default property values if incorrect or empty for these fields.
	// so the errors are ignored here
	r.Version, _ = normalizers[_version].Do(interim.Version)
	r.Edition, _ = normalizers[_edition].Do(interim.Edition)
	r.Mode, _ = normalizers[_mode].Do(interim.Mode)
	r.IgnoreErrors = interim.IgnoreErrors

	_ = interim.Imports.Decode(&r.Imports)
	_ = interim.Groups.Decode(&r.Groups)
	_ = interim.Rules.Decode(&r.Rules)
	_ = interim.Labels.Decode(&r.Labels)
	_ = interim.Workflows.Decode(&r.Workflows)

	return nil
}
