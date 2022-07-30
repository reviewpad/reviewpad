// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllMatches(t *testing.T) {
	wantRuleNames := []string{
		"rule1",
		"rule2",
		"rule3",
		"rule4",
		"rule-five",
	}

	groups := []PadGroup{}
	rules := []PadRule{
		{
			Spec: `$rule("rule1") && $rule("rule2")`,
		},
		{
			Spec: `true`,
		},
		{
			Spec: `$rule("rule3") && $test()`,
		},
		{
			Spec: `$test() && $rule("rule4")`,
		},
		{
			Spec: `$test1()`,
		},
		{
			Spec: `$rule("rule-five") && $isElementOf($author(), $group("group-name"))`,
		},
	}
	workflows := []PadWorkflow{}

	gotRuleNames := getAllRuleCalls(groups, rules, workflows)

	assert.Equal(t, wantRuleNames, gotRuleNames)
}
