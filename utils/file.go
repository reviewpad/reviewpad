// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import "strings"

func FileExt(fp string) string {
	strs := strings.Split(fp, ".")
	str := strings.Join(strs[1:], ".")

	if str != "" {
		str = "." + str
	}

	return str
}
