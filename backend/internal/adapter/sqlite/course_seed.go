package sqlite

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
)

type courseWord struct {
	kazakh        string
	russian       string
	transcription string
}

type courseExercise struct {
	kind    string
	prompt  string
	answer  string
	options []string
}

type courseLesson struct {
	titleKK     string
	titleRU     string
	description string
	exercises   []courseExercise
}

type courseSituation struct {
	slug        string
	titleKK     string
	titleRU     string
	description string
	words       []courseWord
	lessons     []courseLesson
}

// SeedCourse inserts situational vocabulary and the unit/lesson/exercise path
// when the database has no situations yet. Existing category seeds are left alone.
func SeedCourse(ctx context.Context, db *sql.DB) (bool, error) {
	var count int
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM situations`).Scan(&count); err != nil {
		return false, err
	}
	if count > 0 {
		return false, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer rollback(tx)

	for i, situation := range courseStarter {
		situationID, err := insertSituation(ctx, tx, situation, i+1)
		if err != nil {
			return false, err
		}
		if err := insertWords(ctx, tx, situationID, situation.words); err != nil {
			return false, err
		}
		unitID, err := insertUnit(ctx, tx, situationID, situation, i+1)
		if err != nil {
			return false, err
		}
		if err := insertLessons(ctx, tx, unitID, situation.lessons); err != nil {
			return false, err
		}
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func insertSituation(ctx context.Context, tx *sql.Tx, situation courseSituation, position int) (int64, error) {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO situations (slug, title_kk, title_ru, description, position)
		VALUES (?, ?, ?, ?, ?)`,
		situation.slug, situation.titleKK, situation.titleRU, situation.description, position)
	if err != nil {
		return 0, fmt.Errorf("seed situation %s: %w", situation.slug, err)
	}
	return res.LastInsertId()
}

func insertWords(ctx context.Context, tx *sql.Tx, situationID int64, words []courseWord) error {
	for i, word := range words {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO situation_words (situation_id, kazakh, russian, transcription, position)
			VALUES (?, ?, ?, ?, ?)`,
			situationID, word.kazakh, word.russian, word.transcription, i+1); err != nil {
			return fmt.Errorf("seed situation word %s: %w", word.kazakh, err)
		}
	}
	return nil
}

func insertUnit(ctx context.Context, tx *sql.Tx, situationID int64, situation courseSituation, position int) (int64, error) {
	res, err := tx.ExecContext(ctx, `
		INSERT INTO units (situation_id, title_kk, title_ru, description, position)
		VALUES (?, ?, ?, ?, ?)`,
		situationID, situation.titleKK, situation.titleRU, situation.description, position)
	if err != nil {
		return 0, fmt.Errorf("seed unit %s: %w", situation.slug, err)
	}
	return res.LastInsertId()
}

func insertLessons(ctx context.Context, tx *sql.Tx, unitID int64, lessons []courseLesson) error {
	for i, lesson := range lessons {
		res, err := tx.ExecContext(ctx, `
			INSERT INTO lessons (unit_id, title_kk, title_ru, description, position)
			VALUES (?, ?, ?, ?, ?)`,
			unitID, lesson.titleKK, lesson.titleRU, lesson.description, i+1)
		if err != nil {
			return fmt.Errorf("seed lesson %s: %w", lesson.titleKK, err)
		}
		lessonID, err := res.LastInsertId()
		if err != nil {
			return err
		}
		for j, exercise := range lesson.exercises {
			if exercise.options == nil {
				exercise.options = []string{}
			}
			raw, err := json.Marshal(exercise.options)
			if err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO exercises (lesson_id, kind, prompt, answer, options_json, position)
				VALUES (?, ?, ?, ?, ?, ?)`,
				lessonID, exercise.kind, exercise.prompt, exercise.answer, string(raw), j+1); err != nil {
				return fmt.Errorf("seed exercise %s: %w", exercise.prompt, err)
			}
		}
	}
	return nil
}

var courseStarter = []courseSituation{
	{
		slug:        "cafe",
		titleKK:     "Дәмхана",
		titleRU:     "В кафе",
		description: "Заказ, счёт и короткие реакции за столом",
		words: []courseWord{
			{"Мәзір беріңізші", "Принесите меню, пожалуйста", "mäzir beriñizshi"},
			{"Бір шай, өтінемін", "Один чай, пожалуйста", "bir shai, ötinemin"},
			{"Есепшотты әкеліңізші", "Принесите счёт, пожалуйста", "esepshotty äkeliñizshi"},
			{"Дәмді екен", "Вкусно", "dämdi eken"},
			{"Рақмет", "Спасибо", "raqmet"},
		},
		lessons: []courseLesson{
			{
				titleKK: "Тапсырыс", titleRU: "Заказ",
				description: "Как попросить меню и напиток",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Принесите меню, пожалуйста", answer: "Мәзір беріңізші", options: []string{"Мәзір беріңізші", "Бір шай, өтінемін", "Рақмет"}},
					{kind: "choice", prompt: "Один чай, пожалуйста", answer: "Бір шай, өтінемін", options: []string{"Бір шай, өтінемін", "Есепшотты әкеліңізші", "Дәмді екен"}},
					{kind: "translate", prompt: "Принесите меню, пожалуйста", answer: "Мәзір беріңізші"},
				},
			},
			{
				titleKK: "Есеп", titleRU: "Счёт",
				description: "Как попросить счёт и ответить за столом",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Принесите счёт, пожалуйста", answer: "Есепшотты әкеліңізші", options: []string{"Есепшотты әкеліңізші", "Мәзір беріңізші", "Рақмет"}},
					{kind: "choice", prompt: "Вкусно", answer: "Дәмді екен", options: []string{"Дәмді екен", "Бір шай, өтінемін", "Рақмет"}},
					{kind: "translate", prompt: "Спасибо", answer: "Рақмет"},
				},
			},
		},
	},
	{
		slug:        "shop",
		titleKK:     "Дүкен",
		titleRU:     "В магазине",
		description: "Цена, выбор и оплата",
		words: []courseWord{
			{"Бұл қанша тұрады?", "Сколько это стоит?", "bul qansha turady?"},
			{"Қымбат", "Дорого", "qymbat"},
			{"Арзанырақ бар ма?", "Есть подешевле?", "arzanyraq bar ma?"},
			{"Мынаны алайын", "Я возьму это", "mynany alaiyn"},
			{"Картамен төлеймін", "Заплачу картой", "kartamen töleimin"},
		},
		lessons: []courseLesson{
			{
				titleKK: "Баға", titleRU: "Цена",
				description: "Как спросить цену и сказать, что дорого",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Сколько это стоит?", answer: "Бұл қанша тұрады?", options: []string{"Бұл қанша тұрады?", "Қымбат", "Мынаны алайын"}},
					{kind: "choice", prompt: "Дорого", answer: "Қымбат", options: []string{"Қымбат", "Бұл қанша тұрады?", "Арзанырақ бар ма?"}},
					{kind: "translate", prompt: "Есть подешевле?", answer: "Арзанырақ бар ма?"},
				},
			},
			{
				titleKK: "Төлем", titleRU: "Оплата",
				description: "Как взять вещь и заплатить",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Я возьму это", answer: "Мынаны алайын", options: []string{"Мынаны алайын", "Картамен төлеймін", "Қымбат"}},
					{kind: "choice", prompt: "Заплачу картой", answer: "Картамен төлеймін", options: []string{"Картамен төлеймін", "Бұл қанша тұрады?", "Мынаны алайын"}},
					{kind: "translate", prompt: "Сколько это стоит?", answer: "Бұл қанша тұрады?"},
				},
			},
		},
	},
	{
		slug:        "road",
		titleKK:     "Көлік",
		titleRU:     "В дороге",
		description: "Такси и билет",
		words: []courseWord{
			{"Такси шақырыңызшы", "Вызовите такси, пожалуйста", "taksi shaqyryñyzshy"},
			{"Вокзалға барайық", "Поехали на вокзал", "vokzalğa baraiyq"},
			{"Осында тоқтаңыз", "Остановите здесь", "osynda toqtañyz"},
			{"Қанша уақыт кетеді?", "Сколько времени займёт?", "qansha uaqyt ketedi?"},
			{"Билет бар ма?", "Есть билет?", "bilet bar ma?"},
		},
		lessons: []courseLesson{
			{
				titleKK: "Такси", titleRU: "Такси",
				description: "Как вызвать машину и назвать место",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Вызовите такси, пожалуйста", answer: "Такси шақырыңызшы", options: []string{"Такси шақырыңызшы", "Вокзалға барайық", "Билет бар ма?"}},
					{kind: "choice", prompt: "Поехали на вокзал", answer: "Вокзалға барайық", options: []string{"Вокзалға барайық", "Осында тоқтаңыз", "Такси шақырыңызшы"}},
					{kind: "translate", prompt: "Остановите здесь", answer: "Осында тоқтаңыз"},
				},
			},
			{
				titleKK: "Билет", titleRU: "Билет",
				description: "Как спросить время и билет",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Сколько времени займёт?", answer: "Қанша уақыт кетеді?", options: []string{"Қанша уақыт кетеді?", "Билет бар ма?", "Осында тоқтаңыз"}},
					{kind: "choice", prompt: "Есть билет?", answer: "Билет бар ма?", options: []string{"Билет бар ма?", "Такси шақырыңызшы", "Қанша уақыт кетеді?"}},
					{kind: "translate", prompt: "Вызовите такси, пожалуйста", answer: "Такси шақырыңызшы"},
				},
			},
		},
	},
}
