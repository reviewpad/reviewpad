// Copyright 2022 Explore.dev Unipessoal Lda. All Rights Reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package fmtio_test

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/reviewpad/reviewpad/v2/utils/fmtio"
	"github.com/stretchr/testify/assert"
)

func TestLogPrintLn(t *testing.T) {
	var buf bytes.Buffer
    log.SetOutput(&buf)

   	fmtio.LogPrintln("test", "test value: %v", "TEST")

	wantLog := fmt.Sprintf("%v [test] test value: TEST\n", time.Now().Format("2006/01/02 15:04:05"))
	gotLog := buf.String()

	assert.Equal(t, wantLog, gotLog)
}