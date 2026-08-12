package homepage

import (
	"nalakarsa/internal/dto"
	homerepository "nalakarsa/internal/repository/homepage"
)

type HomepageService interface {
	GetHero() (*dto.HomepageHeroResponse, error)
	GetSections() ([]dto.HomepageSectionResponse, error)
	GetTestimonials() ([]dto.HomepageTestimonialResponse, error)
}

type homepageService struct {
	repo homerepository.HomepageRepository
}

func NewHomepageService(repo homerepository.HomepageRepository) HomepageService {
	return &homepageService{
		repo: repo,
	}
}

var defaultHero = dto.HomepageHeroResponse{
	Headline:     "Kolaborasi Riset jadi Solusi Nyata",
	SubHeadline:  "Temukan akademisi, praktisi, dan profesional untuk mendorong riset yang berdampak.",
	CallToAction: "Jelajahi Diskusi",
	CTAURL:       "/discussions",
	ImageURL:     "/assets/homepage/hero.jpg",
}

var defaultSections = []dto.HomepageSectionResponse{
	{
		Key:       "benefits",
		Title:     "Mengapa Nalakarsa?",
		Subtitle:  "Platform penghubung ekosistem riset dan industri.",
		Content:   "Temukan kolaborasi lintas peran untuk mengubah ide riset menjadi implementasi nyata.",
		ImageURL:  "/assets/homepage/benefits.jpg",
		SortOrder: 1,
	},
	{
		Key:       "highlights",
		Title:     "Proyek & Kolaborasi",
		Subtitle:  "Akses peluang yang relevan.",
		Content:   "Lihat proyek, ikuti diskusi, dan bangun koneksi yang membantu kemajuan industri.",
		ImageURL:  "/assets/homepage/highlights.jpg",
		SortOrder: 2,
	},
}

var defaultTestimonials = []dto.HomepageTestimonialResponse{
	{
		Name:      "Dr. Raka Santoso",
		Role:      "Dosen",
		Company:   "Universitas Teknologi Nusantara",
		Message:   "Nalakarsa membantu saya menemukan mitra untuk mengimplementasikan riset ke pilot project.",
		AvatarURL: "/assets/testimonials/raka.jpg",
	},
	{
		Name:      "Sari Wijaya",
		Role:      "Praktisi",
		Company:   "AgriTech Indonesia",
		Message:   "Kolaborasi jadi lebih cepat karena diskusi dan proyek terstruktur dalam satu platform.",
		AvatarURL: "/assets/testimonials/sari.jpg",
	},
}

func (s *homepageService) GetHero() (*dto.HomepageHeroResponse, error) {
	hero, err := s.repo.GetHero()
	if err != nil {
		return nil, err
	}
	if hero != nil {
		return &dto.HomepageHeroResponse{
			Headline:     hero.Headline,
			SubHeadline:  hero.SubHeadline,
			CallToAction: hero.CallToAction,
			CTAURL:       hero.CTAURL,
			ImageURL:     hero.ImageURL,
		}, nil
	}
	return &defaultHero, nil
}

func (s *homepageService) GetSections() ([]dto.HomepageSectionResponse, error) {
	sections, err := s.repo.ListSections()
	if err != nil {
		return nil, err
	}
	if len(sections) > 0 {
		res := make([]dto.HomepageSectionResponse, 0, len(sections))
		for _, section := range sections {
			res = append(res, dto.HomepageSectionResponse{
				Key:       section.Key,
				Title:     section.Title,
				Subtitle:  section.Subtitle,
				Content:   section.Content,
				ImageURL:  section.ImageURL,
				LinkLabel: section.LinkLabel,
				LinkURL:   section.LinkURL,
				SortOrder: section.SortOrder,
			})
		}
		return res, nil
	}
	return defaultSections, nil
}

func (s *homepageService) GetTestimonials() ([]dto.HomepageTestimonialResponse, error) {
	testimonials, err := s.repo.ListTestimonials()
	if err != nil {
		return nil, err
	}
	if len(testimonials) > 0 {
		res := make([]dto.HomepageTestimonialResponse, 0, len(testimonials))
		for _, t := range testimonials {
			res = append(res, dto.HomepageTestimonialResponse{
				Name:      t.Name,
				Role:      t.Role,
				Company:   t.Company,
				Message:   t.Message,
				AvatarURL: t.AvatarURL,
			})
		}
		return res, nil
	}
	return defaultTestimonials, nil
}
