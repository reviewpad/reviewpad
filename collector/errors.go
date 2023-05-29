// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package collector

// errorsToIgnore is a list of errors that should be ignored when executing reviewpad.
var errorsToIgnore = []string{
	// Error related with the failing of the summarize built-in.
	// This is usally caused by the fact that the system is overloaded.
	"summarize failed on request",
	// Errors related with bad use of a built-in.
	"failed to build AST on input",
	"type inference failed",
}
