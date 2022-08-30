// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file
package utils

import "strings"

func MinifyQuery(query string) string {
	minifiedQuery := strings.ReplaceAll(query, "\n", "")
	minifiedQuery = strings.ReplaceAll(minifiedQuery, "\t", "")
	minifiedQuery = strings.ReplaceAll(minifiedQuery, " ", "")
	minifiedQuery += "\n"
	return minifiedQuery
}
