package plugins_aladino_actions

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/reviewpad/reviewpad/v3/lang/aladino"
	"github.com/reviewpad/reviewpad/v3/utils"
	"github.com/shurcooL/githubv4"
)

func AssignReviewerWithPolicy() *aladino.BuiltInAction {
	return &aladino.BuiltInAction{
		Type: aladino.BuildFunctionType([]aladino.Type{
			aladino.BuildStringType(),
			aladino.BuildArrayOfType(aladino.BuildStringType()),
			aladino.BuildIntType(),
		}, nil),
		Code: assignReviewerWithPolicyCode,
	}
}

func assignReviewerWithPolicyCode(e aladino.Env, args []aladino.Value) error {
	policy := args[0].(*aladino.StringValue)
	availableReviewers := args[1].(*aladino.ArrayValue).Vals
	requiredReviewers := args[2].(*aladino.IntValue).Val
	return assignReviewersWithPolicy(e, policy.Val, availableReviewers, requiredReviewers)
}

func assignReviewersWithPolicy(e aladino.Env, policy string, availableReviewers []aladino.Value, requiredReviewers int) error {
	if len(availableReviewers) == 0 {
		return fmt.Errorf("assignReviewer: list of reviewers can't be empty")
	}
	if requiredReviewers == 0 {
		return fmt.Errorf("assignReviewer: required reviewers can't be 0")
	}

	// Remove pull request author from provided reviewers list
	for index, reviewer := range availableReviewers {
		if reviewer.(*aladino.StringValue).Val == *e.GetPullRequest().User.Login {
			availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
			break
		}
	}

	totalAvailableReviewers := len(availableReviewers)
	if requiredReviewers > totalAvailableReviewers {
		log.Printf("assignReviewer: more reviewers required (%d) than available (%d)", requiredReviewers, totalAvailableReviewers)
		requiredReviewers = totalAvailableReviewers
	}

	prNum := utils.GetPullRequestNumber(e.GetPullRequest())
	owner := utils.GetPullRequestBaseOwnerName(e.GetPullRequest())
	repo := utils.GetPullRequestBaseRepoName(e.GetPullRequest())

	reviewers := []string{}

	reviews, _, err := e.GetClient().PullRequests.ListReviews(e.GetCtx(), owner, repo, prNum, nil)
	if err != nil {
		return err
	}

	// Re-request current reviewers if mention on the provided reviewers list
	for _, review := range reviews {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == *review.User.Login {
				requiredReviewers--
				if *review.State != "APPROVED" {
					reviewers = append(reviewers, *review.User.Login)
				}
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	// Skip current requested reviewers if mention on the provided reviewers list
	currentRequestedReviewers := e.GetPullRequest().RequestedReviewers
	for _, requestedReviewer := range currentRequestedReviewers {
		for index, availableReviewer := range availableReviewers {
			if availableReviewer.(*aladino.StringValue).Val == *requestedReviewer.Login {
				requiredReviewers--
				availableReviewers = append(availableReviewers[:index], availableReviewers[index+1:]...)
				break
			}
		}
	}

	switch policy {
	case "random":
		// Select random reviewers from the list of all provided reviewers
		for i := 0; i < requiredReviewers; i++ {
			selectedElementIndex := utils.GenerateRandom(len(availableReviewers))

			selectedReviewer := availableReviewers[selectedElementIndex]
			availableReviewers = append(availableReviewers[:selectedElementIndex], availableReviewers[selectedElementIndex+1:]...)

			reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
		}
	case "round-robin":
		// using the issue number and the number of required reviewers,
		// determine reviewers
		startPos := ((prNum - 1) * requiredReviewers)
		if startPos < 0 {
			startPos = 0
		}
		for i := 0; i < requiredReviewers; i++ {
			pos := (startPos + i) % len(availableReviewers)
			selectedReviewer := availableReviewers[pos]
			reviewers = append(reviewers, selectedReviewer.(*aladino.StringValue).Val)
		}
	case "reviewpad":
		// TODO: determine reviewers with reviewerScores
	}

	if len(reviewers) == 0 {
		log.Printf("assignReviewer: skipping request reviewers. the pull request already has reviewers")
		return nil
	}

	_, _, err = e.GetClient().PullRequests.RequestReviewers(e.GetCtx(), owner, repo, prNum, github.ReviewersRequest{
		Reviewers: reviewers,
	})

	return err
}

type githubUserStatus struct {
	Emoji                        *string
	Message                      *string
	IndicatesLimitedAvailability bool
}

func (s *githubUserStatus) isHoliday() bool {
	messageLower := strings.ToLower(derefString(s.Message))
	return s.IndicatesLimitedAvailability ||
		derefString(s.Emoji) == ":palm_tree:" ||
		strings.Contains(messageLower, "vacation") ||
		strings.Contains(messageLower, "holiday")
}

// query ($owner: String!, $repo: String!) {
//   repository(owner: $owner, name:$repo) {
//     collaborators {
//       nodes {
//         login
//         status {
//           emoji
//           message
//           indicatesLimitedAvailability
//         }
//       }
//     }
//   }
// }
type repositoryCollaboratorsWithStatusesQuery struct {
	Repository struct {
		Collaborators struct {
			Nodes []struct {
				Login  *string
				Status *githubUserStatus
			}
		}
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

// query ($login: String!) {
//   user(login: $login) {
//     status {
//       emoji
//       message
//       indicatesLimitedAvailability
//     }
//   }
// }
type userStatusQuery struct {
	User struct {
		Status *githubUserStatus
	} `graphql:"user(login: $login)"`
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func reviewerScores(e aladino.Env, owner, repo string, availableReviewers []string) error {
	// Determine in a "smart" way which reviewers are more available to perform the review,
	// and which are more adequate to review the specific pull request.
	// Parameters we have:
	// - workload: sum of the logsize (= log10(size) of PR, capped at 5)
	//             of all the user's current PRs in the repo he's a reviewer for.
	// - holiday: is the user on holiday? (1 or 0)
	// - blame: what part of the changed files in the current PR did the user write? (between 0 and 1)
	// these come together to form the score of a user. Once the score of
	// all potential reviewers is determined, the reviewers are sorted in
	// ascending order and the first n are set as the reviewers.
	// Thus:
	// (holiday * 100000) + (workload/(0.8+blame/5)
	// users on holiday fall at the bottom of the list almost always,
	// and the ranking is determined by the user's workload, which is
	// increased by up to 25% if the user is unfamiliar with the code of the PR.

	// We fetch concurrently the holiday status, workload and blame of each user.
	// The goroutines report to their respective channels,
	// and can report back critical errors with the critical channels.
	// Every channel operation must be guarded against ctx.Done
	type holidayStatus struct {
		login   string
		holiday bool
	}

	holidayStatuses := make(chan holidayStatus, 8)
	critical := make(chan error)
	ctx, canc := context.WithCancel(e.GetCtx())
	defer canc()

	go func() {
		// We first fetch the repo collaborators, as the availableReviewers
		// is likely a subset of the collaborators.
		var queryResult repositoryCollaboratorsWithStatusesQuery
		err := e.GetClientGQL().Query(ctx, &queryResult, map[string]interface{}{
			"owner": githubv4.String(owner),
			"repo":  githubv4.String(repo),
		})
		if err != nil {
			select {
			case <-ctx.Done():
			case critical <- fmt.Errorf("error fetching repository collaborators with status: %w", err):
			}
			return
		}

		collaborators := queryResult.Repository.Collaborators.Nodes
	AvailableReviewersLoop:
		for _, reviewer := range availableReviewers {
			// Check whether we have already the reviewer in the fetched collaborators,
			// and in that case, use its status to report back on whether they are on
			// holiday.
			reviewerLower := strings.ToLower(reviewer)
			for _, collaborator := range collaborators {
				if strings.ToLower(*collaborator.Login) == reviewerLower {
					status := collaborator.Status
					select {
					case holidayStatuses <- holidayStatus{
						login:   reviewer,
						holiday: status.isHoliday(),
					}:
					case <-ctx.Done():
						return
					}
					continue AvailableReviewersLoop
				}
			}

			// The collaborator was not found in the fetched collaborators.
			// Fetch individual user status.
			var userStatus userStatusQuery
			err = e.GetClientGQL().Query(ctx, &userStatus, map[string]interface{}{
				"login": githubv4.String(reviewer),
			})
			if err != nil {
				// HACK(Morgan B): ideally we'd use the "type" field of an
				// error object returned by the graphql API, however this is
				// not supported by the graphql library we're using, and on
				// top of it, a fix for it has been added in a PR which has
				// been open for the last 4 years at the time of writing:
				// https://github.com/shurcooL/graphql/pull/33
				// The solution thus would be to change from using the
				// current githubv4 library to a generic graphql library,
				// which would allow us also to generate a custom query
				// to let us query multiple users at the same time using
				// aliases instead of having to repeat this query for every
				// user not found in the collaborators.
				if strings.Contains(err.Error(), "Could not resolve to a User with the login of") {
					continue
				}
				select {
				case <-ctx.Done():
				case critical <- fmt.Errorf("error fetching user status: %w", err):
				}
			}
		}
		close(holidayStatuses)
	}()

	prs, err := utils.GetPullRequests(e.GetCtx(), e.GetClient(), owner, repo)
	if err != nil {
		return fmt.Errorf("assignReviewerWithPolicy: error getting pull requests: %w", err)
	}

	type reviewerParams struct {
		workload float64
		holiday  bool
		blame    float64
	}

	_ = prs
	return nil
}
