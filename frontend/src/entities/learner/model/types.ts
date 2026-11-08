export type DailyGoal = 5 | 10 | 15;

/** Стартовый уровень, который ученик выбирает на онбординге. */
export type Level = "A0" | "A1" | "A2";

/** Зачем ученику казахский — подбираем темы уроков. */
export type Motivation = "work" | "study" | "life";

/**
 * Локальный прогресс по старым урокам-категориям.
 * Prefs и XP/жизни/стрик живут в Go (/me/profile).
 */
export type Learner = {
  email: string;
  /** Категории, по которым урок уже пройден. */
  completedCategoryIds: number[];
};

/** Что ученик выбирает сам при регистрации — уходит в /auth/register. */
export type LearnerPrefs = {
  dailyGoal: DailyGoal;
  level: Level;
  motivation: Motivation;
};

export const MAX_LIVES = 5;
/** Одна жизнь восстанавливается за полчаса. */
export const LIFE_REGEN_MS = 30 * 60 * 1000;
export const XP_PER_CORRECT_ANSWER = 5;
