// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

const (
	PROFESSIONAL_EDITION string = "professional"
	TEAM_EDITION         string = "team"
	SILENT_MODE          string = "silent"
	VERBOSE_MODE         string = "verbose"
)

type PadImport struct {
	Url string `yaml:"url"`
}

func (p PadImport) equals(o PadImport) bool {
	return p.Url == o.Url
}

type PadRule struct {
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Spec        string `yaml:"spec"`
}

func (p PadRule) equals(o PadRule) bool {
	return p.Kind == o.Kind &&
		p.Description == o.Description &&
		p.Spec == o.Spec
}

var kinds = []string{"patch", "author"}

type PadWorkflowRule struct {
	Rule         string   `yaml:"rule"`
	ExtraActions []string `yaml:"extra-actions"`
}

func (p PadWorkflowRule) equals(o PadWorkflowRule) bool {
	if p.Rule != o.Rule {
		return false
	}

	if len(p.ExtraActions) != len(o.ExtraActions) {
		return false
	}
	for i, pE := range p.ExtraActions {
		oE := o.ExtraActions[i]
		if pE != oE {
			return false
		}
	}

	return true
}

type PadLabel struct {
	Color       string `yaml:"color"`
	Description string `yaml:"description"`
}

func (p PadLabel) equals(o PadLabel) bool {
	if p.Color != o.Color {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	return true
}

type PadWorkflow struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	AlwaysRun   bool              `yaml:"always-run"`
	Rules       []PadWorkflowRule `yaml:"if"`
	Actions     []string          `yaml:"then"`
}

func (p PadWorkflow) equals(o PadWorkflow) bool {
	if p.Name != o.Name {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	if len(p.Rules) != len(o.Rules) {
		return false
	}

	for i, pP := range p.Rules {
		oP := o.Rules[i]
		if !pP.equals(oP) {
			return false
		}
	}

	if len(p.Actions) != len(o.Actions) {
		return false
	}

	if p.AlwaysRun != o.AlwaysRun {
		return false
	}

	for i, pA := range p.Actions {
		oA := o.Actions[i]
		if pA != oA {
			return false
		}
	}

	return true
}

type PadGroup struct {
	Description string `yaml:"description"`
	Kind        string `yaml:"kind"`
	Type        string `yaml:"type"`
	Spec        string `yaml:"spec"`
	Param       string `yaml:"param"`
	Where       string `yaml:"where"`
}

func (p PadGroup) equals(o PadGroup) bool {
	return p.Description == o.Description &&
		p.Kind == o.Kind &&
		p.Type == o.Type &&
		p.Spec == o.Spec &&
		p.Where == o.Where
}

type ReviewpadFile struct {
	Version   string              `yaml:"api-version"`
	Edition   string              `yaml:"edition"`
	Mode      string              `yaml:"mode"`
	Imports   []PadImport         `yaml:"imports"`
	Groups    map[string]PadGroup `yaml:"groups"`
	Rules     map[string]PadRule  `yaml:"rules"`
	Labels    map[string]PadLabel `yaml:"labels"`
	Workflows []PadWorkflow       `yaml:"workflows"`
}

func (r *ReviewpadFile) equals(o *ReviewpadFile) bool {
	if r.Version != o.Version {
		return false
	}

	if r.Edition != o.Edition {
		return false
	}

	if r.Mode != o.Mode {
		return false
	}

	if len(r.Imports) != len(o.Imports) {
		return false
	}
	for i, rI := range r.Imports {
		oI := o.Imports[i]
		if !rI.equals(oI) {
			return false
		}
	}

	if len(r.Rules) != len(o.Rules) {
		return false
	}
	for i, rR := range r.Rules {
		oR := o.Rules[i]
		if !rR.equals(oR) {
			return false
		}
	}

	if len(r.Labels) != len(o.Labels) {
		return false
	}

	for i, rL := range r.Labels {
		oL := o.Labels[i]
		if !rL.equals(oL) {
			return false
		}
	}

	if len(r.Workflows) != len(o.Workflows) {
		return false
	}
	for i, rP := range r.Workflows {
		oP := o.Workflows[i]
		if !rP.equals(oP) {
			return false
		}
	}

	if len(r.Groups) != len(o.Groups) {
		return false
	}
	for i, rG := range r.Groups {
		oG := o.Groups[i]
		if !rG.equals(oG) {
			return false
		}
	}

	return true
}

func (r *ReviewpadFile) appendLabels(o *ReviewpadFile) {
	if r.Labels == nil {
		r.Labels = make(map[string]PadLabel)
	}

	for labelName, label := range o.Labels {
		r.Labels[labelName] = label
	}
}

func (r *ReviewpadFile) appendRules(o *ReviewpadFile) {
	if r.Rules == nil {
		r.Rules = make(map[string]PadRule)
	}

	for ruleName, rule := range o.Rules {
		r.Rules[ruleName] = rule
	}
}

func (r *ReviewpadFile) appendWorkflows(o *ReviewpadFile) {
	if r.Workflows == nil {
		r.Workflows = make([]PadWorkflow, 0)
	}

	r.Workflows = append(r.Workflows, o.Workflows...)
}
