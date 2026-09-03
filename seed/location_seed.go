package seed

import (
	"errors"
	"strings"

	"nalakarsa/internal/model"

	"gorm.io/gorm"
)

type cityGroup struct {
	province string
	cities   []string
}

var provinces = []string{
	"Aceh",
	"Bali",
	"Kepulauan Bangka Belitung",
	"Banten",
	"Bengkulu",
	"DI Yogyakarta",
	"DKI Jakarta",
	"Gorontalo",
	"Jambi",
	"Jawa Barat",
	"Jawa Tengah",
	"Jawa Timur",
	"Kalimantan Barat",
	"Kalimantan Selatan",
	"Kalimantan Tengah",
	"Kalimantan Timur",
	"Kalimantan Utara",
	"Kepulauan Riau",
	"Lampung",
	"Maluku",
	"Maluku Utara",
	"Nusa Tenggara Barat",
	"Nusa Tenggara Timur",
	"Papua",
	"Papua Barat",
	"Riau",
	"Sulawesi Barat",
	"Sulawesi Selatan",
	"Sulawesi Tengah",
	"Sulawesi Tenggara",
	"Sulawesi Utara",
	"Sumatera Barat",
	"Sumatera Selatan",
	"Sumatera Utara",
}

var cityGroups = []cityGroup{
	{province: "Aceh", cities: []string{"Banda Aceh", "Lhokseumawe", "Sabang"}},
	{province: "Bali", cities: []string{"Denpasar", "Singaraja", "Gianyar", "Kuta"}},
	{province: "Kepulauan Bangka Belitung", cities: []string{"Pangkalpinang", "Tanjung Pandan", "Sungai Liat"}},
	{province: "Banten", cities: []string{"Serang", "Tangerang", "Cilegon", "Tangerang Selatan"}},
	{province: "Bengkulu", cities: []string{"Bengkulu"}},
	{province: "DI Yogyakarta", cities: []string{"Yogyakarta"}},
	{province: "DKI Jakarta", cities: []string{"Jakarta Pusat", "Jakarta Selatan", "Jakarta Barat", "Jakarta Utara", "Jakarta Timur", "Jakarta Barat"}},
	{province: "Gorontalo", cities: []string{"Gorontalo"}},
	{province: "Jambi", cities: []string{"Jambi", "Sungai Penuh", "Muara Bungo"}},
	{province: "Jawa Barat", cities: []string{"Bandung", "Bekasi", "Bogor", "Depok", "Cirebon", "Tasikmalaya"}},
	{province: "Jawa Tengah", cities: []string{"Semarang", "Solo", "Salatiga", "Pekalongan", "Magelang"}},
	{province: "Jawa Timur", cities: []string{"Surabaya", "Malang", "Kediri", "Sidoarjo", "Pasuruan"}},
	{province: "Kalimantan Barat", cities: []string{"Pontianak", "Singkawang", "Sungai Kakap"}},
	{province: "Kalimantan Selatan", cities: []string{"Banjarmasin", "Banjarbaru", "Martapura"}},
	{province: "Kalimantan Tengah", cities: []string{"Palangkaraya", "Buntok", "Pangkalan Bun"}},
	{province: "Kalimantan Timur", cities: []string{"Samarinda", "Balikpapan", "Bontang", "Sangatta"}},
	{province: "Kalimantan Utara", cities: []string{"Tarakan", "Tanjung Selor"}},
	{province: "Kepulauan Riau", cities: []string{"Tanjungpinang", "Batam", "Tarempa"}},
	{province: "Lampung", cities: []string{"Bandar Lampung", "Metro", "Pringsewu"}},
	{province: "Maluku", cities: []string{"Ambon", "Tual"}},
	{province: "Maluku Utara", cities: []string{"Sofifi", "Ternate", "Tidore"}},
	{province: "Nusa Tenggara Barat", cities: []string{"Mataram", "Bima", "Bima"}},
	{province: "Nusa Tenggara Timur", cities: []string{"Kupang", "Labuan Bajo", "Maumere"}},
	{province: "Papua", cities: []string{"Jayapura", "Timika", "Wamena"}},
	{province: "Papua Barat", cities: []string{"Manokwari", "Sorong", "Fakfak"}},
	{province: "Riau", cities: []string{"Pekanbaru", "Dumai"}},
	{province: "Sulawesi Barat", cities: []string{"Mamuju"}},
	{province: "Sulawesi Selatan", cities: []string{"Makassar", "Parepare", "Palopo"}},
	{province: "Sulawesi Tengah", cities: []string{"Palu"}},
	{province: "Sulawesi Tenggara", cities: []string{"Kendari", "Kolaka"}},
	{province: "Sulawesi Utara", cities: []string{"Manado", "Bitung", "Tomohon"}},
	{province: "Sumatera Barat", cities: []string{"Padang", "Bukittinggi", "Payakumbuh"}},
	{province: "Sumatera Selatan", cities: []string{"Palembang", "Prabumulih", "Lubuk Linggau"}},
	{province: "Sumatera Utara", cities: []string{"Medan", "Binjai", "Tebing Tinggi"}},
	{province: "Sulawesi Tengah", cities: []string{"Palu", "Tolitoli"}},
}

func SeedLocations(db *gorm.DB) error {
	for _, provinceName := range provinces {
		var province model.Location
		err := db.Where("LOWER(name) = LOWER(?) AND type = ?", strings.TrimSpace(provinceName), "province").First(&province).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		province = model.Location{
			Name:        provinceName,
			Type:        "province",
			CountryCode: "ID",
			Country:     "Indonesia",
			IsActive:    true,
		}
		if err := db.Create(&province).Error; err != nil {
			return err
		}
	}
	provinceByName := make(map[string]model.Location)
	var existingProvinces []model.Location
	if err := db.Where("type = ?", "province").Find(&existingProvinces).Error; err != nil {
		return err
	}
	for _, p := range existingProvinces {
		provinceByName[p.Name] = p
	}

	for _, group := range cityGroups {
		prov, ok := provinceByName[group.province]
		if !ok {
			continue
		}

		for _, cityName := range group.cities {
			var city model.Location
			err := db.Where("LOWER(name) = LOWER(?) AND type = ? AND province_id = ?", strings.TrimSpace(cityName), "city", prov.ID).First(&city).Error
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			city = model.Location{
				Name:         cityName,
				Type:         "city",
				ProvinceID:   &prov.ID,
				ProvinceName: prov.Name,
				CountryCode:  "ID",
				Country:      "Indonesia",
				IsActive:     true,
			}
			if err := db.Create(&city).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
