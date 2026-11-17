package domain

import (
	"encoding/json"
	"sort"
	"strings"
)

// ExerciseItem is one labelled piece a learner can pick, drag or match.
// The id is local to the exercise, not a database key.
type ExerciseItem struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

// ExercisePayload is everything about an exercise the client may see.
// Only the fields of the exercise's own kind are filled.
type ExercisePayload struct {
	Options    []ExerciseItem `json:"options,omitempty"`
	PairsLeft  []ExerciseItem `json:"pairs_left,omitempty"`
	PairsRight []ExerciseItem `json:"pairs_right,omitempty"`
	Tiles      []ExerciseItem `json:"tiles,omitempty"`
	// Speaker and Hint belong to a dialogue: who is talking, and the Russian
	// meaning of their line so the learner can follow the scene.
	Speaker string `json:"speaker,omitempty"`
	Hint    string `json:"hint,omitempty"`
	// Goal tells the learner what to say, in Russian. Help is an optional nudge
	// for harder turns: word meanings, never the finished reply.
	Goal string `json:"goal,omitempty"`
	Help string `json:"help,omitempty"`
}

// ExerciseAnswer is a learner's answer in the shape its kind calls for.
type ExerciseAnswer struct {
	// Text answers a translate exercise, and stays for the older clients
	// that post {"answer": "..."}.
	Text     string
	OptionID string
	Pairs    [][2]string
	TileIDs  []string
}

// KnownExerciseKind reports whether the kind is one the course can store.
func KnownExerciseKind(kind string) bool {
	switch kind {
	case ExerciseChoice, ExerciseTranslate, ExerciseMatch, ExerciseAssemble, ExerciseDialogue:
		return true
	}
	return false
}

// CheckAnswer compares a learner's answer with the stored key.
//
// The key is held in Exercise.Answer: a plain line for choice and translate,
// and JSON for the kinds whose answer is more than one value. Keeping it in
// one column means no exercise can be served with its key attached by mistake —
// the payload the client receives simply has no room for it.
func CheckAnswer(exercise Exercise, answer ExerciseAnswer) (bool, string, error) {
	switch exercise.Kind {
	case ExerciseChoice:
		return checkChoice(exercise, answer)
	case ExerciseTranslate:
		got, want := normalizeLine(answer.Text), normalizeLine(exercise.Answer)
		if got == "" || want == "" {
			return false, "", ErrInvalid
		}
		return got == want, exercise.Answer, nil
	case ExerciseMatch:
		return checkMatch(exercise, answer)
	case ExerciseAssemble:
		return checkAssemble(exercise, answer)
	case ExerciseDialogue:
		return checkDialogue(exercise, answer)
	}
	return false, "", ErrInvalid
}

// checkChoice accepts either the option id or, for older clients, its text.
func checkChoice(exercise Exercise, answer ExerciseAnswer) (bool, string, error) {
	payload, err := exercise.DecodePayload()
	if err != nil {
		return false, "", err
	}
	want := strings.TrimSpace(exercise.Answer)
	if want == "" {
		return false, "", ErrInvalid
	}
	expected := want
	for _, option := range payload.Options {
		if option.ID == want {
			expected = option.Text
			break
		}
	}
	if id := strings.TrimSpace(answer.OptionID); id != "" {
		return id == want, expected, nil
	}
	if text := normalizeLine(answer.Text); text != "" {
		return text == normalizeLine(expected), expected, nil
	}
	return false, "", ErrInvalid
}

// checkMatch ignores the order the learner made the pairs in: only the
// set of connections counts.
func checkMatch(exercise Exercise, answer ExerciseAnswer) (bool, string, error) {
	want, err := decodePairs(exercise.Answer)
	if err != nil {
		return false, "", err
	}
	if len(want) == 0 {
		return false, "", ErrInvalid
	}
	expected := describePairs(exercise, want)
	if len(answer.Pairs) == 0 {
		return false, "", ErrInvalid
	}
	return samePairs(want, answer.Pairs), expected, nil
}

// checkAssemble compares the tiles in the order the learner placed them.
func checkAssemble(exercise Exercise, answer ExerciseAnswer) (bool, string, error) {
	var want []string
	if err := json.Unmarshal([]byte(exercise.Answer), &want); err != nil {
		return false, "", ErrInvalid
	}
	if len(want) == 0 {
		return false, "", ErrInvalid
	}
	expected := describeTiles(exercise, want)
	if len(answer.TileIDs) == 0 {
		return false, "", ErrInvalid
	}
	if len(answer.TileIDs) != len(want) {
		return false, expected, nil
	}
	for i, id := range want {
		if answer.TileIDs[i] != id {
			return false, expected, nil
		}
	}
	return true, expected, nil
}

// checkDialogue accepts any of the replies the exercise lists. The key is a
// list of tile-id sequences, one per accepted reply, and tiles are compared by
// their text, so a word that appears twice in the bank is interchangeable.
// A right answer echoes the learner's own reply; a wrong one lists what would
// have fit, so there is more than one way to read the scene.
func checkDialogue(exercise Exercise, answer ExerciseAnswer) (bool, string, error) {
	var variants [][]string
	if err := json.Unmarshal([]byte(exercise.Answer), &variants); err != nil || len(variants) == 0 {
		return false, "", ErrInvalid
	}
	if len(answer.TileIDs) == 0 {
		return false, "", ErrInvalid
	}
	payload, err := exercise.DecodePayload()
	if err != nil {
		return false, "", err
	}
	text := itemText(payload.Tiles)
	line := func(ids []string) []string {
		out := make([]string, 0, len(ids))
		for _, id := range ids {
			out = append(out, strings.ToLower(text[id]))
		}
		return out
	}
	got := line(answer.TileIDs)
	shown := make([]string, 0, len(variants))
	for _, variant := range variants {
		shown = append(shown, describeTiles(exercise, variant))
		want := line(variant)
		if len(want) != len(got) {
			continue
		}
		same := true
		for i := range want {
			if want[i] != got[i] {
				same = false
				break
			}
		}
		if same {
			return true, describeTiles(exercise, answer.TileIDs), nil
		}
	}
	return false, strings.Join(shown, " / "), nil
}

// DecodePayload reads the client-visible part of the exercise.
func (e Exercise) DecodePayload() (ExercisePayload, error) {
	var payload ExercisePayload
	if len(e.Payload) == 0 {
		return payload, nil
	}
	if err := json.Unmarshal(e.Payload, &payload); err != nil {
		return ExercisePayload{}, ErrInvalid
	}
	return payload, nil
}

func decodePairs(raw string) ([][2]string, error) {
	var pairs [][2]string
	if err := json.Unmarshal([]byte(raw), &pairs); err != nil {
		return nil, ErrInvalid
	}
	return pairs, nil
}

func samePairs(want, got [][2]string) bool {
	if len(want) != len(got) {
		return false
	}
	key := func(pairs [][2]string) []string {
		out := make([]string, 0, len(pairs))
		for _, pair := range pairs {
			out = append(out, pair[0]+"\x00"+pair[1])
		}
		sort.Strings(out)
		return out
	}
	a, b := key(want), key(got)
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// describePairs turns the stored key into something a learner can read,
// so the client has one line to show after a wrong answer.
func describePairs(exercise Exercise, pairs [][2]string) string {
	payload, err := exercise.DecodePayload()
	if err != nil {
		return ""
	}
	left := itemText(payload.PairsLeft)
	right := itemText(payload.PairsRight)
	parts := make([]string, 0, len(pairs))
	for _, pair := range pairs {
		parts = append(parts, left[pair[0]]+" — "+right[pair[1]])
	}
	return strings.Join(parts, ", ")
}

func describeTiles(exercise Exercise, ids []string) string {
	payload, err := exercise.DecodePayload()
	if err != nil {
		return ""
	}
	text := itemText(payload.Tiles)
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, text[id])
	}
	return strings.Join(parts, " ")
}

func itemText(items []ExerciseItem) map[string]string {
	out := make(map[string]string, len(items))
	for _, item := range items {
		out[item.ID] = item.Text
	}
	return out
}

func normalizeLine(value string) string {
	return strings.ToLower(strings.Join(strings.Fields(value), " "))
}
