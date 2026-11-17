package postgres

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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
	// goal is what the learner must say, help an optional nudge, and replies
	// are the accepted answers as plain lines; tiles and variants are built
	// from them once all scenes are known.
	speaker  string
	hint     string
	goal     string
	help     string
	replies  []string
	variants [][]string
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
	case "assemble", "dialogue":
		out["tiles"] = jsonItems(e.tiles)
		order := any(e.order)
		if e.kind == "dialogue" {
			out["speaker"] = e.speaker
			out["hint"] = e.hint
			out["goal"] = e.goal
			if e.help != "" {
				out["help"] = e.help
			}
			order = e.variants
		}
		raw, err := json.Marshal(order)
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

// defaulted returns fallback for an empty value.
func defaulted(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

// courseLockKey serializes course seeding between instances that start at the
// same time (a deploy overlapping with a restart).
const courseLockKey = 7_104_001

// courseHash fingerprints the course as it is written in code. Any edit to the
// lessons, the exercises or their keys changes it, so nobody has to remember
// to bump a version by hand.
func courseHash() string {
	sum := sha256.Sum256([]byte(fmt.Sprintf("%#v", courseStarter)))
	return hex.EncodeToString(sum[:])
}

// SeedCourse keeps the stored course equal to the one written in code.
//
// When the fingerprint of the code differs from the one stored (or none is
// stored, as in a database seeded before fingerprints existed), the course is
// replaced: situations, units, lessons and exercises are rebuilt, and the
// attempts that pointed at the old exercises go with them. Users, profiles,
// XP, lives and streaks stay; lesson progress starts over, because it is
// derived from those attempts.
//
// It reports whether anything was written. The work runs in one transaction
// under an advisory lock, so a failure leaves the old course in place and two
// instances starting together do not seed twice.
func SeedCourse(ctx context.Context, db *sql.DB) (bool, error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer rollback(tx)

	if _, err := tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock($1)`, courseLockKey); err != nil {
		return false, err
	}
	want := courseHash()
	var stored string
	err = tx.QueryRowContext(ctx, `SELECT value FROM course_meta WHERE key = 'course_hash'`).Scan(&stored)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	if stored == want {
		return false, nil
	}

	// Everything of the course hangs off situations with ON DELETE CASCADE.
	if _, err := tx.ExecContext(ctx, `DELETE FROM situations`); err != nil {
		return false, err
	}
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
	if _, err := tx.ExecContext(ctx, `
		INSERT INTO course_meta (key, value) VALUES ('course_hash', $1)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value`, want); err != nil {
		return false, err
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
