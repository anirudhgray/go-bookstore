package recommender

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
)

// GetUserLikesDislikes retrieves the book IDs that a user has liked and disliked based on their review ratings.
// It queries the database for reviews associated with the given user ID and categorizes the books into liked and disliked lists.
func GetUserLikesDislikes(userID uint) (liked []uint, notLiked []uint, err error) {
	var reviews []models.Review
	result := database.DB.Where("user_id = ?", userID).Find(&reviews)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	var likedBooks []uint
	var notLikedBooks []uint

	for _, review := range reviews {
		if review.Rating >= 4 {
			likedBooks = append(likedBooks, review.BookID)
		} else {
			notLikedBooks = append(notLikedBooks, review.BookID)
		}
	}

	return likedBooks, notLikedBooks, nil
}

// GetUsersWhoReviewedBook retrieves the user IDs of users who have reviewed a specific book.
// It queries the database for reviews associated with the given book ID and returns the user IDs of users who have reviewed the book.
func GetUsersWhoReviewedBook(bookID uint) (userslist []uint, err error) {
	var reviews []models.Review
	result := database.DB.Where("book_id = ?", bookID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	var users []uint

	for _, review := range reviews {
		users = append(users, review.UserID)
	}

	return users, nil
}

// GetUsersWithSimilarInteractions retrieves user IDs who have interacted with books that are similar to the given user's liked and disliked books.
// It queries the database for users who have reviewed books that match the liked and disliked books of the given user.
func GetUsersWithSimilarInteractions(likedBooks []uint, dislikedBooks []uint, userID uint) (userslist []uint, err error) {
	var users []uint

	for _, bookID := range likedBooks {
		reviewedUsers, err := GetUsersWhoReviewedBook(bookID)
		if err != nil {
			return nil, err
		}
		users = append(users, reviewedUsers...)
	}

	for _, bookID := range dislikedBooks {
		reviewedUsers, err := GetUsersWhoReviewedBook(bookID)
		if err != nil {
			return nil, err
		}
		users = append(users, reviewedUsers...)
	}

	// remove duplicates and remove the current user's ID
	currentUserID := userID
	uniqueUsers := make(map[uint]bool)
	var resultUsers []uint

	for _, userID := range users {
		if userID != currentUserID && !uniqueUsers[userID] {
			uniqueUsers[userID] = true
			resultUsers = append(resultUsers, userID)
		}
	}

	return resultUsers, nil
}

func formatFloat(num float64, prc int) string {
	var (
		zero, dot = "0", "."

		str = fmt.Sprintf("%."+strconv.Itoa(prc)+"f", num)
	)

	return strings.TrimRight(strings.TrimRight(str, zero), dot)
}

// CalculateUserSimilarity of two users via modified Jaccard Coefficient.
//
//	S(U1, U2) = (|L1 intersec L2| + |D1 intersect D2| - |L1 intersect D2| - |L2 intersect D1|) / |L1 union L2 union D1 union D2|
func CalculateUserSimilarity(currentUserLiked, currentUserDisliked, otherUserLiked, otherUserDisliked []uint) (similarityCoefficient float64) {
	L1, L2, D1, D2 := currentUserLiked, otherUserLiked, currentUserDisliked, otherUserDisliked

	L1IntersectL2Size := intersectionSize(currentUserLiked, otherUserLiked)
	D1IntersectD2Size := intersectionSize(currentUserDisliked, otherUserDisliked)
	L1IntersectD2Size := intersectionSize(currentUserLiked, otherUserDisliked)
	L2IntersectD1Size := intersectionSize(otherUserLiked, currentUserDisliked)

	numerator := float64(L1IntersectL2Size + D1IntersectD2Size - L1IntersectD2Size - L2IntersectD1Size)
	denominator := float64(unionSize4(L1, L2, D1, D2))

	if numerator == 0 {
		return 0
	}

	similarity := numerator / denominator

	res, _ := strconv.ParseFloat(formatFloat(similarity, 3), 64)

	return res
}

func intersectionSize(set1, set2 []uint) int {
	count := 0
	set2Map := make(map[uint]bool)
	for _, item := range set2 {
		set2Map[item] = true
	}
	for _, item := range set1 {
		if set2Map[item] {
			count++
		}
	}
	return count
}

func unionSize4(arr1, arr2, arr3, arr4 []uint) int {
	unionSet := make(map[uint]bool)

	for _, item := range arr1 {
		unionSet[item] = true
	}
	for _, item := range arr2 {
		unionSet[item] = true
	}
	for _, item := range arr3 {
		unionSet[item] = true
	}
	for _, item := range arr4 {
		unionSet[item] = true
	}

	return len(unionSet)
}

// CalculateSimilaritiesWithOtherUsers returns a mapping of similar users (returned by GetUsersWithSimilarInteractions) to their similarity coefficient with current user.
func CalculateSimilaritiesWithOtherUsers(currentUserID uint, similarUsers []uint, currentUserLiked, currentUserDisliked []uint) (sims map[uint]float64, err error) {
	userSimilarities := make(map[uint]float64)

	for _, otherUserID := range similarUsers {
		otherUserLiked, otherUserDisliked, err := GetUserLikesDislikes(otherUserID)
		if err != nil {
			return nil, err
		}

		similarity := CalculateUserSimilarity(currentUserLiked, currentUserDisliked, otherUserLiked, otherUserDisliked)
		userSimilarities[otherUserID] = similarity
	}

	return userSimilarities, nil
}

// GetUnreviewedBooks fetches books which the current user has not yet reviewed (will select recommendations from among these)
func GetUnreviewedBooks(currentUserID uint) ([]uint, error) {
	userLikedBooks, userDislikedBooks, err := GetUserLikesDislikes(currentUserID)
	if err != nil {
		return nil, err
	}

	userReviewedBookIDs := append(userLikedBooks, userDislikedBooks...)

	var unreviewedBooks []uint

	var allBookIDs []uint
	result := database.DB.Model(&models.Book{}).Pluck("id", &allBookIDs)
	if result.Error != nil {
		return nil, result.Error
	}

	// Filter out the books that the user has reviewed
	for _, bookID := range allBookIDs {
		if !contains(userReviewedBookIDs, bookID) {
			unreviewedBooks = append(unreviewedBooks, bookID)
		}
	}

	return unreviewedBooks, nil
}

func contains(arr []uint, item uint) bool {
	for _, val := range arr {
		if val == item {
			return true
		}
	}
	return false
}

// CalculateRecommendationProbabilities returns mapping on user's unreviewed books to the "probability" of them liking it (i.e., order of recommendation).
//
// ZL = sum of similarity coefficients of other similar users who have liked a particular book.
//
// ML = number of such users as above
//
// ZD = sum of similarity coefficients of other similar users who have disliked a particular book.
//
// MD = number of such users as above
//
//	P(U, B) = (ZL - ZD) / (|ML| + |MD|)
//
// Produces a value between -1 and 1
func CalculateRecommendationProbabilities(currentUserID uint, unreviewedBooks []uint, similarUsers map[uint]float64) map[uint]float64 {
	recommendationProbabilities := make(map[uint]float64)

	for _, bookID := range unreviewedBooks {
		ZL := 0.0
		ZD := 0.0
		ML := 0
		MD := 0

		likedByUsers, dislikedByUsers, err := GetLikersDislikersForBook(bookID)
		if err != nil {
			continue
		}

		for _, userID := range likedByUsers {
			if similarity, ok := similarUsers[userID]; ok {
				ZL += similarity
				ML++
			}
		}

		for _, userID := range dislikedByUsers {
			if similarity, ok := similarUsers[userID]; ok {
				ZD += similarity
				MD++
			}
		}

		if ML+MD == 0 {
			continue
		}

		probability := (ZL - ZD) / float64(ML+MD)
		recommendationProbabilities[bookID] = probability
	}

	return recommendationProbabilities
}

func GetLikersDislikersForBook(bookID uint) ([]uint, []uint, error) {
	var likedByUsers []uint
	var dislikedByUsers []uint

	usersWhoLiked, err := GetUsersWhoLikedBook(bookID)
	if err != nil {
		return nil, nil, err
	}
	likedByUsers = append(likedByUsers, usersWhoLiked...)

	// Fetch users who disliked the book
	usersWhoDisliked, err := GetUsersWhoDislikedBook(bookID)
	if err != nil {
		return nil, nil, err
	}
	dislikedByUsers = append(dislikedByUsers, usersWhoDisliked...)

	return likedByUsers, dislikedByUsers, nil
}

func GetUsersWhoLikedBook(bookID uint) ([]uint, error) {
	var likedByUsers []uint

	var reviews []models.Review
	result := database.DB.Where("book_id = ? AND rating >= 4", bookID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, review := range reviews {
		likedByUsers = append(likedByUsers, review.UserID)
	}

	return likedByUsers, nil
}

func GetUsersWhoDislikedBook(bookID uint) ([]uint, error) {
	var dislikedByUsers []uint

	var reviews []models.Review
	result := database.DB.Where("book_id = ? AND rating < 4", bookID).Find(&reviews)
	if result.Error != nil {
		return nil, result.Error
	}

	for _, review := range reviews {
		dislikedByUsers = append(dislikedByUsers, review.UserID)
	}

	return dislikedByUsers, nil
}

type Recommendation struct {
	Book        models.SafeBook
	Probability float64
}

func GetRecommendedBooksSortedAndPaginated(recommendationProbabilities map[uint]float64, page int, perPage int) []Recommendation {
	var recommendations []Recommendation

	for bookID, probability := range recommendationProbabilities {
		var book models.SafeBook
		database.DB.Model(&models.Book{}).First(&book, bookID)
		recommendations = append(recommendations, Recommendation{
			Book:        book,
			Probability: probability,
		})
	}

	// Sort recommendations by probability in descending order
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Probability > recommendations[j].Probability
	})

	// Paginate the recommendations
	startIdx := (page - 1) * perPage
	endIdx := int(math.Min(float64(startIdx+perPage), float64(len(recommendations))))

	if startIdx < len(recommendations) {
		return recommendations[startIdx:endIdx]
	}
	return []Recommendation{}
}
