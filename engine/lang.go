// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"reflect"

	"github.com/reviewpad/go-lib/entities"
)

const (
	SILENT_MODE  string = "silent"
	VERBOSE_MODE string = "verbose"
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
	Name                 string                      `yaml:"name"`
	On                   []entities.TargetEntityKind `yaml:"on"`
	Description          string                      `yaml:"description"`
	AlwaysRun            bool                        `yaml:"always-run"`
	Rules                []PadWorkflowRule           `yaml:"-"`
	Actions              []string                    `yaml:"-"`
	Runs                 []PadWorkflowRunBlock       `yaml:"-"`
	NonNormalizedRules   any                         `yaml:"if"`
	NonNormalizedActions any                         `yaml:"then"`
	NonNormalizedElse    any                         `yaml:"else"`
	NonNormalizedRun     any                         `yaml:"run"`
}

type PadWorkflowRunBlock struct {
	If      []PadWorkflowRule
	Then    []PadWorkflowRunBlock
	Else    []PadWorkflowRunBlock
	Actions []string
	ForEach *PadWorkflowRunForEachBlock
}

type PadWorkflowRunForEachBlock struct {
	Key   string
	Value string
	In    string
	Do    []PadWorkflowRunBlock
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
	Mode           string              `yaml:"mode"`
	IgnoreErrors   *bool               `yaml:"ignore-errors"`
	MetricsOnMerge *bool               `yaml:"metrics-on-merge"`
	Imports        []PadImport         `yaml:"imports"`
	Extends        []string            `yaml:"extends"`
	Groups         []PadGroup          `yaml:"groups"`
	Checks         map[string]PadCheck `yaml:"checks"`
	Rules          []PadRule           `yaml:"rules"`
	Labels         map[string]PadLabel `yaml:"labels"`
	Workflows      []PadWorkflow       `yaml:"workflows"`
	Pipelines      []PadPipeline       `yaml:"pipelines"`
	Recipes        map[string]*bool    `yaml:"recipes"`
	Dictionaries   []PadDictionary     `yaml:"dictionaries"`
}

type PadCheck struct {
	Severity   string                 `yaml:"severity"`
	Activation string                 `yaml:"activation"`
	Parameters map[string]interface{} `yaml:"parameters"`
}

func (p PadCheck) equals(o PadCheck) bool {
	if p.Severity != o.Severity {
		return false
	}

	if len(p.Parameters) != len(o.Parameters) {
		return false
	}

	for key, value := range p.Parameters {
		if o.Parameters[key] != value {
			return false
		}
	}

	return true
}

type PadDictionary struct {
	Name string            `yaml:"name"`
	Spec map[string]string `yaml:"spec"`
}

func (p PadDictionary) equals(o PadDictionary) bool {
	if p.Name != o.Name {
		return false
	}

	if len(p.Spec) != len(o.Spec) {
		return false
	}

	for key, value := range p.Spec {
		if o.Spec[key] != value {
			return false
		}
	}

	return true
}

type PadPipeline struct {
	Name        string     `yaml:"name"`
	Description string     `yaml:"description"`
	Trigger     string     `yaml:"trigger"`
	Stages      []PadStage `yaml:"stages"`
}

type PadStage struct {
	Actions              []string `yaml:"-"`
	NonNormalizedActions any      `yaml:"actions"`
	Until                string   `yaml:"until"`
}

func (p PadPipeline) equals(o PadPipeline) bool {
	if p.Name != o.Name {
		return false
	}

	if p.Description != o.Description {
		return false
	}

	if p.Trigger != o.Trigger {
		return false
	}

	for i, pS := range p.Stages {
		oS := o.Stages[i]
		if !pS.equals(oS) {
			return false
		}
	}

	return true
}

func (p PadStage) equals(o PadStage) bool {
	if p.Until != o.Until {
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

func (r *ReviewpadFile) equals(o *ReviewpadFile) bool {
	if r.Mode != o.Mode {
		return false
	}

	if r.IgnoreErrors != o.IgnoreErrors {
		return false
	}

	if r.MetricsOnMerge != o.MetricsOnMerge {
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

	if len(r.Dictionaries) != len(o.Dictionaries) {
		return false
	}
	for i, rD := range r.Dictionaries {
		oD := o.Dictionaries[i]
		if !rD.equals(oD) {
			return false
		}
	}

	if len(r.Pipelines) != len(o.Pipelines) {
		return false
	}
	for i, rP := range r.Pipelines {
		oP := o.Pipelines[i]
		if !rP.equals(oP) {
			return false
		}
	}

	if len(r.Checks) != len(o.Checks) {
		return false
	}
	for i, rC := range r.Checks {
		oC := o.Checks[i]
		if !rC.equals(oC) {
			return false
		}
	}

	return reflect.DeepEqual(r.Recipes, o.Recipes)
}

func (r *ReviewpadFile) appendChecks(o *ReviewpadFile) {
	if r.Checks == nil {
		r.Checks = make(map[string]PadCheck)
	}

	for checkName, check := range o.Checks {
		r.Checks[checkName] = check
	}
}

func (r *ReviewpadFile) appendLabels(o *ReviewpadFile) {
	if r.Labels == nil {
		r.Labels = make(map[string]PadLabel)
	}

	for labelName, label := range o.Labels {
		r.Labels[labelName] = label
	}
}

func (r *ReviewpadFile) appendGroups(o *ReviewpadFile) {
	updatedGroups := make([]PadGroup, 0)

	for _, group := range r.Groups {
		if _, ok := findGroup(o.Groups, group.Name); !ok {
			updatedGroups = append(updatedGroups, group)
		}
	}

	r.Groups = append(updatedGroups, o.Groups...)
}

func (r *ReviewpadFile) appendRules(o *ReviewpadFile) {
	updatedRules := make([]PadRule, 0)

	for _, rule := range r.Rules {
		if _, ok := findRule(o.Rules, rule.Name); !ok {
			updatedRules = append(updatedRules, rule)
		}
	}

	r.Rules = append(updatedRules, o.Rules...)
}

func (r *ReviewpadFile) appendWorkflows(o *ReviewpadFile) {
	updatedWorkflows := make([]PadWorkflow, 0)

	for _, workflow := range r.Workflows {
		if _, ok := findWorkflow(o.Workflows, workflow.Name); !ok {
			updatedWorkflows = append(updatedWorkflows, workflow)
		}
	}

	r.Workflows = append(updatedWorkflows, o.Workflows...)
}

func (r *ReviewpadFile) appendPipelines(o *ReviewpadFile) {
	updatedPipelines := make([]PadPipeline, 0)

	for _, pipeline := range r.Pipelines {
		if _, ok := findPipeline(o.Pipelines, pipeline.Name); !ok {
			updatedPipelines = append(updatedPipelines, pipeline)
		}
	}

	r.Pipelines = append(updatedPipelines, o.Pipelines...)
}

func (r *ReviewpadFile) appendRecipes(o *ReviewpadFile) {
	if r.Recipes == nil {
		r.Recipes = make(map[string]*bool)
	}

	for name, active := range o.Recipes {
		r.Recipes[name] = active
	}
}

func (r *ReviewpadFile) appendDictionaries(o *ReviewpadFile) {
	updatedDictionaries := make([]PadDictionary, 0)

	for _, dictionary := range r.Dictionaries {
		if _, ok := findDictionary(o.Dictionaries, dictionary.Name); !ok {
			updatedDictionaries = append(updatedDictionaries, dictionary)
		}
	}

	r.Dictionaries = append(updatedDictionaries, o.Dictionaries...)
}

func (r *ReviewpadFile) extend(o *ReviewpadFile) {
	if o.Mode != "" {
		r.Mode = o.Mode
	}

	if o.IgnoreErrors != nil {
		r.IgnoreErrors = o.IgnoreErrors
	}

	if o.MetricsOnMerge != nil {
		r.MetricsOnMerge = o.MetricsOnMerge
	}

	r.appendChecks(o)
	r.appendLabels(o)
	r.appendGroups(o)
	r.appendRules(o)
	r.appendWorkflows(o)
	r.appendPipelines(o)
	r.appendRecipes(o)
	r.appendDictionaries(o)
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

func findDictionary(dictionaries []PadDictionary, name string) (*PadDictionary, bool) {
	for _, dictionary := range dictionaries {
		if dictionary.Name == name {
			return &dictionary, true
		}
	}

	return nil, false
}
