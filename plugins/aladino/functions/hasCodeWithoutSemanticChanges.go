// Copyright (C) 2023 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"strings"

	doublestar "github.com/bmatcuk/doublestar/v4"
	"github.com/google/uuid"
	"github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v3/codehost"
	"github.com/reviewpad/reviewpad/v3/codehost/github/target"
	"github.com/reviewpad/reviewpad/v3/handler"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	semantic "github.com/reviewpad/reviewpad/v3/plugins/aladino/semantic"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v3/plugins/aladino/services"
	"github.com/reviewpad/reviewpad/v3/utils"
	"google.golang.org/grpc/metadata"
)

// RequestIDKey identifies request id field in context
const RequestIDKey = "request-id"

func HasCodeWithoutSemanticChanges() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildArrayOfType(aladino.BuildStringType())}, aladino.BuildBoolType()),
		Code:           hasCodeWithoutSemanticChanges,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasCodeWithoutSemanticChanges(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)
	head := pullRequest.PullRequest.GetHead()
	url := head.GetRepo().GetURL()
	log := e.GetLogger()

	// filter patch by file pattern
	patch := pullRequest.Patch
	newPatch := make(map[string]*codehost.File)

	filePatterns := args[0].(*aladino.ArrayValue)
	for fp, file := range patch {
		ignore := false
		for _, v := range filePatterns.Vals {
			filePatternRegex := v.(*aladino.StringValue)

			re, err := doublestar.Match(filePatternRegex.Val, fp)
			if err == nil && re {
				ignore = true
				break
			}
		}
		if !ignore {
			newPatch[fp] = file
			log.Infof("file %s is not ignored", fp)
		} else {
			log.Infof("file %s is ignored", fp)
		}
	}

	baseSymbols, baseFiles, err := semantic.GetSymbolsFromBaseByPatch(e, newPatch)
	if err != nil {
		return nil, err
	}

	headSymbols, headFiles, err := semantic.GetSymbolsFromHeadByPatch(e, newPatch)
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

	req := &api.GetDefinedSymbolsDiffRequest{
		RepoUri: url,
		Old:     baseSymbols,
		New:     headSymbols,
		GitDiff: gitDiff,
	}

	requestID := uuid.New().String()
	ctx := e.GetCtx()
	md := metadata.Pairs(RequestIDKey, requestID)
	reqCtx := metadata.NewOutgoingContext(ctx, md)

	reply, err := diffClient.GetDefinedSymbolsDiff(reqCtx, req)
	if err != nil {
		return nil, err
	}

	// if any symbol is modified, then return false
	res := true
	for _, symbol := range reply.Symbols {
		// log symbol information
		log.WithField("symbol", symbol).Info("semantic diff symbol")
		var sb strings.Builder

		sb.WriteString(fmt.Sprintf("path: %v, ", symbol.MainDefinitionPath))

		if symbol.OldInfo != nil {
			sb.WriteString(fmt.Sprintf("old: %s, ", symbol.OldInfo.Name))
		}

		if symbol.NewInfo != nil {
			sb.WriteString(fmt.Sprintf("old: %s, ", symbol.NewInfo.Name))
		}

		sb.WriteString(fmt.Sprintf("change type: %s", symbol.ChangeType))
		log.Infof(sb.String())

		log.Infof("------------------------------------------------")

		if symbol.ChangeType != entities.ChangeType_UNMODIFIED {
			res = false
		}
	}

	if res {
		log.Infof("no semantic changes detected")
		return aladino.BuildTrueValue(), nil
	} else {
		log.Infof("semantic changes detected")
		return aladino.BuildFalseValue(), nil
	}
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
			old := &entities.GitDiffSpan{}
			if diffBlock.Old != nil {
				old.Start = diffBlock.Old.Start
				old.End = diffBlock.Old.End
			}

			new := &entities.GitDiffSpan{}
			if diffBlock.New != nil {
				new.Start = diffBlock.New.Start
				new.End = diffBlock.New.End
			}

			gitFileDiffs[i] = &entities.GitDiffBlock{
				Old:       old,
				New:       new,
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
