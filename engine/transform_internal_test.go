// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformAladinoExpression(t *testing.T) {
	tests := map[string]struct {
		arg     string
		wantVal string
	}{
		"id": {
			arg:     "$id()",
			wantVal: "$id()",
		},
		"merge": {
			arg:     "$merge()",
			wantVal: "$merge(\"merge\")",
		},
		"size": {
			arg:     "$size()",
			wantVal: "$size([])",
		},
		"issueCountBy simple": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0",
		},
		"issueCountBy": {
			arg:     "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$issueCountBy(\"john\", \"open\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"pullRequestCountBy simple": {
			arg:     "$pullRequestCountBy(\"john\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"pullRequestCountBy nil state": {
			arg:     "$pullRequestCountBy(\"john\", \"\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0",
		},
		"pullRequestCountBy nil dev": {
			arg:     "$pullRequestCountBy(\"\", \"closed\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"closed\") > 0",
		},
		"pullRequestCountBy nil dev and state": {
			arg:     "$pullRequestCountBy(\"\") > 0",
			wantVal: "$pullRequestCountBy(\"\", \"all\") > 0",
		},
		"pullRequestCountBy and issueCountBy": {
			arg:     "$pullRequestCountBy(\"john\") > 0 && true && $issueCountBy(\"dev\") > 0",
			wantVal: "$pullRequestCountBy(\"john\", \"all\") > 0 && true && $issueCountBy(\"dev\", \"all\") > 0",
		},
		"pullRequestCountBy with dev": {
			arg:     `$pullRequestCountBy("john") > 0 && true && $pullRequestCountBy("john", "open") > 1`,
			wantVal: `$pullRequestCountBy("john", "all") > 0 && true && $pullRequestCountBy("john", "open") > 1`,
		},
		"pullRequestCountBy with function call": {
			arg:     `$pullRequestCountBy($author()) > 1 && $pullRequestCountBy($author(), "open") > 2`,
			wantVal: `$pullRequestCountBy($author(), "all") > 1 && $pullRequestCountBy($author(), "open") > 2`,
		},
		"pullRequestCountBy with function call and state": {
			arg:     `$pullRequestCountBy($author(), "open")`,
			wantVal: `$pullRequestCountBy($author(), "open")`,
		},
		"pullRequestCountBy with function call and empty state": {
			arg:     `$pullRequestCountBy($author(), "")`,
			wantVal: `$pullRequestCountBy($author(), "all")`,
		},
		"close": {
			arg:     "$close()",
			wantVal: `$close("", "completed")`,
		},
		"close with comment": {
			arg:     `$close("comment")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with comment and not_planned state reason": {
			arg:     `$close("comment", "not_planned")`,
			wantVal: `$close("comment", "not_planned")`,
		},
		"close with comment and completed state reason": {
			arg:     `$close("comment", "completed")`,
			wantVal: `$close("comment", "completed")`,
		},
		"close with empty comment and completed state reason": {
			arg:     `$close("", "completed")`,
			wantVal: `$close("", "completed")`,
		},
		"close with empty comment and not_planned state reason": {
			arg:     `$close("", "not_planned")`,
			wantVal: `$close("", "not_planned")`,
		},
		"multiple $haveAllChecksRunCompleted with one empty": {
			arg:     `$haveAllChecksRunCompleted() && true && $haveAllChecksRunCompleted(["build"])`,
			wantVal: `$haveAllChecksRunCompleted([], "") && true && $haveAllChecksRunCompleted(["build"], "")`,
		},
		"one $haveAllChecksRunCompleted with arg provided": {
			arg:     `$haveAllChecksRunCompleted(["build", "test"], "")`,
			wantVal: `$haveAllChecksRunCompleted(["build", "test"], "")`,
		},
		"one $haveAllChecksRunCompleted with no arg provided": {
			arg:     `$haveAllChecksRunCompleted()`,
			wantVal: `$haveAllChecksRunCompleted([], "")`,
		},
		"one $haveAllChecksRunCompleted with no empty arg provided": {
			arg:     `$haveAllChecksRunCompleted([], "success")`,
			wantVal: `$haveAllChecksRunCompleted([], "success")`,
		},
		"multiple $haveAllChecksRunCompleted": {
			arg:     `$haveAllChecksRunCompleted() && true && $haveAllChecksRunCompleted(["build"]) && $haveAllChecksRunCompleted(["build", "run"], "completed") && $haveAllChecksRunCompleted([]) && $haveAllChecksRunCompleted([], "")`,
			wantVal: `$haveAllChecksRunCompleted([], "") && true && $haveAllChecksRunCompleted(["build"], "") && $haveAllChecksRunCompleted(["build", "run"], "completed") && $haveAllChecksRunCompleted([], "") && $haveAllChecksRunCompleted([], "")`,
		},
		// TODO: test addDefaultTotalRequestedReviewers
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotVal := transformAladinoExpression(test.arg)
			assert.Equal(t, test.wantVal, gotVal)
		})
	}
}
