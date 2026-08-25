package dto

import "github.com/google/uuid"

type LocationSuggestion struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	ProvinceID   *uuid.UUID `json:"provinceId"`
	ProvinceName string     `json:"provinceName"`
	IsProvince   bool       `json:"isProvince"`
	IsCity       bool       `json:"isCity"`
}
