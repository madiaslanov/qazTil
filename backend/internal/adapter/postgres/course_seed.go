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
	// speaker and hint belong to a dialogue: who talks, and what the line means.
	speaker string
	hint    string
}

// payload builds the client-visible part and the stored key for one exercise.
func (e courseExercise) payload() (string, string, error) {
	out := map[string]any{}
	key := e.answer
	switch e.kind {
	case "choice", "dialogue":
		if e.kind == "dialogue" {
			out["speaker"] = e.speaker
			out["hint"] = e.hint
		}
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
