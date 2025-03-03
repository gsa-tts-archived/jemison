package searching

type RequiredQueryParameters struct {
	Affiliate   string
	SearchQuery string
}

type OptionalQueryParameters struct {
	EnableHighlighting bool
	Offset             int
	SortBy             int
	Sitelimit          int
}

type QueryWebResultsData struct {
	Title           string `json:"title"`
	URL             string `json:"url"`
	Snippet         string `json:"snippet"`
	PublicationDate string `json:"publication_date"`
	ThumbnailURL    string `json:"thumbnail_url"`
}

type QueryWebData struct {
	Total              int                   `json:"total"`
	NextOffset         int                   `json:"next_offset"`
	SpellingCorrection string                `json:"spelling_correction"`
	Results            []QueryWebResultsData `json:"results"`
}

type QueryData struct {
	SearchedQuery            string       `json:"query"`
	Web                      QueryWebData `json:"web"`
	TextBestBets             []string     `json:"text_best_bets"`
	GraphicBestBets          []string     `json:"graphic_best_bets"`
	HealthTopics             []string     `json:"health_topics"`
	JobOpenings              []string     `json:"job_openings"`
	FederalRegisterDocuments []string     `json:"federal_register_documents"`
	RelatedSearchTerms       []string     `json:"related_search_terms"`
}
