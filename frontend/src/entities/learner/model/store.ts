import { readJSON, writeJSON } from "@/shared/lib/storage";
import { MAX_LIVES, XP_PER_CORRECT_ANSWER, type DailyGoal, type Learner } from "./types";

const KEY = "qaztil.learner";

const listeners = new Set<() => void>();
let cache: Learner | null | undefined;

function today(): string {
  return new Date().toISOString().slice(0, 10);
}

function yesterday(): string {
  const date = new Date();
  date.setDate(date.getDate() - 1);
  return date.toISOString().slice(0, 10);
}

function emit() {
  for (const listener of listeners) listener();
}

function persist(learner: Learner) {
  cache = learner;
  writeJSON(KEY, learner);
  emit();
}

export const learnerStore = {
  subscribe(listener: () => void) {
    listeners.add(listener);
    return () => listeners.delete(listener);
  },

  get(): Learner | null {
    if (cache === undefined) {
      cache = readJSON<Learner>(KEY);
    }
    return cache;
  },

  /** На сервере профиля нет — рисуем пустое состояние до гидратации. */
  getServerSnapshot(): Learner | null {
    return null;
  },

  create(email: string, dailyGoal: DailyGoal): Learner {
    const learner: Learner = {
      email,
      dailyGoal,
      xp: 0,
      lives: MAX_LIVES,
      streak: 0,
      lastLessonOn: null,
      completedCategoryIds: [],
    };
    persist(learner);
    return learner;
  },

  setDailyGoal(dailyGoal: DailyGoal) {
    const learner = learnerStore.get();
    if (!learner) return;
    persist({ ...learner, dailyGoal });
  },

  /** Ошибка в вопросе стоит жизни; ниже нуля не уходим. */
  loseLife() {
    const learner = learnerStore.get();
    if (!learner) return;
    persist({ ...learner, lives: Math.max(0, learner.lives - 1) });
  },

  /** Урок закрыт: начисляем XP, продлеваем страйк, восстанавливаем жизни. */
  completeLesson(categoryId: number, correctAnswers: number) {
    const learner = learnerStore.get();
    if (!learner) return;

    const day = today();
    const streak =
      learner.lastLessonOn === day
        ? learner.streak
        : learner.lastLessonOn === yesterday()
          ? learner.streak + 1
          : 1;

    persist({
      ...learner,
      xp: learner.xp + correctAnswers * XP_PER_CORRECT_ANSWER,
      lives: MAX_LIVES,
      streak,
      lastLessonOn: day,
      completedCategoryIds: learner.completedCategoryIds.includes(categoryId)
        ? learner.completedCategoryIds
        : [...learner.completedCategoryIds, categoryId],
    });
  },

  reset() {
    cache = null;
    if (typeof window !== "undefined") {
      try {
        window.localStorage.removeItem(KEY);
      } catch {
        // молча: приватное окно
      }
    }
    emit();
  },
};
