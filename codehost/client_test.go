// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package codehost

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHostInfo(t *testing.T) {
	hostInfo, err := GetHostInfo("https://github.com")
	assert.Nil(t, err)
	assert.Equal(t, &HostInfo{
		HostUri: "https://github.com",
	}, hostInfo)
}
