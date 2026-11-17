import type { Exercise } from "@/entities/exercise";

/** GET /lessons/{id}. Закрытый урок отвечает 403, а не отдаёт содержимое. */
export type Lesson = {
  id: number;
  unit_id: number;
  title_kk: string;
  title_ru: string;
  description: string;
  position: number;
  /** Только completed или current: закрытый приходит как 403. */
  state: "completed" | "current";
  /** Индекс первого неотвеченного упражнения; у пройденного урока 0. */
  resume_position: number;
  correct: number;
  total: number;
  lives: number;
  max_lives: number;
  exercises: Exercise[];
};
