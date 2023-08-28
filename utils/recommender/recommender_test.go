package recommender

import (
	"testing"
)

func TestCalculateUserSimilarity(t *testing.T) {
	testCases := []struct {
		name                string
		currentUserLiked    []uint
		currentUserDisliked []uint
		otherUserLiked      []uint
		otherUserDisliked   []uint
		expectedSimilarity  float64
		expectedTolerance   float64
	}{
		{
			name:                "Case 1: No Interactions",
			currentUserLiked:    []uint{},
			currentUserDisliked: []uint{},
			otherUserLiked:      []uint{},
			otherUserDisliked:   []uint{},
			expectedSimilarity:  0.0,
			expectedTolerance:   1e-6,
		},
		{
			name:                "Case 2: Perfect Match",
			currentUserLiked:    []uint{1, 2, 3},
			currentUserDisliked: []uint{4, 5},
			otherUserLiked:      []uint{1, 2, 3},
			otherUserDisliked:   []uint{4, 5},
			expectedSimilarity:  1.0,
			expectedTolerance:   1e-6,
		},
		{
			name:                "Case 3: Partial Overlap",
			currentUserLiked:    []uint{1, 2, 3},
			currentUserDisliked: []uint{4, 5},
			otherUserLiked:      []uint{3, 4, 5, 6},
			otherUserDisliked:   []uint{1},
			expectedSimilarity:  -0.333,
			expectedTolerance:   1e-6,
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualSimilarity := CalculateUserSimilarity(tc.currentUserLiked, tc.currentUserDisliked, tc.otherUserLiked, tc.otherUserDisliked)
			if !(tc.expectedSimilarity-tc.expectedTolerance <= actualSimilarity && actualSimilarity <= tc.expectedSimilarity+tc.expectedTolerance) {
				t.Errorf("Expected similarity: %f, Actual similarity: %f", tc.expectedSimilarity, actualSimilarity)
			}
		})
	}
}
