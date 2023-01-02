// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package utils

import "os"

func ReadFile(filepath string) ([]byte, error) {
	fileData, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return fileData, nil
}
