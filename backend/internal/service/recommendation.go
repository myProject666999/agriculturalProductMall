package service

import (
	"agricultural-product-mall/internal/model"
	"agricultural-product-mall/pkg/database"
	"math"
	"sync"
	"time"
)

type UserBehaviorScore struct {
	View     float64
	Collect  float64
	Purchase float64
	Rating   float64
}

var defaultBehaviorScore = UserBehaviorScore{
	View:     1.0,
	Collect:  3.0,
	Purchase: 5.0,
	Rating:   2.0,
}

type RecommendationService struct {
	behaviorScore UserBehaviorScore
	mu            sync.RWMutex
}

func NewRecommendationService() *RecommendationService {
	return &RecommendationService{
		behaviorScore: defaultBehaviorScore,
	}
}

func (s *RecommendationService) RecordBehavior(userID, productID uint, behavior string, rating float64) error {
	var score float64
	switch behavior {
	case "view":
		score = s.behaviorScore.View
	case "collect":
		score = s.behaviorScore.Collect
	case "purchase":
		score = s.behaviorScore.Purchase
	case "rating":
		score = s.behaviorScore.Rating * rating
	default:
		score = 1.0
	}

	behaviorRecord := model.UserBehavior{
		UserID:    userID,
		ProductID: productID,
		Behavior:  behavior,
		Score:     score,
		CreatedAt: time.Now(),
	}

	return database.DB.Create(&behaviorRecord).Error
}

func (s *RecommendationService) GetUserItemMatrix() (map[uint]map[uint]float64, error) {
	var behaviors []model.UserBehavior
	if err := database.DB.Find(&behaviors).Error; err != nil {
		return nil, err
	}

	userItemMatrix := make(map[uint]map[uint]float64)
	for _, behavior := range behaviors {
		if _, exists := userItemMatrix[behavior.UserID]; !exists {
			userItemMatrix[behavior.UserID] = make(map[uint]float64)
		}
		userItemMatrix[behavior.UserID][behavior.ProductID] += behavior.Score
	}

	return userItemMatrix, nil
}

func (s *RecommendationService) cosineSimilarity(vec1, vec2 map[uint]float64) float64 {
	var dotProduct, norm1, norm2 float64

	for item, score1 := range vec1 {
		if score2, exists := vec2[item]; exists {
			dotProduct += score1 * score2
		}
		norm1 += score1 * score1
	}

	for _, score2 := range vec2 {
		norm2 += score2 * score2
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return dotProduct / (math.Sqrt(norm1) * math.Sqrt(norm2))
}

func (s *RecommendationService) FindSimilarUsers(userID uint, userItemMatrix map[uint]map[uint]float64, topN int) []uint {
	userVector := userItemMatrix[userID]
	if userVector == nil {
		return []uint{}
	}

	type similarity struct {
		userID uint
		score  float64
	}

	var similarities []similarity
	for otherUserID, otherVector := range userItemMatrix {
		if otherUserID == userID {
			continue
		}
		sim := s.cosineSimilarity(userVector, otherVector)
		if sim > 0 {
			similarities = append(similarities, similarity{userID: otherUserID, score: sim})
		}
	}

	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[i].score < similarities[j].score {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	var similarUsers []uint
	for i := 0; i < len(similarities) && i < topN; i++ {
		similarUsers = append(similarUsers, similarities[i].userID)
	}

	return similarUsers
}

func (s *RecommendationService) UserBasedRecommend(userID uint, limit int) ([]model.RecommendRecord, error) {
	userItemMatrix, err := s.GetUserItemMatrix()
	if err != nil {
		return nil, err
	}

	userVector := userItemMatrix[userID]
	similarUsers := s.FindSimilarUsers(userID, userItemMatrix, 10)

	scores := make(map[uint]float64)
	weights := make(map[uint]float64)

	for _, similarUserID := range similarUsers {
		simVector := userItemMatrix[similarUserID]
		sim := s.cosineSimilarity(userVector, simVector)

		for productID, score := range simVector {
			if userVector[productID] > 0 {
				continue
			}
			scores[productID] += score * sim
			weights[productID] += sim
		}
	}

	var recommendations []model.RecommendRecord
	for productID, score := range scores {
		if weights[productID] > 0 {
			normalizedScore := score / weights[productID]
			recommendations = append(recommendations, model.RecommendRecord{
				UserID:        userID,
				ProductID:     productID,
				Score:         normalizedScore,
				RecommendType: "user",
				CreatedAt:     time.Now(),
			})
		}
	}

	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (s *RecommendationService) GetItemUserMatrix(userItemMatrix map[uint]map[uint]float64) map[uint]map[uint]float64 {
	itemUserMatrix := make(map[uint]map[uint]float64)
	for userID, items := range userItemMatrix {
		for itemID, score := range items {
			if _, exists := itemUserMatrix[itemID]; !exists {
				itemUserMatrix[itemID] = make(map[uint]float64)
			}
			itemUserMatrix[itemID][userID] = score
		}
	}
	return itemUserMatrix
}

func (s *RecommendationService) FindSimilarItems(productID uint, itemUserMatrix map[uint]map[uint]float64, topN int) []uint {
	itemVector := itemUserMatrix[productID]
	if itemVector == nil {
		return []uint{}
	}

	type similarity struct {
		itemID uint
		score  float64
	}

	var similarities []similarity
	for otherItemID, otherVector := range itemUserMatrix {
		if otherItemID == productID {
			continue
		}
		sim := s.cosineSimilarity(itemVector, otherVector)
		if sim > 0 {
			similarities = append(similarities, similarity{itemID: otherItemID, score: sim})
		}
	}

	for i := 0; i < len(similarities); i++ {
		for j := i + 1; j < len(similarities); j++ {
			if similarities[i].score < similarities[j].score {
				similarities[i], similarities[j] = similarities[j], similarities[i]
			}
		}
	}

	var similarItems []uint
	for i := 0; i < len(similarities) && i < topN; i++ {
		similarItems = append(similarItems, similarities[i].itemID)
	}

	return similarItems
}

func (s *RecommendationService) ItemBasedRecommend(userID uint, limit int) ([]model.RecommendRecord, error) {
	userItemMatrix, err := s.GetUserItemMatrix()
	if err != nil {
		return nil, err
	}

	itemUserMatrix := s.GetItemUserMatrix(userItemMatrix)
	userVector := userItemMatrix[userID]

	scores := make(map[uint]float64)

	for ratedProductID, rating := range userVector {
		similarItems := s.FindSimilarItems(ratedProductID, itemUserMatrix, 10)
		for _, similarProductID := range similarItems {
			if userVector[similarProductID] > 0 {
				continue
			}
			sim := s.cosineSimilarity(
				itemUserMatrix[ratedProductID],
				itemUserMatrix[similarProductID],
			)
			scores[similarProductID] += rating * sim
		}
	}

	var recommendations []model.RecommendRecord
	for productID, score := range scores {
		recommendations = append(recommendations, model.RecommendRecord{
			UserID:        userID,
			ProductID:     productID,
			Score:         score,
			RecommendType: "item",
			CreatedAt:     time.Now(),
		})
	}

	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (s *RecommendationService) HybridRecommend(userID uint, limit int) ([]model.RecommendRecord, error) {
	userBasedRecs, err := s.UserBasedRecommend(userID, limit)
	if err != nil {
		return nil, err
	}

	itemBasedRecs, err := s.ItemBasedRecommend(userID, limit)
	if err != nil {
		return nil, err
	}

	combinedScores := make(map[uint]float64)
	for _, rec := range userBasedRecs {
		combinedScores[rec.ProductID] += rec.Score * 0.5
	}
	for _, rec := range itemBasedRecs {
		combinedScores[rec.ProductID] += rec.Score * 0.5
	}

	var hotProducts []model.Product
	database.DB.Where("status = ?", 1).
		Order("sales DESC, rating DESC").
		Limit(limit).
		Find(&hotProducts)

	for _, p := range hotProducts {
		if _, exists := combinedScores[p.ID]; !exists {
			combinedScores[p.ID] = float64(p.Sales) * 0.1
		}
	}

	var recommendations []model.RecommendRecord
	for productID, score := range combinedScores {
		recommendations = append(recommendations, model.RecommendRecord{
			UserID:        userID,
			ProductID:     productID,
			Score:         score,
			RecommendType: "hybrid",
			CreatedAt:     time.Now(),
		})
	}

	for i := 0; i < len(recommendations); i++ {
		for j := i + 1; j < len(recommendations); j++ {
			if recommendations[i].Score < recommendations[j].Score {
				recommendations[i], recommendations[j] = recommendations[j], recommendations[i]
			}
		}
	}

	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations, nil
}

func (s *RecommendationService) GetRecommendProducts(userID uint, limit int) ([]model.Product, error) {
	var recommendations []model.RecommendRecord
	var err error

	if userID > 0 {
		recommendations, err = s.HybridRecommend(userID, limit)
		if err != nil {
			return nil, err
		}
	}

	var productIDs []uint
	for _, rec := range recommendations {
		productIDs = append(productIDs, rec.ProductID)
	}

	if len(productIDs) == 0 {
		var hotProducts []model.Product
		err = database.DB.Where("status = ?", 1).
			Order("sales DESC, rating DESC").
			Limit(limit).
			Find(&hotProducts).Error
		return hotProducts, err
	}

	var products []model.Product
	err = database.DB.Where("id IN ? AND status = ?", productIDs, 1).
		Preload("Category").
		Preload("Merchant").
		Find(&products).Error

	if err != nil {
		return nil, err
	}

	productMap := make(map[uint]model.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	var sortedProducts []model.Product
	for _, rec := range recommendations {
		if p, exists := productMap[rec.ProductID]; exists {
			sortedProducts = append(sortedProducts, p)
		}
	}

	return sortedProducts, nil
}
