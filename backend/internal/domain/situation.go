package domain

// Situation is a real-life context. The learner picks one and studies its words.
type Situation struct {
	ID          int64
	Slug        string
	TitleKK     string
	TitleRU     string
	Description string
	Position    int
	WordCount   int
}

// SituationWord is a Kazakh word or line used in one situation.
type SituationWord struct {
	ID            int64
	SituationID   int64
	Kazakh        string
	Russian       string
	Transcription string
	Position      int
}
