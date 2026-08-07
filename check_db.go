//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"log"
	"nalakarsa/internal/config"
	"nalakarsa/internal/database"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Critical Error: %v", err)
	}

	type Column struct {
		ColumnName string `gorm:"column:column_name"`
	}
	var columns []Column
	if err := db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name='profiles'").Scan(&columns).Error; err != nil {
		log.Fatalf("Failed to fetch columns: %v", err)
	}
	
	fmt.Println("Columns in profiles table:")
	for _, c := range columns {
		fmt.Println("-", c.ColumnName)
	}
}
