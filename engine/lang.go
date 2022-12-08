// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import "github.com/reviewpad/reviewpad/v3/handler"

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
	Name        string `yaml:"name"`
	Kind        string `yaml:"kind"`
	Description string `yaml:"description"`
	Spec        string `yaml:"spec"`
}

func (p PadRule) equals(o PadRule) bool {
	if p.Name != o.Name {
		return false
	}

	if p.Kind != o.Kind {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	if p.Spec != o.Spec {
		return false
	}

	return true
}

var kinds = []string{"patch", "author"}

type PadWorkflowRule struct {
	Rule         string   `yaml:"rule"`
	ExtraActions []string `yaml:"extra-actions" mapstructure:"extra-actions"`
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
	Name        string `yaml:"name"`
	Color       string `yaml:"color"`
	Description string `yaml:"description"`
}

func (p PadLabel) equals(o PadLabel) bool {
	if p.Name != o.Name {
		return false
	}

	if p.Color != o.Color {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	return true
}

type PadWorkflow struct {
	Name               string                     `yaml:"name"`
	On                 []handler.TargetEntityKind `yaml:"on"`
	Description        string                     `yaml:"description"`
	AlwaysRun          bool                       `yaml:"always-run"`
	Rules              []PadWorkflowRule          `yaml:"-"`
	Actions            []string                   `yaml:"then"`
	NonNormalizedRules []interface{}              `yaml:"if"`
}

func (p PadWorkflow) equals(o PadWorkflow) bool {
	if p.Name != o.Name {
		return false
	}

	if len(p.On) != len(o.On) {
		return false
	}

	for i, pO := range p.On {
		oO := o.On[i]
		if pO != oO {
			return false
		}
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
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Kind        string `yaml:"kind"`
	Type        string `yaml:"type"`
	Spec        string `yaml:"spec"`
	Param       string `yaml:"param"`
	Where       string `yaml:"where"`
}

func (p PadGroup) equals(o PadGroup) bool {
	if p.Name != o.Name {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	if p.Kind != o.Kind {
		return false
	}

	if p.Type != o.Type {
		return false
	}

	if p.Spec != o.Spec {
		return false
	}

	if p.Where != o.Where {
		return false
	}

	return true
}

type ReviewpadFile struct {
	Version      string              `yaml:"api-version"`
	Edition      string              `yaml:"edition"`
	Mode         string              `yaml:"mode"`
	IgnoreErrors bool                `yaml:"ignore-errors"`
	Imports      []PadImport         `yaml:"imports"`
	Extends      []string            `yaml:"extends"`
	Groups       []PadGroup          `yaml:"groups"`
	Rules        []PadRule           `yaml:"rules"`
	Labels       map[string]PadLabel `yaml:"labels"`
	Workflows    []PadWorkflow       `yaml:"workflows"`
	Pipelines    []PadPipeline       `yaml:"pipelines"`
}

type PadPipeline struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Trigger     string     `yaml:"trigger"`
	Stages      []PadStage `yaml:"stages"`
}

type PadStage struct {
	Actions []string `yaml:"actions"`
	Until   string   `yaml:"until"`
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

	if r.IgnoreErrors != o.IgnoreErrors {
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

	if len(r.Extends) != len(o.Extends) {
		return false
	}
	for i, rE := range r.Extends {
		oE := o.Extends[i]
		if rE != oE {
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
		r.Rules = make([]PadRule, 0)
	}

	r.Rules = append(r.Rules, o.Rules...)
}

func (r *ReviewpadFile) appendGroups(o *ReviewpadFile) {
	if r.Groups == nil {
		r.Groups = make([]PadGroup, 0)
	}

	r.Groups = append(r.Groups, o.Groups...)
}

func (r *ReviewpadFile) appendWorkflows(o *ReviewpadFile) {
	if r.Workflows == nil {
		r.Workflows = make([]PadWorkflow, 0)
	}

	r.Workflows = append(r.Workflows, o.Workflows...)
}

func (r *ReviewpadFile) appendPipelines(o *ReviewpadFile) {
	if r.Pipelines == nil {
		r.Pipelines = make([]PadPipeline, 0)
	}

	r.Pipelines = append(r.Pipelines, o.Pipelines...)
}

func (r *ReviewpadFile) extends(o *ReviewpadFile) {
	// extend labels
	if r.Labels == nil {
		r.Labels = make(map[string]PadLabel)
	}

	for labelName, label := range o.Labels {
		if _, ok := r.Labels[labelName]; !ok {
			r.Labels[labelName] = label
		}
	}

	// extend groups
	if r.Groups == nil {
		r.Groups = make([]PadGroup, 0)
	}

	filteredGroups := make([]PadGroup, 0)

	for _, group := range o.Groups {
		if _, ok := findGroup(r.Groups, group.Name); !ok {
			filteredGroups = append(filteredGroups, group)
		}
	}

	r.Groups = append(filteredGroups, r.Groups...)

	// extend rules
	if r.Rules == nil {
		r.Rules = make([]PadRule, 0)
	}

	filteredRules := make([]PadRule, 0)

	for _, rule := range o.Rules {
		if _, ok := findRule(r.Rules, rule.Name); !ok {
			filteredRules = append(filteredRules, rule)
		}
	}

	r.Rules = append(filteredRules, r.Rules...)

	// extend workflows
	if r.Workflows == nil {
		r.Workflows = make([]PadWorkflow, 0)
	}

	filteredWorkflows := make([]PadWorkflow, 0)

	for _, workflow := range o.Workflows {
		if _, ok := findWorkflow(r.Workflows, workflow.Name); !ok {
			filteredWorkflows = append(filteredWorkflows, workflow)
		}
	}

	r.Workflows = append(filteredWorkflows, r.Workflows...)

	// extend pipelines
	if r.Pipelines == nil {
		r.Pipelines = make([]PadPipeline, 0)
	}

	filteredPipelines := make([]PadPipeline, 0)

	for _, pipeline := range o.Pipelines {
		if _, ok := findPipeline(r.Pipelines, pipeline.Name); !ok {
			filteredPipelines = append(filteredPipelines, pipeline)
		}
	}

	r.Pipelines = append(filteredPipelines, r.Pipelines...)
}

func findGroup(groups []PadGroup, name string) (*PadGroup, bool) {
	for _, group := range groups {
		if group.Name == name {
			return &group, true
		}
	}

	return nil, false
}

func findRule(rules []PadRule, name string) (*PadRule, bool) {
	for _, rule := range rules {
		if rule.Name == name {
			return &rule, true
		}
	}

	return nil, false
}

func findWorkflow(workflows []PadWorkflow, name string) (*PadWorkflow, bool) {
	for _, workflow := range workflows {
		if workflow.Name == name {
			return &workflow, true
		}
	}

	return nil, false
}

func findPipeline(pipelines []PadPipeline, name string) (*PadPipeline, bool) {
	for _, pipeline := range pipelines {
		if pipeline.Name == name {
			return &pipeline, true
		}
	}

	return nil, false
}
