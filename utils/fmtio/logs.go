// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package fmtio

import (
	"fmt"
	"log"
)

func LogPrintln(context string, format string, a ...interface{}) {
	log.Printf("[%v] %v", context, fmt.Sprintf(format, a...))
}
