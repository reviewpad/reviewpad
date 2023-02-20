// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"strings"

	"github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	semantic "github.com/reviewpad/reviewpad/v3/plugins/aladino/semantic"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v3/plugins/aladino/services"
	"github.com/reviewpad/reviewpad/v3/utils"
)

func HasOnlyCodeChanges() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type: aladino.BuildFunctionType(
			[]aladino.Type{
				aladino.BuildArrayOfType(aladino.BuildStringType()),
				aladino.BuildArrayOfType(aladino.BuildStringType()),
			}, aladino.BuildBoolType()),
		Code:           hasOnlyCodeChanges,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasOnlyCodeChanges(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	head := pullRequest.PullRequest.GetHead()
	url := head.GetRepo().GetURL()

	baseSymbols, baseFiles, err := semantic.GetSymbolsFromBase(e)
	if err != nil {
		return nil, err
	}

	headSymbols, headFiles, err := semantic.GetSymbolsFromHead(e)
	if err != nil {
		return nil, err
	}

	gitDiff, err := getRawDiff(e, baseFiles, headFiles)
	if err != nil {
		return nil, err
	}

	service, ok := e.GetBuiltIns().Services[plugins_aladino_services.DIFF_SERVICE_KEY]
	if !ok {
		return nil, fmt.Errorf("diff service not found")
	}

	diffClient := service.(api.DiffClient)

	req := &api.EnhancedDiffRequest{
		Old:     baseSymbols,
		New:     headSymbols,
		RawDiff: gitDiff,
		RepoUri: url,
	}

	reply, err := diffClient.EnhancedDiff(e.GetCtx(), req)
	if err != nil {
		return nil, err
	}

	// if any symbol is modified, then return false
	for _, symbol := range reply.Diff.Symbols {
		if symbol.ChangeType != entities.ChangeType_UNMODIFIED {
			return aladino.BuildFalseValue(), nil
		}
	}

	return aladino.BuildTrueValue(), nil
}

func getRawDiff(e aladino.Env, baseFiles map[string]string, headFiles map[string]string) (*entities.GitDiff, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	head := pullRequest.PullRequest.GetHead()
	headCommit := head.GetSHA()

	base := pullRequest.PullRequest.GetBase()
	baseCommit := base.GetSHA()

	files := make(map[string]*entities.FileInfoDiff)
	diffs := make(map[string]*entities.GitFileDiffs)

	for fp, commitFile := range pullRequest.Patch {
		gitFileDiffs := make([]*entities.GitDiffBlock, len(commitFile.Diff))

		for i, diffBlock := range commitFile.Diff {
			gitFileDiffs[i] = &entities.GitDiffBlock{
				Old: &entities.GitDiffSpan{
					Start: diffBlock.Old.Start,
					End:   diffBlock.Old.End,
				},
				New: &entities.GitDiffSpan{
					Start: diffBlock.New.Start,
					End:   diffBlock.New.End,
				},
				IsContext: diffBlock.IsContext,
			}
		}

		diffs[fp] = &entities.GitFileDiffs{
			Blocks: gitFileDiffs,
		}

		files[fp] = &entities.FileInfoDiff{
			OldInfo: &entities.FileInfo{
				Path:      commitFile.Repr.GetPreviousFilename(),
				Extension: utils.FileExt(commitFile.Repr.GetPreviousFilename()),
				BlobId:    commitFile.Repr.GetBlobURL(),
				NumLines:  getLineCount(baseFiles[fp]),
			},
			NewInfo: &entities.FileInfo{
				Path:      commitFile.Repr.GetFilename(),
				Extension: utils.FileExt(commitFile.Repr.GetFilename()),
				BlobId:    commitFile.Repr.GetBlobURL(),
				NumLines:  getLineCount(headFiles[fp]),
			},
			NumAddedLines:   int32(commitFile.Repr.GetAdditions()),
			NumRemovedLines: int32(commitFile.Repr.GetDeletions()),
			// FIXME: this is not correct, we need to check if the file is binary
			// IsBinary: false,
		}
	}

	res := &entities.GitDiff{
		OldCommitId: baseCommit,
		NewCommitId: headCommit,
		Files:       files,
		Diffs:       diffs,
	}

	return res, nil
}

func getLineCount(contents string) int32 {
	numLines := int32(strings.Count(contents, "\n") + 1)
	//Ignoring end of line whitespaces
	if len(contents) > 0 && rune(contents[len(contents)-1]) == '\n' {
		numLines--
	}

	return numLines
}
