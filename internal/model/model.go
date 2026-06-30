package model

// Lead is a customer contact/inquiry captured from the AHV Worldwide website.
type Lead struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Country   string `json:"country"`
	Inquiry   string `json:"inquiry"`
	Message   string `json:"message"`
	Source    string `json:"source"`
	Status    string `json:"status"`
	CreatedAt string `json:"createdAt"`
}

// SiteSettings is a flat key-value map of configurable site content.
type SiteSettings map[string]string
