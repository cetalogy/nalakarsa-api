package seed

import (
	"nalakarsa/internal/model"

	"gorm.io/gorm"
)

func SeedExpertise(db *gorm.DB) error {
	items := []model.Expertise{
		{Name: "Artificial Intelligence", Category: "Technology"},
		{Name: "Machine Learning", Category: "Technology"},
		{Name: "Data Science", Category: "Technology"},
		{Name: "Software Engineering", Category: "Technology"},
		{Name: "Cybersecurity", Category: "Technology"},
		{Name: "Cloud Computing", Category: "Technology"},
		{Name: "Internet of Things", Category: "Technology"},
		{Name: "Education", Category: "Education"},
		{Name: "Public Health", Category: "Health"},
		{Name: "Medicine", Category: "Health"},
		{Name: "Agriculture", Category: "Agriculture"},
		{Name: "Environmental Science", Category: "Science"},
		{Name: "Engineering", Category: "Engineering"},
		{Name: "Business Management", Category: "Business"},
		{Name: "Economics", Category: "Social Science"},
		{Name: "Law", Category: "Social Science"},
		{Name: "Communication", Category: "Social Science"},
		{Name: "Psychology", Category: "Social Science"},
		{Name: "Sociology", Category: "Social Science"},
		{Name: "Arts and Design", Category: "Arts"},
	}
	for _, item := range items {
		if err := db.Where("category = ? AND name = ? AND specification = ?", item.Category, item.Name, item.Specification).FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}
