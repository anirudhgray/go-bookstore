package recommender

import (
	"math"
	"sort"

	"github.com/anirudhgray/balkan-assignment/infra/database"
	"github.com/anirudhgray/balkan-assignment/models"
)

func GetUserLikesDislikes(userID uint) ([]uint, []uint, error) {
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

func GetUsersWhoReviewedBook(bookID uint) ([]uint, error) {
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

func GetUsersWithSimilarInteractions(likedBooks []uint, dislikedBooks []uint, userID uint) ([]uint, error) {
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

// S(U1, U2) = (|L1 intersec L2| + |D1 intersect D2| - |L1 intersect D2| - |L2 intersect D1|) / |L1 union L2 union D1 union D2|
func CalculateUserSimilarity(currentUserLiked, currentUserDisliked, otherUserLiked, otherUserDisliked []uint) float64 {
	L1, L2, D1, D2 := currentUserLiked, otherUserLiked, currentUserDisliked, otherUserDisliked

	L1IntersectL2Size := intersectionSize(currentUserLiked, otherUserLiked)
	D1IntersectD2Size := intersectionSize(currentUserDisliked, otherUserDisliked)
	L1IntersectD2Size := intersectionSize(currentUserLiked, otherUserDisliked)
	L2IntersectD1Size := intersectionSize(otherUserLiked, currentUserDisliked)

	numerator := float64(L1IntersectL2Size + D1IntersectD2Size - L1IntersectD2Size - L2IntersectD1Size)
	denominator := float64(unionSize4(L1, L2, D1, D2))

	similarity := numerator / denominator

	return similarity
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

func CalculateSimilaritiesWithOtherUsers(currentUserID uint, similarUsers []uint, currentUserLiked, currentUserDisliked []uint) (map[uint]float64, error) {
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

// P(U, B) = (ZL - ZD) / (|ML| + |MD|)
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
	Book        models.Book
	Probability float64
}

func GetRecommendedBooksSortedAndPaginated(recommendationProbabilities map[uint]float64, page int, perPage int) []Recommendation {
	var recommendations []Recommendation

	for bookID, probability := range recommendationProbabilities {
		var book models.Book
		database.DB.First(&book, bookID)
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

	return recommendations[startIdx:endIdx]
}
