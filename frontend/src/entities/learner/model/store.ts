import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import {
  LIFE_REGEN_MS,
  MAX_LIVES,
  XP_PER_CORRECT_ANSWER,
  type DailyGoal,
  type Learner,
  type LearnerPrefs,
} from "./types";

type LearnerState = {
  learner: Learner | null;
  /** persist поднимается вручную в провайдере, до этого состояние пустое. */
  hydrated: boolean;
  create: (email: string, prefs: LearnerPrefs) => void;
  setDailyGoal: (dailyGoal: DailyGoal) => void;
  loseLife: () => void;
  /** Начисляет жизни, которые успели восстановиться к моменту now. */
  regenerateLives: (now: number) => void;
  completeLesson: (categoryId: number, correctAnswers: number) => void;
  reset: () => void;
  markHydrated: () => void;
};

function day(offset = 0): string {
  const date = new Date();
  date.setDate(date.getDate() + offset);
  return date.toISOString().slice(0, 10);
}

export const useLearnerStore = create<LearnerState>()(
  persist(
    (set) => ({
      learner: null,
      hydrated: false,

      create: (email, prefs) =>
        set({
          learner: {
            email,
            ...prefs,
            xp: 0,
            lives: MAX_LIVES,
            streak: 0,
            lastLessonOn: null,
            completedCategoryIds: [],
          },
        }),

      setDailyGoal: (dailyGoal) =>
        set(({ learner }) => ({
          learner: learner ? { ...learner, dailyGoal } : null,
        })),

      loseLife: () =>
        set(({ learner }) => {
          if (!learner) return { learner };
          return {
            learner: {
              ...learner,
              lives: Math.max(0, learner.lives - 1),
              // Таймер стартует с первой потерянной жизни и дальше не сдвигается.
              nextLifeAt: learner.nextLifeAt ?? Date.now() + LIFE_REGEN_MS,
            },
          };
        }),

      regenerateLives: (now) =>
        set(({ learner }) => {
          if (!learner?.nextLifeAt || now < learner.nextLifeAt) {
            return { learner };
          }
          const earned = 1 + Math.floor((now - learner.nextLifeAt) / LIFE_REGEN_MS);
          const lives = Math.min(MAX_LIVES, learner.lives + earned);
          return {
            learner: {
              ...learner,
              lives,
              nextLifeAt:
                lives === MAX_LIVES
                  ? null
                  : learner.nextLifeAt + earned * LIFE_REGEN_MS,
            },
          };
        }),

      completeLesson: (categoryId, correctAnswers) =>
        set(({ learner }) => {
          if (!learner) return { learner };
          const today = day();
          const streak =
            learner.lastLessonOn === today
              ? learner.streak
              : learner.lastLessonOn === day(-1)
                ? learner.streak + 1
                : 1;

          return {
            learner: {
              ...learner,
              xp: learner.xp + correctAnswers * XP_PER_CORRECT_ANSWER,
              lives: MAX_LIVES,
              nextLifeAt: null,
              streak,
              lastLessonOn: today,
              completedCategoryIds: learner.completedCategoryIds.includes(
                categoryId,
              )
                ? learner.completedCategoryIds
                : [...learner.completedCategoryIds, categoryId],
            },
          };
        }),

      reset: () => set({ learner: null }),

      markHydrated: () => set({ hydrated: true }),
    }),
    {
      name: "qaztil.learner",
      storage: createJSONStorage(() => localStorage),
      partialize: ({ learner }) => ({ learner }),
      skipHydration: true,
      onRehydrateStorage: () => (state) => {
        state?.markHydrated();
      },
    },
  ),
);
