// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package utils

import (
	"math/rand"
	"time"
)

func GenerateRandom(size int) int {
	seedNow := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(seedNow)

	return randGen.Intn(size)
}

func AbsInt32(value int32) int32 {
	if value >= 0 {
		return value
	}
	return -value
}
