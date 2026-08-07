package dto

type HomepageHeroResponse struct {
	Headline     string `json:"headline"`
	SubHeadline  string `json:"sub_headline"`
	CallToAction string `json:"call_to_action"`
	CTAURL       string `json:"cta_url"`
	ImageURL     string `json:"image_url"`
}

type HomepageSectionResponse struct {
	Key       string `json:"key"`
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Content   string `json:"content"`
	ImageURL  string `json:"image_url"`
	LinkLabel string `json:"link_label,omitempty"`
	LinkURL   string `json:"link_url,omitempty"`
	SortOrder int    `json:"sort_order"`
}

type HomepageTestimonialResponse struct {
	Name      string `json:"name"`
	Role      string `json:"role"`
	Company   string `json:"company"`
	Message   string `json:"message"`
	AvatarURL string `json:"avatar_url"`
}
