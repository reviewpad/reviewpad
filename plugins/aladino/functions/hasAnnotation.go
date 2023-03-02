// Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

package plugins_aladino_functions

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/reviewpad/api/go/entities"
	api "github.com/reviewpad/api/go/services"
	"github.com/reviewpad/reviewpad/v4/codehost"
	"github.com/reviewpad/reviewpad/v4/codehost/github/target"
	"github.com/reviewpad/reviewpad/v4/handler"
	"github.com/reviewpad/reviewpad/v4/lang/aladino"
	plugins_aladino_services "github.com/reviewpad/reviewpad/v4/plugins/aladino/services"
)

func HasAnnotation() *aladino.BuiltInFunction {
	return &aladino.BuiltInFunction{
		Type:           aladino.BuildFunctionType([]aladino.Type{aladino.BuildStringType()}, aladino.BuildBoolType()),
		Code:           hasAnnotationCode,
		SupportedKinds: []handler.TargetEntityKind{handler.PullRequest},
	}
}

func hasAnnotationCode(e aladino.Env, args []aladino.Value) (aladino.Value, error) {
	annotation := args[0].(*aladino.StringValue).Val

	symbolsByFileName, err := getSymbolsFromPatch(e)
	if err != nil {
		return nil, err
	}

	for _, fileSymbols := range symbolsByFileName {
		for _, symbol := range fileSymbols.Symbols {
			for _, symbolComment := range symbol.CodeComments {
				comment := symbolComment.Code
				if commentHasAnnotation(comment, annotation) {
					return aladino.BuildTrueValue(), nil
				}
			}
		}
	}

	return aladino.BuildFalseValue(), nil
}

func commentHasAnnotation(comment, annotation string) bool {
	normalizedComment := strings.ToLower(strings.Trim(comment, " \n"))
	normalizedAnnotation := strings.ToLower(annotation)
	anPrefix := "reviewpad-an: "
	if strings.HasPrefix(normalizedComment, anPrefix) {
		rest := normalizedComment[len(anPrefix)-1:]
		restUntilNewLine := strings.Split(rest, "\n")[0]

		annotations := strings.Fields(restUntilNewLine)
		for _, annot := range annotations {
			if normalizedAnnotation == annot {
				return true
			}
		}
	}

	return false
}

func getSymbolsFromPatch(e aladino.Env) (map[string]*entities.Symbols, error) {
	pullRequest := e.GetTarget().(*target.PullRequestTarget)

	res := make(map[string]*entities.Symbols)

	head := pullRequest.PullRequest.GetHead()
	base := pullRequest.PullRequest.GetBase()
	url := head.GetRepo().GetURL()
	patch := pullRequest.Patch

	lastCommit := head.SHA

	for fp, commitFile := range patch {
		blob, err := e.GetGithubClient().DownloadContents(e.GetCtx(), fp, head)
		if err != nil {
			// If fails to download the file from head, then tries to download it from base.
			// This happens when a file has been removed.
			blob, err = e.GetGithubClient().DownloadContents(e.GetCtx(), fp, base)

			if err != nil {
				return nil, err
			}
		}

		// Convert the bytes array to string, remove the escape and convert again to bytes
		str, _ := strconv.Unquote(strings.Replace(strconv.Quote(string(blob)), `\\u`, `\u`, -1))
		blob = []byte(str)

		blocks := getBlocks(commitFile)

		service, ok := e.GetBuiltIns().Services[plugins_aladino_services.SEMANTIC_SERVICE_KEY]
		if !ok {
			return nil, fmt.Errorf("semantic service not found")
		}

		semanticClient := service.(api.SemanticClient)
		req := &api.GetSymbolsRequest{
			Uri:      url,
			CommitId: *lastCommit,
			Filepath: fp,
			Blob:     blob,
			BlobId:   *commitFile.Repr.SHA,
			Diff:     &entities.ResolveFileDiff{Blocks: blocks},
		}
		reply, err := semanticClient.GetSymbols(e.GetCtx(), req)
		if err != nil {
			return nil, err
		}

		res[fp] = reply.Symbols
	}

	return res, nil
}

func getBlocks(commitFile *codehost.File) []*entities.ResolveBlockDiff {
	res := make([]*entities.ResolveBlockDiff, len(commitFile.Diff))
	for i, block := range commitFile.Diff {
		if block.New != nil {
			res[i] = &entities.ResolveBlockDiff{
				Start: block.New.Start,
				End:   block.New.End,
			}
		} else {
			// invariant: block.Old != nil
			res[i] = &entities.ResolveBlockDiff{
				Start: block.Old.Start,
				End:   block.Old.End,
			}
		}
	}

	return res
}
