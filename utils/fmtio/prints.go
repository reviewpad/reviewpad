// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package fmtio

import "fmt"

func Sprintf(context string, format string, a ...interface{}) string {
	return fmt.Sprintf("[%v] %v", context, fmt.Sprintf(format, a...))
}

func Sprint(context string, val string) string {
	return fmt.Sprintf("[%v] %v", context, fmt.Sprint(val))
}
