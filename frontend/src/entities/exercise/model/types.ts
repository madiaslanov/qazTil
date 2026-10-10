/*
 * Упражнения урока. Формы совпадают с контрактом бэкенда (BE-1/BE-2):
 * когда появится POST /exercises/{id}/answers, поменяется только адаптер.
 */

export type ExerciseItem = { id: string; text: string };

export type MatchExercise = {
  id: string;
  kind: "match";
  pairs_left: ExerciseItem[];
  pairs_right: ExerciseItem[];
};

export type AssembleExercise = {
  id: string;
  kind: "assemble";
  prompt: string;
  tiles: ExerciseItem[];
};

export type Exercise = MatchExercise | AssembleExercise;

export type ExerciseAnswer =
  | { kind: "match"; pairs: [string, string][] }
  | { kind: "assemble"; tile_ids: string[] };

export type ExerciseResult = {
  correct: boolean;
  correct_answer: string;
};
