import type {
  AssembleExercise,
  Exercise,
  ExerciseAnswer,
  ExerciseResult,
  MatchExercise,
} from "../model/types";

/*
 * Временный адаптер: упражнения собираются из слов категории и
 * проверяются в браузере. Ключи ответов живут здесь, а не в упражнении,
 * чтобы UI работал с теми же данными, что потом придут с сервера.
 */

/** Что адаптеру нужно от слова. Совпадает с полями /words. */
type Word = {
  id: number;
  kazakh: string;
  russian: string;
  example_kk: string;
  example_ru: string;
};

const PAIRS = 4;
const EXTRA_TILES = 2;

const keys = new Map<string, { pairs?: Map<string, string>; order?: string[]; answer: string }>();

function shuffle<T>(items: T[]): T[] {
  const copy = [...items];
  for (let i = copy.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [copy[i], copy[j]] = [copy[j], copy[i]];
  }
  return copy;
}

function buildMatch(words: Word[]): MatchExercise | null {
  const picked = shuffle(words).slice(0, PAIRS);
  if (picked.length < 2) return null;
  const id = `match-${picked.map((word) => word.id).join("-")}`;
  keys.set(id, {
    pairs: new Map(picked.map((word) => [`l${word.id}`, `r${word.id}`])),
    answer: picked.map((word) => `${word.kazakh} — ${word.russian}`).join(" · "),
  });
  return {
    id,
    kind: "match",
    pairs_left: shuffle(picked.map((word) => ({ id: `l${word.id}`, text: word.kazakh }))),
    pairs_right: shuffle(picked.map((word) => ({ id: `r${word.id}`, text: word.russian }))),
  };
}

function buildAssemble(words: Word[]): AssembleExercise | null {
  const source = shuffle(words).find(
    (word) => word.example_kk && word.example_ru && word.example_kk.split(" ").length > 1,
  );
  if (!source) return null;
  const parts = source.example_kk.split(/\s+/);
  const distractors = shuffle(
    words.filter((word) => word.id !== source.id && !word.kazakh.includes(" ")),
  )
    .slice(0, EXTRA_TILES)
    .map((word) => word.kazakh);
  const tiles = [...parts, ...distractors].map((text, index) => ({
    id: `t${index}`,
    text,
  }));
  const id = `assemble-${source.id}`;
  keys.set(id, {
    order: parts.map((_, index) => `t${index}`),
    answer: source.example_kk,
  });
  return { id, kind: "assemble", prompt: source.example_ru, tiles: shuffle(tiles) };
}

/** Упражнения для урока по словам категории: пара и сборка фразы, если хватает данных. */
export function buildLocalExercises(words: Word[]): Exercise[] {
  return [buildMatch(words), buildAssemble(words)].filter(
    (exercise): exercise is Exercise => exercise !== null,
  );
}

/** Проверка ответа — повторяет то, что потом будет делать сервер. */
export function checkLocalExercise(
  exercise: Exercise,
  answer: ExerciseAnswer,
): ExerciseResult {
  const key = keys.get(exercise.id);
  if (!key) return { correct: false, correct_answer: "" };

  if (answer.kind === "match") {
    const correct =
      answer.pairs.length === key.pairs?.size &&
      answer.pairs.every(([left, right]) => key.pairs?.get(left) === right);
    return { correct, correct_answer: key.answer };
  }

  // Плитки с одинаковым текстом взаимозаменяемы — сравниваем по тексту.
  const tiles = exercise.kind === "assemble" ? exercise.tiles : [];
  const text = (ids: string[]) =>
    ids.map((id) => tiles.find((tile) => tile.id === id)?.text).join(" ");
  return {
    correct: text(answer.tile_ids) === text(key.order ?? []),
    correct_answer: key.answer,
  };
}

/** Верна ли одна пара — для подсветки прямо во время игры. */
export function isLocalPair(exercise: MatchExercise, left: string, right: string) {
  return keys.get(exercise.id)?.pairs?.get(left) === right;
}
