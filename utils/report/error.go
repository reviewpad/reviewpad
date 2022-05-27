// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file

package report

import "fmt"

func Error(format string, a ...interface{}) string {
	return fmt.Sprintf("Error occurred! Details:\n%v", fmt.Sprintf(format, a...))
}
