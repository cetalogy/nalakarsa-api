package seed

import (
	"errors"
	"strings"

	"nalakarsa/internal/model"

	"gorm.io/gorm"
)

var defaultSkills = []string{
	"Python", "JavaScript", "TypeScript", "Go", "Java", "C#", "C++", "PHP", "Kotlin", "Swift",
	"Artificial Intelligence", "Machine Learning", "Deep Learning", "Data Analysis", "Data Engineering",
	"Data Visualization", "Statistics", "Natural Language Processing", "Computer Vision", "Robotics",
	"Embedded Systems", "Cybersecurity", "Cloud Computing", "DevOps", "Backend Development",
	"Frontend Development", "Full Stack Development", "Mobile App Development", "UI/UX Design", "Product Design",
	"Software Testing", "Quality Assurance", "Database Administration", "SQL", "PostgreSQL", "API Development",
	"Microservices", "System Architecture", "Blockchain", "Geographic Information System", "Remote Sensing",
	"Renewable Energy", "Agricultural Technology", "Biotechnology", "Food Technology", "Environmental Management",
	"Project Management", "Business Analysis", "Digital Marketing", "Content Writing", "Copywriting",
	"Graphic Design", "Video Production", "Photography", "Public Speaking", "Research Methodology",
	"Academic Writing", "Curriculum Development", "Instructional Design", "Entrepreneurship", "Financial Analysis",
	"Accounting", "Legal Research", "Human Resources", "Supply Chain Management", "Operations Management",
	"Community Development", "Public Policy", "Healthcare Management", "Clinical Research",
}

func SeedSkills(db *gorm.DB) error {
	for _, name := range defaultSkills {
		var existing model.Skill
		err := db.Where("lower(name) = lower(?)", name).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&model.Skill{Name: strings.TrimSpace(name), IsActive: true}).Error; err != nil {
			return err
		}
	}
	return nil
}
