export type DailyGoal = 5 | 10 | 15;

/**
 * Профиль ученика. В Go API авторизации и геймификации нет,
 * поэтому всё это живёт в браузере.
 */
export type Learner = {
  email: string;
  dailyGoal: DailyGoal;
  xp: number;
  lives: number;
  streak: number;
  /** Дата последнего пройденного урока, YYYY-MM-DD. */
  lastLessonOn: string | null;
  /** Категории, по которым урок уже пройден. */
  completedCategoryIds: number[];
};

export const MAX_LIVES = 5;
export const XP_PER_CORRECT_ANSWER = 5;
