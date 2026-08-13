package homepage

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type HomepageHero struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Headline     string         `gorm:"type:varchar(255);not null"`
	SubHeadline  string         `gorm:"type:text;not null"`
	CallToAction string         `gorm:"type:varchar(120);default:'';not null"`
	CTAURL       string         `gorm:"type:varchar(255);default:'';not null"`
	ImageURL     string         `gorm:"type:varchar(255);default:'';not null"`
	IsActive     bool           `gorm:"default:true;not null"`
	CreatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt    time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type HomepageSection struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Key       string         `gorm:"type:varchar(120);uniqueIndex;not null"`
	Title     string         `gorm:"type:varchar(255);not null"`
	Subtitle  string         `gorm:"type:varchar(255);default:'';not null"`
	Content   string         `gorm:"type:text;default:'';not null"`
	ImageURL  string         `gorm:"type:varchar(255);default:'';not null"`
	LinkLabel string         `gorm:"type:varchar(120);default:'';not null"`
	LinkURL   string         `gorm:"type:varchar(255);default:'';not null"`
	SortOrder int            `gorm:"default:0;not null"`
	IsActive  bool           `gorm:"default:true;not null"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type HomepageTestimonial struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string         `gorm:"type:varchar(255);not null"`
	Role      string         `gorm:"type:varchar(255);default:'';not null"`
	Company   string         `gorm:"type:varchar(255);default:'';not null"`
	Message   string         `gorm:"type:text;not null"`
	AvatarURL string         `gorm:"type:varchar(255);default:'';not null"`
	IsActive  bool           `gorm:"default:true;not null"`
	SortOrder int            `gorm:"default:0;not null"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
