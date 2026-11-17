package postgres

import (
	"encoding/json"
	"strings"
	"testing"
)

// Тест читает сид так, как его отдаст API, и проверяет то, что в тексте
// легко испортить руками: ключи ответов, дубли и порядок плиток.
func TestCourseContentIsConsistent(t *testing.T) {
	for _, situation := range courseStarter {
		if len(situation.lessons) < 3 {
			t.Errorf("%s: уроков %d, ждали не меньше трёх", situation.slug, len(situation.lessons))
		}
		last := situation.lessons[len(situation.lessons)-1]
		for _, exercise := range last.exercises {
			if exercise.kind != "dialogue" {
				t.Errorf("%s: последний урок %q должен быть диалогом, а в нём %s", situation.slug, last.titleRU, exercise.kind)
			}
		}

		for _, lesson := range situation.lessons {
			if n := len(lesson.exercises); n < 5 || n > 8 {
				t.Errorf("%s/%s: упражнений %d, ждали от 5 до 8", situation.slug, lesson.titleRU, n)
			}
			seen := map[string]bool{}
			for _, exercise := range lesson.exercises {
				id := lesson.titleRU + " | " + exercise.kind + " | " + exercise.prompt
				if seen[id] {
					t.Errorf("повтор в уроке: %s", id)
				}
				seen[id] = true
				checkExercise(t, id, exercise)
			}
		}
	}
}

func checkExercise(t *testing.T, id string, e courseExercise) {
	t.Helper()
	switch e.kind {
	case "translate":
		t.Errorf("%s: печатать перевод больше не нужно", id)
	case "choice":
		count := 0
		for _, option := range e.options {
			if option == e.answer {
				count++
			}
		}
		if count != 1 {
			t.Errorf("%s: верный вариант встречается %d раз", id, count)
		}
		if len(e.options) < 3 {
			t.Errorf("%s: вариантов %d, ждали не меньше трёх", id, len(e.options))
		}
		unique := map[string]bool{}
		for _, option := range e.options {
			if unique[option] {
				t.Errorf("%s: вариант %q повторён", id, option)
			}
			unique[option] = true
		}
	case "assemble", "dialogue":
		if e.kind == "dialogue" && (e.speaker == "" || e.hint == "") {
			t.Errorf("%s: у реплики нет говорящего или подсказки", id)
		}
		payload, key, err := e.payload()
		if err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		text := map[string]string{}
		var decoded struct {
			Tiles []map[string]string `json:"tiles"`
		}
		if err := json.Unmarshal([]byte(payload), &decoded); err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		for _, item := range decoded.Tiles {
			text[item["id"]] = item["text"]
		}
		var order []string
		if err := json.Unmarshal([]byte(key), &order); err != nil {
			t.Fatalf("%s: %v", id, err)
		}
		words := make([]string, 0, len(order))
		for _, tileID := range order {
			word, ok := text[tileID]
			if !ok {
				t.Fatalf("%s: в ключе плитка %s, которой нет", id, tileID)
			}
			words = append(words, word)
		}
		if len(text) <= len(order) {
			t.Errorf("%s: в банке нет лишнего слова", id)
		}
		bank := make([]string, 0, len(decoded.Tiles))
		for _, item := range decoded.Tiles {
			bank = append(bank, item["text"])
		}
		if inOrder(bank, words) {
			t.Errorf("%s: слова уже стоят по порядку (%s)", id, strings.Join(words, " "))
		}
	case "match":
		if len(e.pairs) != len(e.pairsLeft) || len(e.pairsLeft) != len(e.pairsRight) {
			t.Errorf("%s: колонки и ключ разной длины", id)
		}
		right := map[string]string{}
		for _, item := range e.pairsRight {
			right[item.id] = item.text
		}
		left := map[string]string{}
		for _, item := range e.pairsLeft {
			left[item.id] = item.text
		}
		for _, p := range e.pairs {
			if left[p[0]] == "" || right[p[1]] == "" {
				t.Errorf("%s: ключ ссылается на несуществующую пару %v", id, p)
			}
		}
		// Правая колонка не должна стоять напротив своих пар.
		same := 0
		for i := range e.pairsLeft {
			if e.pairsRight[i].id == strings.Replace(e.pairsLeft[i].id, "l", "r", 1) {
				same++
			}
		}
		if same == len(e.pairsLeft) {
			t.Errorf("%s: пары стоят друг напротив друга", id)
		}
	default:
		t.Errorf("%s: неизвестный тип %q", id, e.kind)
	}
}
