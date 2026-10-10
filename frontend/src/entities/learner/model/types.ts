export type DailyGoal = 5 | 10 | 15;

/** Стартовый уровень, который ученик выбирает на онбординге. */
export type Level = "A0" | "A1" | "A2";

/** Зачем ученику казахский — подбираем темы уроков. */
export type Motivation = "work" | "study" | "life";

/**
 * Профиль ученика. В Go API геймификации нет,
 * поэтому всё это живёт в браузере.
 */
export type Learner = {
  email: string;
  dailyGoal: DailyGoal;
  /** У профилей до редизайна уровня и мотивации нет. */
  level?: Level;
  motivation?: Motivation;
  xp: number;
  lives: number;
  streak: number;
  /** Дата последнего пройденного урока, YYYY-MM-DD. */
  lastLessonOn: string | null;
  /** Категории, по которым урок уже пройден. */
  completedCategoryIds: number[];
  /** Когда восстановится следующая жизнь, ms. null — запас полный. */
  nextLifeAt?: number | null;
};

/** Что ученик выбирает сам при регистрации. */
export type LearnerPrefs = Pick<Learner, "dailyGoal" | "level" | "motivation">;

export const MAX_LIVES = 5;
/** Одна жизнь восстанавливается за полчаса. */
export const LIFE_REGEN_MS = 30 * 60 * 1000;
export const XP_PER_CORRECT_ANSWER = 5;
