package engine

import (
	"gopkg.in/yaml.v3"
)

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

	r.Version, _ = normalizers[_defaultApiVersion].Do(interim.Version)
	r.Edition, _ = normalizers[_defaultEdition].Do(interim.Edition)
	r.Mode, _ = normalizers[_defaultMode].Do(interim.Mode)
	r.IgnoreErrors = interim.IgnoreErrors

	_ = interim.Imports.Decode(&r.Imports)
	_ = interim.Groups.Decode(&r.Groups)
	_ = interim.Rules.Decode(&r.Rules)
	_ = interim.Labels.Decode(&r.Labels)
	_ = interim.Workflows.Decode(&r.Workflows)

	return nil
}
