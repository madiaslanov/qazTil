/*
 * Упражнения урока — ровно то, что отдаёт GET /lessons/{id}.
 * Правильного ответа в них нет: проверка на сервере.
 * Идентификаторы упражнений — числа; id внутри упражнения (вариант, колонка,
 * плитка) — локальные строки вроде "o1", "l2", "t3".
 */

export type ExerciseItem = { id: string; text: string };

type ExerciseBase = {
  id: number;
  lesson_id: number;
  position: number;
};

export type ChoiceExercise = ExerciseBase & {
  kind: "choice";
  prompt: string;
  options: ExerciseItem[];
};

export type TranslateExercise = ExerciseBase & {
  kind: "translate";
  prompt: string;
};

export type MatchExercise = ExerciseBase & {
  kind: "match";
  prompt?: string;
  pairs_left: ExerciseItem[];
  pairs_right: ExerciseItem[];
};

export type AssembleExercise = ExerciseBase & {
  kind: "assemble";
  prompt: string;
  tiles: ExerciseItem[];
};

/**
 * Сцена: собеседник говорит реплику (prompt, по-казахски), ученик собирает
 * ответ из слов, как в assemble. hint — русский смысл реплики, чтобы можно
 * было следить за сценой.
 */
export type DialogueExercise = ExerciseBase & {
  kind: "dialogue";
  prompt: string;
  speaker: string;
  hint: string;
  tiles: ExerciseItem[];
};

export type Exercise =
  | ChoiceExercise
  | TranslateExercise
  | MatchExercise
  | AssembleExercise
  | DialogueExercise;

/** Тело POST /exercises/{id}/answers: форма зависит от типа. */
export type ExerciseAnswer =
  | { kind: "choice"; option_id: string; answer: string }
  | { kind: "dialogue"; tile_ids: string[] }
  | { kind: "translate"; answer: string }
  | { kind: "match"; pairs: [string, string][] }
  | { kind: "assemble"; tile_ids: string[] };

/** Один ответ на все типы. `expected` приходит всегда и готово к показу. */
export type ExerciseResult = {
  correct: boolean;
  expected: string;
  /** Остаток жизней после этого ответа. */
  lives: number;
  max_lives: number;
  /** Ушла последняя жизнь: урок остановлен. */
  lesson_failed: boolean;
  /** На все упражнения урока дан хотя бы один ответ. */
  lesson_completed: boolean;
  /** Сколько XP принёс именно этот ответ, а не общая сумма. */
  xp_earned: number;
  streak: number;
  /** false, если урок уже закрывали сегодня. */
  streak_increased: boolean;
  /** Открывшийся урок; null в конце пути. */
  next_lesson_id: number | null;
};
