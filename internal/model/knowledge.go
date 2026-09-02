package model

import (
	"time"

	"github.com/google/uuid"
)

type KnowledgeDomain struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex"`
	IsActive  bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type KnowledgeSubdomain struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	DomainID  uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_knowledge_subdomain_name"`
	Name      string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_knowledge_subdomain_name"`
	IsActive  bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}

type KnowledgeField struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	SubdomainID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex:idx_knowledge_field_name"`
	Name        string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_knowledge_field_name"`
	IsActive    bool      `gorm:"type:boolean;default:true;not null"`
	CreatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
