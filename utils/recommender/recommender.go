package recommender

import (
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
