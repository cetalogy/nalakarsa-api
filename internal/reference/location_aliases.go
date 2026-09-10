package reference

import "strings"

var locationAliases = map[string]map[string]string{
	"province": {
		"jabar":                      "Jawa Barat",
		"jawa barat":                 "Jawa Barat",
		"jateng":                     "Jawa Tengah",
		"jawa tengah":                "Jawa Tengah",
		"jatim":                      "Jawa Timur",
		"jawa timur":                 "Jawa Timur",
		"dki":                        "DKI Jakarta",
		"dki jakarta":                "DKI Jakarta",
		"diy":                        "DI Yogyakarta",
		"daerah istimewa yogyakarta": "DI Yogyakarta",
		"sumut":                      "Sumatera Utara",
		"sumbar":                     "Sumatera Barat",
		"sumsel":                     "Sumatera Selatan",
	},
	"city": {
		"bdg":             "Bandung",
		"kota bandung":    "Bandung",
		"jkt":             "Jakarta Pusat",
		"jogja":           "Yogyakarta",
		"yogya":           "Yogyakarta",
		"kota yogyakarta": "Yogyakarta",
		"surakarta":       "Surakarta",
		"solo":            "Surakarta",
	},
}

func NormalizeLocationSearch(value, locationType string) string {
	search := strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(value)), " "))
	if aliases, ok := locationAliases[strings.ToLower(strings.TrimSpace(locationType))]; ok {
		if canonical, ok := aliases[search]; ok {
			return canonical
		}
	}
	return value
}
