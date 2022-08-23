// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// Proprietary and confidential

package handler

import (
	"fmt"
	"log"
)

func Log(format string, a ...interface{}) {
	log.Printf("[host-event-handler] %v", fmt.Sprintf(format, a...))
}
