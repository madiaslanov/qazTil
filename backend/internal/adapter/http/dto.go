package httpapi

import "time"

type ErrorResponse struct {
	Error string `json:"error" example:"не найдено"`
} // @name ErrorResponse

type CategoryRequest struct {
	NameKK      string `json:"name_kk" example:"Сәлемдесу"`
	NameRU      string `json:"name_ru" example:"Приветствия"`
	Description string `json:"description" example:"Короткие фразы"`
} // @name CategoryRequest

type CategoryResponse struct {
	ID          int64  `json:"id" example:"1"`
	NameKK      string `json:"name_kk" example:"Сәлемдесу"`
	NameRU      string `json:"name_ru" example:"Приветствия"`
	Description string `json:"description" example:"Короткие фразы"`
} // @name CategoryResponse

type WordRequest struct {
	CategoryID    int64  `json:"category_id" example:"1"`
	Kazakh        string `json:"kazakh" example:"Сәлем"`
	Russian       string `json:"russian" example:"Привет"`
	Transcription string `json:"transcription" example:"sälem"`
	ExampleKK     string `json:"example_kk" example:"Сәлем!"`
	ExampleRU     string `json:"example_ru" example:"Привет!"`
} // @name WordRequest

type WordResponse struct {
	ID            int64  `json:"id" example:"1"`
	CategoryID    int64  `json:"category_id" example:"1"`
	Kazakh        string `json:"kazakh" example:"Сәлем"`
	Russian       string `json:"russian" example:"Привет"`
	Transcription string `json:"transcription" example:"sälem"`
	ExampleKK     string `json:"example_kk" example:"Сәлем!"`
	ExampleRU     string `json:"example_ru" example:"Привет!"`
} // @name WordResponse

type StartQuizRequest struct {
	CategoryID int64 `json:"category_id" example:"1"`
	Size       int   `json:"size" example:"5"`
} // @name StartQuizRequest

type AnswerRequest struct {
	QuestionID    int64 `json:"question_id" example:"1"`
	SelectedIndex int   `json:"selected_index" example:"0"`
} // @name AnswerRequest

type QuestionResponse struct {
	ID            int64    `json:"id" example:"1"`
	Prompt        string   `json:"prompt" example:"Привет"`
	Options       []string `json:"options"`
	SelectedIndex *int     `json:"selected_index,omitempty"`
	Correct       *bool    `json:"correct,omitempty"`
	CorrectIndex  *int     `json:"correct_index,omitempty"`
} // @name QuestionResponse

type ScoreResponse struct {
	Correct int `json:"correct" example:"1"`
	Total   int `json:"total" example:"5"`
} // @name ScoreResponse

type QuizResponse struct {
	ID         int64              `json:"id" example:"1"`
	CategoryID int64              `json:"category_id" example:"1"`
	CreatedAt  time.Time          `json:"created_at"`
	FinishedAt *time.Time         `json:"finished_at,omitempty"`
	Score      ScoreResponse      `json:"score"`
	Questions  []QuestionResponse `json:"questions"`
} // @name QuizResponse

type AnswerResponse struct {
	Correct      bool         `json:"correct" example:"true"`
	CorrectIndex int          `json:"correct_index" example:"2"`
	Quiz         QuizResponse `json:"quiz"`
} // @name AnswerResponse

type SituationResponse struct {
	ID          int64  `json:"id" example:"1"`
	Slug        string `json:"slug" example:"cafe"`
	TitleKK     string `json:"title_kk" example:"Дәмхана"`
	TitleRU     string `json:"title_ru" example:"В кафе"`
	Description string `json:"description" example:"Заказ и счёт"`
	Position    int    `json:"position" example:"1"`
	WordCount   int    `json:"word_count" example:"5"`
} // @name SituationResponse

type SituationWordResponse struct {
	ID            int64  `json:"id" example:"1"`
	SituationID   int64  `json:"situation_id" example:"1"`
	Kazakh        string `json:"kazakh" example:"Рақмет"`
	Russian       string `json:"russian" example:"Спасибо"`
	Transcription string `json:"transcription" example:"raqmet"`
	Position      int    `json:"position" example:"1"`
} // @name SituationWordResponse

type SituationDetailResponse struct {
	SituationResponse
	Words []SituationWordResponse `json:"words"`
} // @name SituationDetailResponse

type UnitResponse struct {
	ID          int64  `json:"id" example:"1"`
	SituationID int64  `json:"situation_id" example:"1"`
	TitleKK     string `json:"title_kk" example:"Дәмхана"`
	TitleRU     string `json:"title_ru" example:"В кафе"`
	Description string `json:"description" example:"Заказ и счёт"`
	Position    int    `json:"position" example:"1"`
	LessonCount int    `json:"lesson_count" example:"2"`
} // @name UnitResponse

type LessonResponse struct {
	ID          int64  `json:"id" example:"1"`
	UnitID      int64  `json:"unit_id" example:"1"`
	TitleKK     string `json:"title_kk" example:"Тапсырыс"`
	TitleRU     string `json:"title_ru" example:"Заказ"`
	Description string `json:"description" example:"Как попросить меню"`
	Position    int    `json:"position" example:"1"`
} // @name LessonResponse

type UnitDetailResponse struct {
	UnitResponse
	Lessons []LessonResponse `json:"lessons"`
} // @name UnitDetailResponse

type ExerciseResponse struct {
	ID       int64    `json:"id" example:"1"`
	LessonID int64    `json:"lesson_id" example:"1"`
	Kind     string   `json:"kind" example:"choice" enums:"choice,translate"`
	Prompt   string   `json:"prompt" example:"Спасибо"`
	Options  []string `json:"options"`
	Position int      `json:"position" example:"1"`
} // @name ExerciseResponse

type LessonDetailResponse struct {
	LessonResponse
	Exercises []ExerciseResponse `json:"exercises"`
} // @name LessonDetailResponse

type ExerciseAnswerRequest struct {
	Answer string `json:"answer" example:"Рақмет"`
} // @name ExerciseAnswerRequest

type ExerciseAnswerResponse struct {
	Correct  bool   `json:"correct" example:"true"`
	Expected string `json:"expected" example:"Рақмет"`
} // @name ExerciseAnswerResponse

type ProgressResponse struct {
	CategoryID     int64      `json:"category_id" example:"1"`
	CategoryNameKK string     `json:"category_name_kk" example:"Сәлемдесу"`
	CategoryNameRU string     `json:"category_name_ru" example:"Приветствия"`
	TotalAnswers   int        `json:"total_answers" example:"4"`
	CorrectAnswers int        `json:"correct_answers" example:"3"`
	LastStudiedAt  *time.Time `json:"last_studied_at,omitempty"`
} // @name ProgressResponse
