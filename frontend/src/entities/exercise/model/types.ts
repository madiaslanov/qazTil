/*
 * Упражнения урока. Формы совпадают с контрактом бэкенда (BE-1/BE-2):
 * когда появится POST /exercises/{id}/answers, поменяется только адаптер.
 */

export type ExerciseItem = { id: string; text: string };

export type MatchExercise = {
  id: string;
  kind: "match";
  /** Порядок упражнения внутри урока. Локальные упражнения его не задают. */
  position?: number;
  pairs_left: ExerciseItem[];
  pairs_right: ExerciseItem[];
};

export type AssembleExercise = {
  id: string;
  kind: "assemble";
  position?: number;
  prompt: string;
  tiles: ExerciseItem[];
};

/** Выбор одного варианта. Появился в контракте API, локально не строится. */
export type ChoiceExercise = {
  id: string;
  kind: "choice";
  position: number;
  prompt: string;
  options: ExerciseItem[];
};

export type Exercise = MatchExercise | AssembleExercise;

/** Упражнение урока с сервера: все три типа из контракта. */
export type LessonExercise = Exercise | ChoiceExercise;

export type ExerciseAnswer =
  | { kind: "match"; pairs: [string, string][] }
  | { kind: "assemble"; tile_ids: string[] };

/** Ответ на упражнение урока с сервера. */
export type LessonAnswer =
  | ExerciseAnswer
  | { kind: "choice"; option_id: string };

export type ExerciseResult = {
  correct: boolean;
  correct_answer: string;
  /** Приходят только от сервера. */
  hearts_left?: number;
  lesson_failed?: boolean;
};
