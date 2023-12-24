package validators

import "project/internal/model"

func RemoveDuplicateContributors(contributors []model.ContributorWithStatus) []model.ContributorWithStatus {
	seen := make(map[string]bool)
	result := []model.ContributorWithStatus{}

	for _, contributor := range contributors {
		// Assuming Email is the field you want to check for duplicates
		email := contributor.Email

		// Check if the email has been seen before
		if !seen[email] {
			result = append(result, contributor)
			// Mark the email as seen
			seen[email] = true
		}
	}

	return result
}
