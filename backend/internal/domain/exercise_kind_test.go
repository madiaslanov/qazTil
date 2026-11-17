package domain_test

import (
	"testing"

	"github.com/madiaslanov/qazTil/internal/domain"
)

func matchExercise() domain.Exercise {
	return domain.Exercise{
		ID: 1, LessonID: 1, Kind: domain.ExerciseMatch,
		Payload: []byte(`{
			"pairs_left":[{"id":"l1","text":"Рақмет"},{"id":"l2","text":"Кешіріңіз"}],
			"pairs_right":[{"id":"r1","text":"Извините"},{"id":"r2","text":"Спасибо"}]
		}`),
		Answer: `[["l1","r2"],["l2","r1"]]`,
	}
}

func assembleExercise() domain.Exercise {
	return domain.Exercise{
		ID: 2, LessonID: 1, Kind: domain.ExerciseAssemble,
		Payload: []byte(`{"tiles":[
			{"id":"t1","text":"кофе"},{"id":"t2","text":"Мен"},
			{"id":"t3","text":"ішемін"},{"id":"t4","text":"шай"}
		]}`),
		Answer: `["t2","t1","t3"]`,
	}
}

func choiceExercise() domain.Exercise {
	return domain.Exercise{
		ID: 3, LessonID: 1, Kind: domain.ExerciseChoice,
		Payload: []byte(`{"options":[
			{"id":"o1","text":"Рақмет"},{"id":"o2","text":"Сәлем"}
		]}`),
		Answer: "o1",
	}
}

func TestMatchIgnoresPairOrder(t *testing.T) {
	// Сопоставление — это набор связей, а не последовательность:
	// ученик мог соединить пары в любом порядке.
	cases := map[string][][2]string{
		"в порядке ключа": {{"l1", "r2"}, {"l2", "r1"}},
		"в обратном":      {{"l2", "r1"}, {"l1", "r2"}},
	}
	for name, pairs := range cases {
		t.Run(name, func(t *testing.T) {
			correct, expected, err := domain.CheckAnswer(matchExercise(), domain.ExerciseAnswer{Pairs: pairs})
			if err != nil {
				t.Fatalf("ошибка: %v", err)
			}
			if !correct {
				t.Fatal("ждали верный ответ")
			}
			if expected != "Рақмет — Спасибо, Кешіріңіз — Извините" {
				t.Fatalf("подсказка = %q", expected)
			}
		})
	}
}

func TestMatchRejectsWrongAndPartial(t *testing.T) {
	cases := map[string][][2]string{
		"перепутаны пары": {{"l1", "r1"}, {"l2", "r2"}},
		"только одна":     {{"l1", "r2"}},
		"лишняя пара":     {{"l1", "r2"}, {"l2", "r1"}, {"l1", "r1"}},
	}
	for name, pairs := range cases {
		t.Run(name, func(t *testing.T) {
			correct, _, err := domain.CheckAnswer(matchExercise(), domain.ExerciseAnswer{Pairs: pairs})
			if err != nil {
				t.Fatalf("ошибка: %v", err)
			}
			if correct {
				t.Fatal("ждали неверный ответ")
			}
		})
	}
}

func TestAssembleKeepsOrder(t *testing.T) {
	// Сборка фразы — наоборот: порядок и есть ответ.
	correct, expected, err := domain.CheckAnswer(assembleExercise(),
		domain.ExerciseAnswer{TileIDs: []string{"t2", "t1", "t3"}})
	if err != nil {
		t.Fatalf("ошибка: %v", err)
	}
	if !correct {
		t.Fatal("ждали верный ответ")
	}
	if expected != "Мен кофе ішемін" {
		t.Fatalf("подсказка = %q", expected)
	}

	wrong := map[string][]string{
		"переставлены":       {"t1", "t2", "t3"},
		"взята лишняя":       {"t2", "t1", "t3", "t4"},
		"не хватает плитки":  {"t2", "t1"},
		"только отвлекающая": {"t4"},
	}
	for name, ids := range wrong {
		t.Run(name, func(t *testing.T) {
			correct, _, err := domain.CheckAnswer(assembleExercise(), domain.ExerciseAnswer{TileIDs: ids})
			if err != nil {
				t.Fatalf("ошибка: %v", err)
			}
			if correct {
				t.Fatal("ждали неверный ответ")
			}
		})
	}
}

func TestChoiceAcceptsIDAndLegacyText(t *testing.T) {
	correct, expected, err := domain.CheckAnswer(choiceExercise(), domain.ExerciseAnswer{OptionID: "o1"})
	if err != nil || !correct {
		t.Fatalf("по id: correct=%v err=%v", correct, err)
	}
	if expected != "Рақмет" {
		t.Fatalf("подсказка = %q", expected)
	}
	if correct, _, _ := domain.CheckAnswer(choiceExercise(), domain.ExerciseAnswer{OptionID: "o2"}); correct {
		t.Fatal("чужой вариант принят за верный")
	}
	// Старые клиенты шлют текст вместо идентификатора.
	if correct, _, _ := domain.CheckAnswer(choiceExercise(), domain.ExerciseAnswer{Text: "  рақмет "}); !correct {
		t.Fatal("текстовый ответ старого клиента отклонён")
	}
}

func TestLegacyChoiceWithoutPayloadStillWorks(t *testing.T) {
	// Упражнения, засеянные до появления payload: ключ хранит текст.
	legacy := domain.Exercise{ID: 9, Kind: domain.ExerciseChoice, Answer: "Рақмет"}
	correct, expected, err := domain.CheckAnswer(legacy, domain.ExerciseAnswer{Text: "Рақмет"})
	if err != nil || !correct || expected != "Рақмет" {
		t.Fatalf("correct=%v expected=%q err=%v", correct, expected, err)
	}
}

func TestEmptyAndUnknownAnswersAreRejected(t *testing.T) {
	if _, _, err := domain.CheckAnswer(matchExercise(), domain.ExerciseAnswer{}); err == nil {
		t.Fatal("пустое сопоставление принято")
	}
	if _, _, err := domain.CheckAnswer(assembleExercise(), domain.ExerciseAnswer{}); err == nil {
		t.Fatal("пустая сборка принята")
	}
	unknown := domain.Exercise{ID: 8, Kind: "doodle", Answer: "x"}
	if _, _, err := domain.CheckAnswer(unknown, domain.ExerciseAnswer{Text: "x"}); err == nil {
		t.Fatal("неизвестный тип принят")
	}
}

func TestKnownExerciseKind(t *testing.T) {
	for _, kind := range []string{domain.ExerciseChoice, domain.ExerciseTranslate,
		domain.ExerciseMatch, domain.ExerciseAssemble} {
		if !domain.KnownExerciseKind(kind) {
			t.Fatalf("%q должен быть известен", kind)
		}
	}
	if domain.KnownExerciseKind("doodle") {
		t.Fatal("doodle не должен быть известен")
	}
}

func TestDialogueIsCheckedLikeAssemble(t *testing.T) {
	exercise := assembleExercise()
	exercise.Kind = domain.ExerciseDialogue
	exercise.Payload = []byte(`{"speaker":"Даяшы","hint":"Что желаете?","tiles":[
		{"id":"t1","text":"шай"},{"id":"t2","text":"Бір"},{"id":"t3","text":"өтінемін"},{"id":"t4","text":"самса"}
	]}`)
	exercise.Answer = `["t2","t1","t3"]`

	correct, expected, err := domain.CheckAnswer(exercise, domain.ExerciseAnswer{TileIDs: []string{"t2", "t1", "t3"}})
	if err != nil || !correct || expected != "Бір шай өтінемін" {
		t.Fatalf("верный ответ: correct=%v expected=%q err=%v", correct, expected, err)
	}
	correct, _, err = domain.CheckAnswer(exercise, domain.ExerciseAnswer{TileIDs: []string{"t1", "t2", "t3"}})
	if err != nil || correct {
		t.Fatalf("порядок плиток и есть ответ: correct=%v err=%v", correct, err)
	}
	if !domain.KnownExerciseKind(domain.ExerciseDialogue) {
		t.Fatal("dialogue должен быть известным типом")
	}
}
