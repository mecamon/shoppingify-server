package top_items

import "github.com/mecamon/shoppingify-server/models"

func addPercentages(topItem, allItem []models.TopItemDTO) []models.TopItemDTO {
	var total = 0
	var completedItems []models.TopItemDTO

	for _, c := range allItem {
		total += c.SumQuantity
	}

	for _, c := range allItem {
		percentage := c.SumQuantity * 100 / total
		item := models.TopItemDTO{
			ID:          c.ID,
			ItemID:      c.ItemID,
			Name:        c.Name,
			SumQuantity: c.SumQuantity,
			Percentage:  percentage,
		}
		for _, tc := range topItem {
			if tc.ID == c.ID {
				completedItems = append(completedItems, item)
			}
		}
	}
	return completedItems
}
