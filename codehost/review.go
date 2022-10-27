package codehost

func HasReview(reviews []*Review, userLogin string) bool {
	for _, review := range reviews {
		if review.User.Login == userLogin {
			return true
		}
	}

	return false
}

func LastReview(reviews []*Review, userLogin string) *Review {
	var latestReview *Review
	for _, review := range reviews {
		if review.User.Login == userLogin {
			if latestReview == nil || latestReview.SubmittedAt.Before(*review.SubmittedAt) {
				latestReview = review
			}
		}
	}

	return latestReview
}
