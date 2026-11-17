package postgres

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

type seedItem struct {
	id   string
	text string
}

type courseExercise struct {
	kind   string
	prompt string
	// answer is the key: a line for choice and translate, item ids for the rest.
	answer  string
	options []string
	// pairsLeft and pairsRight belong to a match exercise. The key is built
	// from pairs: each left id is joined to the right id next to it.
	pairsLeft  []seedItem
	pairsRight []seedItem
	pairs      [][2]string
	// tiles and order belong to an assemble exercise.
	tiles []seedItem
	order []string
}

// payload builds the client-visible part and the stored key for one exercise.
func (e courseExercise) payload() (string, string, error) {
	out := map[string]any{}
	key := e.answer
	switch e.kind {
	case "choice":
		if len(e.options) > 0 && len(e.pairsLeft) == 0 {
			items := make([]seedItem, 0, len(e.options))
			for i, text := range e.options {
				items = append(items, seedItem{id: fmt.Sprintf("o%d", i+1), text: text})
			}
			out["options"] = jsonItems(items)
			// The key moves from the line itself to the id holding it, so the
			// client answers with an id and never has to send the text back.
			for _, item := range items {
				if item.text == e.answer {
					key = item.id
					break
				}
			}
		}
	case "match":
		out["pairs_left"] = jsonItems(e.pairsLeft)
		out["pairs_right"] = jsonItems(e.pairsRight)
		raw, err := json.Marshal(e.pairs)
		if err != nil {
			return "", "", err
		}
		key = string(raw)
	case "assemble":
		out["tiles"] = jsonItems(e.tiles)
		raw, err := json.Marshal(e.order)
		if err != nil {
			return "", "", err
		}
		key = string(raw)
	}
	raw, err := json.Marshal(out)
	if err != nil {
		return "", "", err
	}
	return string(raw), key, nil
}

func jsonItems(items []seedItem) []map[string]string {
	out := make([]map[string]string, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]string{"id": item.id, "text": item.text})
	}
	return out
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
	// level is where in the course this unit sits, and topic is the learner
	// motivation it serves. Both come from the onboarding answers.
	level   string
	topic   string
	words   []courseWord
	lessons []courseLesson
}

// SeedCourse inserts situational vocabulary and the unit/lesson/exercise path
// when the database has no situations yet. Existing category seeds are left alone.
func defaulted(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// SeedCourse inserts every situation the build knows that the database does not
// have yet, matching on slug. Situations already stored are left untouched,
// so new content reaches a running deployment without wiping what learners did.
func SeedCourse(ctx context.Context, db *sql.DB) (bool, error) {
	known, err := storedSlugs(ctx, db)
	if err != nil {
		return false, err
	}
	missing := make([]courseSituation, 0, len(courseStarter))
	for _, situation := range courseStarter {
		if !known[situation.slug] {
			missing = append(missing, situation)
		}
	}
	if len(missing) == 0 {
		return false, nil
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer rollback(tx)

	// Positions continue after what is already stored, so an added unit lands
	// at the end of the path instead of fighting an existing one for a slot.
	offset := len(known)
	for i, situation := range missing {
		situationID, err := insertSituation(ctx, tx, situation, offset+i+1)
		if err != nil {
			return false, err
		}
		if err := insertWords(ctx, tx, situationID, situation.words); err != nil {
			return false, err
		}
		unitID, err := insertUnit(ctx, tx, situationID, situation, offset+i+1)
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
	id, err := insertID(ctx, tx, `
		INSERT INTO situations (slug, title_kk, title_ru, description, position)
		VALUES ($1, $2, $3, $4, $5)`,
		situation.slug, situation.titleKK, situation.titleRU, situation.description, position)
	if err != nil {
		return 0, fmt.Errorf("seed situation %s: %w", situation.slug, err)
	}
	return id, nil
}

func insertWords(ctx context.Context, tx *sql.Tx, situationID int64, words []courseWord) error {
	for i, word := range words {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO situation_words (situation_id, kazakh, russian, transcription, position)
			VALUES ($1, $2, $3, $4, $5)`,
			situationID, word.kazakh, word.russian, word.transcription, i+1); err != nil {
			return fmt.Errorf("seed situation word %s: %w", word.kazakh, err)
		}
	}
	return nil
}

func insertUnit(ctx context.Context, tx *sql.Tx, situationID int64, situation courseSituation, position int) (int64, error) {
	id, err := insertID(ctx, tx, `
		INSERT INTO units (situation_id, title_kk, title_ru, description, position, level, topic)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		situationID, situation.titleKK, situation.titleRU, situation.description, position,
		defaulted(situation.level, "A0"), defaulted(situation.topic, "life"))
	if err != nil {
		return 0, fmt.Errorf("seed unit %s: %w", situation.slug, err)
	}
	return id, nil
}

func insertLessons(ctx context.Context, tx *sql.Tx, unitID int64, lessons []courseLesson) error {
	for i, lesson := range lessons {
		lessonID, err := insertID(ctx, tx, `
			INSERT INTO lessons (unit_id, title_kk, title_ru, description, position)
			VALUES ($1, $2, $3, $4, $5)`,
			unitID, lesson.titleKK, lesson.titleRU, lesson.description, i+1)
		if err != nil {
			return fmt.Errorf("seed lesson %s: %w", lesson.titleKK, err)
		}
		for j, exercise := range lesson.exercises {
			if exercise.options == nil {
				exercise.options = []string{}
			}
			raw, err := json.Marshal(exercise.options)
			if err != nil {
				return err
			}
			payload, key, err := exercise.payload()
			if err != nil {
				return fmt.Errorf("seed exercise %s: %w", exercise.prompt, err)
			}
			if _, err := tx.ExecContext(ctx, `
				INSERT INTO exercises (lesson_id, kind, prompt, answer, options_json, position, payload)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`,
				lessonID, exercise.kind, exercise.prompt, key, string(raw), j+1, payload); err != nil {
				return fmt.Errorf("seed exercise %s: %w", exercise.prompt, err)
			}
		}
	}
	return nil
}

// Путь идёт от A0 к A1: сначала бытовые ситуации, затем работа и учёба.
var courseStarter = append(append([]courseSituation{}, courseLife...), courseWork, courseStudy)

var courseLife = []courseSituation{
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
					{kind: "assemble", prompt: "Соберите: «Принесите меню, пожалуйста»",
						tiles: []seedItem{{"t1", "беріңізші"}, {"t2", "шай"}, {"t3", "Мәзір"}},
						order: []string{"t3", "t1"}},
				},
			},
			{
				titleKK: "Есеп", titleRU: "Счёт",
				description: "Как попросить счёт и ответить за столом",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Принесите счёт, пожалуйста", answer: "Есепшотты әкеліңізші", options: []string{"Есепшотты әкеліңізші", "Мәзір беріңізші", "Рақмет"}},
					{kind: "choice", prompt: "Вкусно", answer: "Дәмді екен", options: []string{"Дәмді екен", "Бір шай, өтінемін", "Рақмет"}},
					{kind: "assemble", prompt: "Соберите: «Большое спасибо»",
						tiles: []seedItem{{"t1", "рақмет"}, {"t2", "Кешіріңіз"}, {"t3", "Көп"}},
						order: []string{"t3", "t1"}},
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
					{kind: "assemble", prompt: "Соберите: «Есть подешевле?»",
						tiles: []seedItem{{"t1", "бар"}, {"t2", "қымбат"}, {"t3", "ма?"}, {"t4", "Арзанырақ"}},
						order: []string{"t4", "t1", "t3"}},
				},
			},
			{
				titleKK: "Төлем", titleRU: "Оплата",
				description: "Как взять вещь и заплатить",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Я возьму это", answer: "Мынаны алайын", options: []string{"Мынаны алайын", "Картамен төлеймін", "Қымбат"}},
					{kind: "choice", prompt: "Заплачу картой", answer: "Картамен төлеймін", options: []string{"Картамен төлеймін", "Бұл қанша тұрады?", "Мынаны алайын"}},
					{kind: "assemble", prompt: "Соберите: «Сколько это стоит?»",
						tiles: []seedItem{{"t1", "қанша"}, {"t2", "Қымбат"}, {"t3", "тұрады?"}, {"t4", "Бұл"}},
						order: []string{"t4", "t1", "t3"}},
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
					{kind: "assemble", prompt: "Соберите: «Остановите здесь»",
						tiles: []seedItem{{"t1", "тоқтаңыз"}, {"t2", "Такси"}, {"t3", "Осында"}},
						order: []string{"t3", "t1"}},
				},
			},
			{
				titleKK: "Билет", titleRU: "Билет",
				description: "Как спросить время и билет",
				exercises: []courseExercise{
					{kind: "choice", prompt: "Сколько времени займёт?", answer: "Қанша уақыт кетеді?", options: []string{"Қанша уақыт кетеді?", "Билет бар ма?", "Осында тоқтаңыз"}},
					{kind: "choice", prompt: "Есть билет?", answer: "Билет бар ма?", options: []string{"Билет бар ма?", "Такси шақырыңызшы", "Қанша уақыт кетеді?"}},
					{kind: "assemble", prompt: "Соберите: «Вызовите такси, пожалуйста»",
						tiles: []seedItem{{"t1", "шақырыңызшы"}, {"t2", "Осында"}, {"t3", "Такси"}},
						order: []string{"t3", "t1"}},
				},
			},
		},
	},
}

// storedSlugs lists the situations already in the database.
func storedSlugs(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, `SELECT slug FROM situations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var slug string
		if err := rows.Scan(&slug); err != nil {
			return nil, err
		}
		out[slug] = true
	}
	return out, rows.Err()
}

// Юниты под мотивации work и study. Стартовые три (кафе, магазин, дорога)
// закрывают только life, поэтому ответ ученика на вопрос «зачем вам казахский»
// до сих пор не на что было применить.
var courseWork = courseSituation{
	slug:        "jumys",
	titleKK:     "Жұмыста",
	titleRU:     "На работе",
	description: "Знакомство с коллегами, встречи и просьбы",
	level:       "A1",
	topic:       "work",
	words: []courseWord{
		{"Мен бағдарламашымын", "Я программист", "men bağdarlamashymyn"},
		{"Қай бөлімде жұмыс істейсіз?", "В каком отделе вы работаете?", "qai bölimde jumys isteisiz"},
		{"Кездесу сағат үште", "Встреча в три часа", "kezdesu sağat üshte"},
		{"Көмектесе аласыз ба?", "Можете помочь?", "kömektese alasyz ba"},
		{"Есепті жібердім", "Я отправил отчёт", "esepti jiberdim"},
	},
	lessons: []courseLesson{
		{
			titleKK: "Танысу", titleRU: "Знакомство",
			description: "Как представиться и спросить о работе коллеги",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Я программист",
					answer:  "Мен бағдарламашымын",
					options: []string{"Мен бағдарламашымын", "Кездесу сағат үште", "Рақмет"}},
				{kind: "choice", prompt: "В каком отделе вы работаете?",
					answer:  "Қай бөлімде жұмыс істейсіз?",
					options: []string{"Қай бөлімде жұмыс істейсіз?", "Көмектесе аласыз ба?", "Есепті жібердім"}},
				{kind: "match", prompt: "Соедините фразу с переводом",
					pairsLeft: []seedItem{
						{"l1", "Мен бағдарламашымын"},
						{"l2", "Қай бөлімде жұмыс істейсіз?"},
						{"l3", "Көмектесе аласыз ба?"},
					},
					pairsRight: []seedItem{
						{"r1", "Можете помочь?"},
						{"r2", "Я программист"},
						{"r3", "В каком отделе вы работаете?"},
					},
					pairs: [][2]string{{"l1", "r2"}, {"l2", "r3"}, {"l3", "r1"}}},
				{kind: "assemble", prompt: "Соберите: «Я программист»",
					tiles: []seedItem{
						{"t1", "бағдарламашымын"}, {"t2", "Мен"},
						{"t3", "жұмыс"}, {"t4", "істеймін"},
					},
					order: []string{"t2", "t1"}},
				{kind: "assemble", prompt: "Соберите: «Можете помочь?»",
					tiles: []seedItem{{"t1", "аласыз"}, {"t2", "жұмыс"}, {"t3", "ба?"}, {"t4", "Көмектесе"}},
					order: []string{"t4", "t1", "t3"}},
			},
		},
		{
			titleKK: "Кездесу", titleRU: "Встреча",
			description: "Как договориться о времени",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Встреча в три часа",
					answer:  "Кездесу сағат үште",
					options: []string{"Кездесу сағат үште", "Есепті жібердім", "Мен бағдарламашымын"}},
				{kind: "assemble", prompt: "Соберите: «Встреча в три часа»",
					tiles: []seedItem{
						{"t1", "сағат"}, {"t2", "Кездесу"},
						{"t3", "үште"}, {"t4", "ертең"},
					},
					order: []string{"t2", "t1", "t3"}},
				{kind: "match", prompt: "Соедините фразу с переводом",
					pairsLeft: []seedItem{
						{"l1", "Кездесу сағат үште"},
						{"l2", "Есепті жібердім"},
					},
					pairsRight: []seedItem{
						{"r1", "Я отправил отчёт"},
						{"r2", "Встреча в три часа"},
					},
					pairs: [][2]string{{"l1", "r2"}, {"l2", "r1"}}},
				{kind: "assemble", prompt: "Соберите: «Встреча завтра»",
					tiles: []seedItem{{"t1", "ертең"}, {"t2", "үште"}, {"t3", "Кездесу"}},
					order: []string{"t3", "t1"}},
			},
		},
		{
			titleKK: "Есеп беру", titleRU: "Отчёт",
			description: "Как сказать, что работа сделана",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Я отправил отчёт",
					answer:  "Есепті жібердім",
					options: []string{"Есепті жібердім", "Көмектесе аласыз ба?", "Кездесу сағат үште"}},
				{kind: "assemble", prompt: "Соберите: «Я отправил отчёт»",
					tiles: []seedItem{
						{"t1", "жібердім"}, {"t2", "Есепті"}, {"t3", "алдым"},
					},
					order: []string{"t2", "t1"}},
				{kind: "assemble", prompt: "Соберите: «Отчёт готов»",
					tiles: []seedItem{{"t1", "дайын"}, {"t2", "жібердім"}, {"t3", "Есеп"}},
					order: []string{"t3", "t1"}},
				{kind: "choice", prompt: "Можете помочь?",
					answer:  "Көмектесе аласыз ба?",
					options: []string{"Көмектесе аласыз ба?", "Есепті жібердім", "Қай бөлімде жұмыс істейсіз?"}},
			},
		},
	},
}

var courseStudy = courseSituation{
	slug:        "oqu",
	titleKK:     "Оқуда",
	titleRU:     "На учёбе",
	description: "Занятия, библиотека и просьба повторить",
	level:       "A1",
	topic:       "study",
	words: []courseWord{
		{"Сабақ қашан басталады?", "Когда начинается занятие?", "sabaq qashan bastalady"},
		{"Дәрісхана қай жерде?", "Где аудитория?", "därishana qai jerde"},
		{"Кітап алғым келеді", "Хочу взять книгу", "kitap alğym keledi"},
		{"Түсінбедім, қайталаңызшы", "Не понял, повторите пожалуйста", "tüsinbedim, qaitalañyzshy"},
		{"Емтихан қиын болды", "Экзамен был трудным", "emtihan qiyn boldy"},
	},
	lessons: []courseLesson{
		{
			titleKK: "Сабақта", titleRU: "На занятии",
			description: "Как спросить про время и аудиторию",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Когда начинается занятие?",
					answer:  "Сабақ қашан басталады?",
					options: []string{"Сабақ қашан басталады?", "Дәрісхана қай жерде?", "Кітап алғым келеді"}},
				{kind: "choice", prompt: "Где аудитория?",
					answer:  "Дәрісхана қай жерде?",
					options: []string{"Дәрісхана қай жерде?", "Емтихан қиын болды", "Сабақ қашан басталады?"}},
				{kind: "match", prompt: "Соедините фразу с переводом",
					pairsLeft: []seedItem{
						{"l1", "Сабақ қашан басталады?"},
						{"l2", "Дәрісхана қай жерде?"},
						{"l3", "Түсінбедім, қайталаңызшы"},
					},
					pairsRight: []seedItem{
						{"r1", "Где аудитория?"},
						{"r2", "Не понял, повторите пожалуйста"},
						{"r3", "Когда начинается занятие?"},
					},
					pairs: [][2]string{{"l1", "r3"}, {"l2", "r1"}, {"l3", "r2"}}},
				{kind: "assemble", prompt: "Соберите: «Где аудитория?»",
					tiles: []seedItem{
						{"t1", "қай"}, {"t2", "Дәрісхана"},
						{"t3", "жерде"}, {"t4", "қашан"},
					},
					order: []string{"t2", "t1", "t3"}},
				{kind: "assemble", prompt: "Соберите: «Когда начинается занятие?»",
					tiles: []seedItem{{"t1", "қашан"}, {"t2", "кітап"}, {"t3", "басталады?"}, {"t4", "Сабақ"}},
					order: []string{"t4", "t1", "t3"}},
			},
		},
		{
			titleKK: "Кітапханада", titleRU: "В библиотеке",
			description: "Как попросить книгу",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Хочу взять книгу",
					answer:  "Кітап алғым келеді",
					options: []string{"Кітап алғым келеді", "Емтихан қиын болды", "Дәрісхана қай жерде?"}},
				{kind: "assemble", prompt: "Соберите: «Хочу взять книгу»",
					tiles: []seedItem{
						{"t1", "алғым"}, {"t2", "Кітап"},
						{"t3", "келеді"}, {"t4", "оқыдым"},
					},
					order: []string{"t2", "t1", "t3"}},
				{kind: "match", prompt: "Соедините фразу с переводом",
					pairsLeft: []seedItem{
						{"l1", "Кітап алғым келеді"},
						{"l2", "Емтихан қиын болды"},
					},
					pairsRight: []seedItem{
						{"r1", "Экзамен был трудным"},
						{"r2", "Хочу взять книгу"},
					},
					pairs: [][2]string{{"l1", "r2"}, {"l2", "r1"}}},
				{kind: "assemble", prompt: "Соберите: «Я взял книгу»",
					tiles: []seedItem{{"t1", "кітап"}, {"t2", "келеді"}, {"t3", "Мен"}, {"t4", "алдым"}},
					order: []string{"t3", "t1", "t4"}},
			},
		},
		{
			titleKK: "Емтихан", titleRU: "Экзамен",
			description: "Как сказать, что было трудно, и попросить повторить",
			exercises: []courseExercise{
				{kind: "choice", prompt: "Экзамен был трудным",
					answer:  "Емтихан қиын болды",
					options: []string{"Емтихан қиын болды", "Кітап алғым келеді", "Сабақ қашан басталады?"}},
				{kind: "assemble", prompt: "Соберите: «Не понял, повторите пожалуйста»",
					tiles: []seedItem{
						{"t1", "қайталаңызшы"}, {"t2", "Түсінбедім"}, {"t3", "білемін"},
					},
					order: []string{"t2", "t1"}},
				{kind: "assemble", prompt: "Соберите: «Экзамен был трудным»",
					tiles: []seedItem{{"t1", "қиын"}, {"t2", "Сабақ"}, {"t3", "болды"}, {"t4", "Емтихан"}},
					order: []string{"t4", "t1", "t3"}},
				{kind: "choice", prompt: "Не понял, повторите пожалуйста",
					answer:  "Түсінбедім, қайталаңызшы",
					options: []string{"Түсінбедім, қайталаңызшы", "Емтихан қиын болды", "Кітап алғым келеді"}},
			},
		},
	},
}
