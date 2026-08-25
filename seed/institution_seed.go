package seed

import (
	"errors"
	"strings"

	"nalakarsa/internal/model"

	"gorm.io/gorm"
)

func SeedInstitutions(db *gorm.DB) error {
	institutions := []model.Institution{
		{
			Name:        "Universitas Indonesia",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Depok",
			Type:        "university",
		},
		{
			Name:        "Institut Teknologi Bandung",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Bandung",
			Type:        "university",
		},
		{
			Name:        "Institut Teknologi Sepuluh Nopember",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Surabaya",
			Type:        "university",
		},
		{
			Name:        "Universitas Gadjah Mada",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Yogyakarta",
			Type:        "university",
		},
		{
			Name:        "Institut Pertanian Bogor",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Bogor",
			Type:        "university",
		},
		{
			Name:        "Universitas Airlangga",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Surabaya",
			Type:        "university",
		},
		{
			Name:        "Universitas Diponegoro",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Semarang",
			Type:        "university",
		},
		{
			Name:        "Universitas Brawijaya",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Malang",
			Type:        "university",
		},
		{
			Name:        "Universitas Padjadjaran",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Bandung",
			Type:        "university",
		},
		{
			Name:        "University of Indonesia",
			CountryCode: "ID",
			Country:     "Indonesia",
			City:        "Depok",
			Type:        "university",
		},
		{
			Name:        "Massachusetts Institute of Technology",
			CountryCode: "US",
			Country:     "United States",
			City:        "Cambridge",
			Type:        "university",
		},
		{
			Name:        "Stanford University",
			CountryCode: "US",
			Country:     "United States",
			City:        "Stanford",
			Type:        "university",
		},
		{
			Name:        "Harvard University",
			CountryCode: "US",
			Country:     "United States",
			City:        "Cambridge",
			Type:        "university",
		},
		{
			Name:        "California Institute of Technology",
			CountryCode: "US",
			Country:     "United States",
			City:        "Pasadena",
			Type:        "university",
		},
		{
			Name:        "University of Oxford",
			CountryCode: "GB",
			Country:     "United Kingdom",
			City:        "Oxford",
			Type:        "university",
		},
		{
			Name:        "University of Cambridge",
			CountryCode: "GB",
			Country:     "United Kingdom",
			City:        "Cambridge",
			Type:        "university",
		},
		{
			Name:        "Imperial College London",
			CountryCode: "GB",
			Country:     "United Kingdom",
			City:        "London",
			Type:        "university",
		},
		{
			Name:        "National University of Singapore",
			CountryCode: "SG",
			Country:     "Singapore",
			City:        "Singapore",
			Type:        "university",
		},
		{
			Name:        "Nanyang Technological University",
			CountryCode: "SG",
			Country:     "Singapore",
			City:        "Singapore",
			Type:        "university",
		},
		{
			Name:        "University of Tokyo",
			CountryCode: "JP",
			Country:     "Japan",
			City:        "Tokyo",
			Type:        "university",
		},
		{
			Name:        "Seoul National University",
			CountryCode: "KR",
			Country:     "South Korea",
			City:        "Seoul",
			Type:        "university",
		},
		{
			Name:        "ETH Zurich",
			CountryCode: "CH",
			Country:     "Switzerland",
			City:        "Zurich",
			Type:        "university",
		},
		{
			Name:        "Peking University",
			CountryCode: "CN",
			Country:     "China",
			City:        "Beijing",
			Type:        "university",
		},
		{
			Name:        "Tsinghua University",
			CountryCode: "CN",
			Country:     "China",
			City:        "Beijing",
			Type:        "university",
		},
		{
			Name:        "University of Melbourne",
			CountryCode: "AU",
			Country:     "Australia",
			City:        "Melbourne",
			Type:        "university",
		},
		{
			Name:        "University of Sydney",
			CountryCode: "AU",
			Country:     "Australia",
			City:        "Sydney",
			Type:        "university",
		},
		{
			Name:        "Australian National University",
			CountryCode: "AU",
			Country:     "Australia",
			City:        "Canberra",
			Type:        "university",
		},
		{
			Name: "University of Toronto",
			CountryCode: "CA",
			Country:     "Canada",
			City:        "Toronto",
			Type:        "university",
		},
		{
			Name:        "University of Chicago",
			CountryCode: "US",
			Country:     "United States",
			City:        "Chicago",
			Type:        "university",
		},
		{
			Name:        "University of California, Berkeley",
			CountryCode: "US",
			Country:     "United States",
			City:        "Berkeley",
			Type:        "university",
		},
		{
			Name:        "University of California, Los Angeles",
			CountryCode: "US",
			Country:     "United States",
			City:        "Los Angeles",
			Type:        "university",
		},
	}

	for _, institution := range institutions {
		var existing model.Institution
		err := db.Where("LOWER(name) = LOWER(?)", strings.TrimSpace(institution.Name)).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&institution).Error; err != nil {
			return err
		}
	}

	return nil
}
