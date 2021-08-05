package pkg

// An article is meant with 'post' and not the request methods 'POST'
// Just so you know

type Post struct {
	Title string `json:"title"`
	LogoType string `json:"logo-type"`
	BannerType string `json:"banner-type"`
	Description string `json:"description"`
	LongDescription string `json:"long-description"`
}