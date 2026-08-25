package seed

import (
	"log"

	projectcommon "nalakarsa/internal/common/project"
	usercommon "nalakarsa/internal/common/user"
	"nalakarsa/internal/model"
	"nalakarsa/internal/utils"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) error {
	log.Println("Seeding database with default mock data...")

	// 1. Clear existing data
	log.Println("Purging existing data...")
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageSection{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageHero{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.HomepageTestimonial{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.ProjectApplication{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Project{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.DiscussionReply{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Discussion{})
	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.RefreshToken{})

	db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.User{})

	if err := SeedExpertise(db); err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword("password123")
	if err != nil {
		return err
	}

	// 2. Create Users & Profiles
	akademisiUser := &model.User{
		Email:        "dosen@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         usercommon.RoleAkademisi,
		SystemRole:   usercommon.SystemRoleUser,
		Status:       usercommon.StatusActive,
		FirstName:    "Budi",
		LastName:     "Santoso",
		FullName:     "Dr. Ir. Budi Santoso",
		PrefixTitle:  "Dr. Ir.",
		SuffixTitle:  "",
		Affiliation:  "Universitas Indonesia",
		Location:     "Jakarta",
		Expertise:    "Kecerdasan Buatan, Machine Learning",
		Bio:          "Peneliti bidang IoT.",
		Mission:      "Melakukan penelitian terapan dan mencari mitra industrialisasi.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=budi",
	}

	praktisiUser := &model.User{
		Email:        "pengusaha@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         usercommon.RolePraktisi,
		SystemRole:   usercommon.SystemRoleUser,
		Status:       usercommon.StatusActive,
		FirstName:    "Setyo",
		LastName:     "Nugroho",
		FullName:     "Setyo Nugroho",
		PrefixTitle:  "",
		SuffixTitle:  "",
		Affiliation:  "PT Telkom Indonesia",
		Location:     "Surabaya",
		Expertise:    "Data Science, Big Data Analytics",
		Bio:          "Pengusaha pertanian digital.",
		Mission:      "Mendukung otomatisasi pertanian menggunakan teknologi hasil riset.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=hendra",
	}

	profesionalUser := &model.User{
		Email:        "engineer@nalakarsa.id",
		PasswordHash: hashedPassword,
		Role:         usercommon.RoleProfesional,
		SystemRole:   usercommon.SystemRoleUser,
		Status:       usercommon.StatusActive,
		FirstName:    "Hendra",
		LastName:     "Wijaya",
		FullName:     "Hendra Wijaya",
		PrefixTitle:  "",
		SuffixTitle:  "",
		Affiliation:  "Tech Startup Indonesia",
		Location:     "Bandung",
		Expertise:    "Software Engineering, Cloud Architecture",
		Bio:          "Software Engineer spesialis backend.",
		Mission:      "Menyediakan solusi backend handal skala besar.",
		AvatarURL:    "https://api.dicebear.com/7.x/adventurer/svg?seed=setyo",
	}

	if err := db.Create(akademisiUser).Error; err != nil {
		return err
	}
	if err := db.Create(praktisiUser).Error; err != nil {
		return err
	}
	if err := db.Create(profesionalUser).Error; err != nil {
		return err
	}

	// 3. Create Discussion Topics
	disc1 := &model.Discussion{
		UserID:      profesionalUser.ID,
		Title:       "Mengapa Golang Sangat Cocok untuk Microservices?",
		Description: "Dalam ekosistem Nalakarsa, saya ingin berdiskusi mengenai konkurensi di Go. Goroutines dan Channels membuat pemrosesan paralel menjadi sangat efisien dibandingkan thread konvensional. Bagaimana pengalaman rekan-rekan?",
		Category:    "Tech",
		Tags:        "golang,microservices,backend",
		Status:      "open",
	}

	disc2 := &model.Discussion{
		UserID:      akademisiUser.ID,
		Title:       "Hilirisasi Hasil Riset Universitas ke Sektor Industri",
		Description: "Banyak prototipe riset IoT terhenti di laci laboratorium. Kendala utama adalah kurangnya komunikasi dengan pelaku bisnis. Mari kita bahas bagaimana menjembatani gap ini.",
		Category:    "Research",
		Tags:        "riset,industri,hilirisasi",
		Status:      "open",
	}

	if err := db.Create(disc1).Error; err != nil {
		return err
	}
	if err := db.Create(disc2).Error; err != nil {
		return err
	}

	// 4. Create Discussion Replies
	reply1 := &model.DiscussionReply{
		DiscussionID: disc1.ID,
		UserID:       akademisiUser.ID,
		Content:      "Setuju sekali! Kami di lab riset juga menggunakan Go untuk backend sensor gateway kami karena penggunaan memori yang sangat kecil.",
	}

	reply2 := &model.DiscussionReply{
		DiscussionID: disc2.ID,
		UserID:       praktisiUser.ID,
		Content:      "Kami dari sektor swasta siap menampung riset yang sudah matang di TRL 7 ke atas untuk dikomersialkan. Silakan hubungi kami.",
	}

	if err := db.Create(reply1).Error; err != nil {
		return err
	}
	if err := db.Create(reply2).Error; err != nil {
		return err
	}

	// 5. Create Projects (Replacing Collaboration)
	proj := &model.Project{
		OwnerID:       akademisiUser.ID,
		Title:         "Pengembangan Lapangan Sistem IoT Monitoring Tanah Pertanian",
		Description:   "Kami merancang sensor kelembapan tanah berbasis LoRaWAN. Kami mencari Praktisi Agribisnis yang memiliki lahan perkebunan untuk melakukan uji coba lapangan nyata dan validasi pasar.",
		Category:      "Research",
		Needs:         usercommon.RolePraktisi,
		Status:        projectcommon.ProjectStatusOpen,
		FundingStatus: "Self-funded",
		Location:      "Bogor",
	}

	if err := db.Create(proj).Error; err != nil {
		return err
	}

	// 6. Create Application
	app := &model.ProjectApplication{
		ProjectID:   proj.ID,
		ApplicantID: praktisiUser.ID,
		Message:     "Halo Pak Budi, saya Hendra dari PT Tani Maju Digital. Kami memiliki perkebunan teh seluas 3 hektar di daerah Bogor dan sangat tertarik untuk dijadikan lahan riset implementasi alat ini.",
		Status:      projectcommon.ApplicationStatusPending,
	}

	if err := db.Create(app).Error; err != nil {
		return err
	}

	// 7. Seed Homepage landing content
	hero := &model.HomepageHero{
		Headline:     "Kolaborasi Riset jadi Solusi Nyata",
		SubHeadline:  "Temukan akademisi, praktisi, dan profesional untuk mendorong riset yang berdampak.",
		CallToAction: "Jelajahi Diskusi",
		CTAURL:       "/discussions",
		ImageURL:     "/assets/homepage/hero.jpg",
		IsActive:     true,
	}

	sections := []model.HomepageSection{
		{
			Key:       "benefits",
			Title:     "Mengapa Nalakarsa?",
			Subtitle:  "Platform penghubung ekosistem riset dan industri.",
			Content:   "Temukan kolaborasi lintas peran untuk mengubah ide riset menjadi implementasi nyata.",
			ImageURL:  "/assets/homepage/benefits.jpg",
			SortOrder: 1,
			IsActive:  true,
		},
		{
			Key:       "highlights",
			Title:     "Proyek & Kolaborasi",
			Subtitle:  "Akses peluang yang relevan.",
			Content:   "Lihat proyek, ikuti diskusi, dan bangun koneksi yang membantu kemajuan industri.",
			ImageURL:  "/assets/homepage/highlights.jpg",
			SortOrder: 2,
			IsActive:  true,
		},
	}

	testimonials := []model.HomepageTestimonial{
		{
			Name:      "Dr. Raka Santoso",
			Role:      "Dosen",
			Company:   "Universitas Teknologi Nusantara",
			Message:   "Nalakarsa membantu saya menemukan mitra untuk mengimplementasikan riset ke pilot project.",
			AvatarURL: "/assets/testimonials/raka.jpg",
			IsActive:  true,
			SortOrder: 1,
		},
		{
			Name:      "Sari Wijaya",
			Role:      "Praktisi",
			Company:   "AgriTech Indonesia",
			Message:   "Kolaborasi jadi lebih cepat karena diskusi dan proyek terstruktur dalam satu platform.",
			AvatarURL: "/assets/testimonials/sari.jpg",
			IsActive:  true,
			SortOrder: 2,
		},
	}

	if err := db.Create(hero).Error; err != nil {
		return err
	}
	if err := db.Create(&sections).Error; err != nil {
		return err
	}
	if err := db.Create(&testimonials).Error; err != nil {
		return err
	}

	log.Println("Mock data seeding successfully completed.")
	return nil
}
