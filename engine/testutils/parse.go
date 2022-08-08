// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package testutils

import (
	"github.com/reviewpad/reviewpad/v3/engine"
	"gopkg.in/yaml.v3"
)

func ParseReviewpadFile(data []byte) (*engine.ReviewpadFile, error) {
	reviewpadFile := &engine.ReviewpadFile{}
	err := yaml.Unmarshal(data, &reviewpadFile)
	if err != nil {
		return nil, err
	}

	// At the end of loading all imports from the file, its imports are reset to []engine.PadImport{}.
	// However, the parsing of the wanted reviewpad file, sets the imports to []engine.PadImport(nil).
	reviewpadFile.Imports = []engine.PadImport{}

	return reviewpadFile, nil
}
