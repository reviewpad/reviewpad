// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package engine

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateLabelColor(t *testing.T) {
	tests := map[string]struct {
		arg     *PadLabel
		wantErr error
	}{
		"hex value with #": {
			arg:     &PadLabel{Color: "#91cf60"},
			wantErr: nil,
		},
		"hex value without #": {
			arg:     &PadLabel{Color: "91cf60"},
			wantErr: nil,
		},
		"invalid hex value with #": {
			arg:     &PadLabel{Color: "#91cg60"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value without #": {
			arg:     &PadLabel{Color: "91cg60"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value because of size with #": {
			arg:     &PadLabel{Color: "#91cg6"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"invalid hex value because of size without #": {
			arg:     &PadLabel{Color: "91cg6"},
			wantErr: execError("evalLabel: color code not valid"),
		},
		"english color": {
			arg:     &PadLabel{Color: "red"},
			wantErr: execError("evalLabel: color code not valid"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			gotErr := validateLabelColor(test.arg)
			assert.Equal(t, test.wantErr, gotErr)
		})
	}
}
