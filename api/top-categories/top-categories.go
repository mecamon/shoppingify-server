package top_categories

import (
	"github.com/mecamon/shoppingify-server/models"
)

func addPercentages(topCat, allCat []models.TopCategoryDTO) []models.TopCategoryDTO {
	var total = 0
	var completedCat []models.TopCategoryDTO

	for _, c := range allCat {
		total += c.SumQuantity
	}

	for _, c := range allCat {
		percentage := c.SumQuantity * 100 / total
		cat := models.TopCategoryDTO{
			ID:          c.ID,
			CategoryID:  c.CategoryID,
			Name:        c.Name,
			SumQuantity: c.SumQuantity,
			Percentage:  percentage,
		}
		for _, tc := range topCat {
			if tc.ID == c.ID {
				completedCat = append(completedCat, cat)
			}
		}
	}
	return completedCat
}
