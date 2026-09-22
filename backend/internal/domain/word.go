package domain

// Word is a Kazakh vocabulary item with a Russian translation.
type Word struct {
	ID            int64
	CategoryID    int64
	Kazakh        string
	Russian       string
	Transcription string
	ExampleKK     string
	ExampleRU     string
}
