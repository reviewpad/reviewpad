// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

func ElementOf(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func Contains(slice []string, str string) bool {
	for _, elem := range slice {
		if elem == str {
			return true
		}
	}

	return false
}
