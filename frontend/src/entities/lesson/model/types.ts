import type { LessonExercise } from "@/entities/exercise/@x/lesson";
import type { LessonState } from "@/entities/path/@x/lesson";
import type { TodayProgress } from "@/entities/profile/@x/lesson";

/** GET /lessons/{id}. Правильных ответов здесь нет: проверка на сервере. */
export type Lesson = {
  id: string;
  unit_id: string;
  title_kk: string;
  title_ru: string;
  state: LessonState;
  resume_position?: number;
  hearts_left: number;
  exercises: LessonExercise[];
};

/** POST /lessons/{id}/start */
export type LessonStart = {
  resume_position: number;
  hearts_left: number;
  attempt: number;
};

/** POST /lessons/{id}/complete */
export type LessonComplete = {
  xp_earned: number;
  xp_total: number;
  streak_current: number;
  streak_increased: boolean;
  today: TodayProgress;
  /** Только что открывшийся урок или null, если юнит закончился. */
  next_lesson_id: string | null;
};
