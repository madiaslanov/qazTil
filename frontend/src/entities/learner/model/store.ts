import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import {
  MAX_LIVES,
  XP_PER_CORRECT_ANSWER,
  type DailyGoal,
  type Learner,
} from "./types";

type LearnerState = {
  learner: Learner | null;
  /** persist поднимается вручную в провайдере, до этого состояние пустое. */
  hydrated: boolean;
  create: (email: string, dailyGoal: DailyGoal) => void;
  setDailyGoal: (dailyGoal: DailyGoal) => void;
  loseLife: () => void;
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

      create: (email, dailyGoal) =>
        set({
          learner: {
            email,
            dailyGoal,
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
        set(({ learner }) => ({
          learner: learner
            ? { ...learner, lives: Math.max(0, learner.lives - 1) }
            : null,
        })),

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
