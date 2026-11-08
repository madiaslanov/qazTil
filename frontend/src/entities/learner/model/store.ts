import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import type { Learner } from "./types";

type LearnerState = {
  learner: Learner | null;
  /** persist поднимается вручную в провайдере, до этого состояние пустое. */
  hydrated: boolean;
  create: (email: string) => void;
  completeLesson: (categoryId: number) => void;
  reset: () => void;
  markHydrated: () => void;
};

/**
 * Локально держим только пройденные категории.
 * Prefs и геймификация — в /me/profile.
 */
export const useLearnerStore = create<LearnerState>()(
  persist(
    (set) => ({
      learner: null,
      hydrated: false,

      create: (email) =>
        set({
          learner: {
            email,
            completedCategoryIds: [],
          },
        }),

      completeLesson: (categoryId) =>
        set(({ learner }) => {
          if (!learner) return { learner };
          return {
            learner: {
              ...learner,
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
      // Старые профили хранили xp/lives/prefs — оставляем только email и категории.
      merge: (persisted, current) => {
        const raw = persisted as { learner?: Record<string, unknown> } | null;
        const previous = raw?.learner;
        if (!previous || typeof previous.email !== "string") {
          return current;
        }
        const ids = Array.isArray(previous.completedCategoryIds)
          ? previous.completedCategoryIds.filter(
              (id): id is number => typeof id === "number",
            )
          : [];
        return {
          ...current,
          learner: { email: previous.email, completedCategoryIds: ids },
        };
      },
      onRehydrateStorage: () => (state) => {
        state?.markHydrated();
      },
    },
  ),
);
