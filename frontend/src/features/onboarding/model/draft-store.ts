import { create } from "zustand";

import type { DailyGoal, Level, Motivation } from "@/entities/learner";

type OnboardingDraft = {
  level: Level;
  motivation: Motivation;
  dailyGoal: DailyGoal;
  setLevel: (level: Level) => void;
  setMotivation: (motivation: Motivation) => void;
  setDailyGoal: (dailyGoal: DailyGoal) => void;
};

/**
 * Выбор ученика между шагами онбординга. Профиль заводится только
 * после регистрации, до этого ответы живут здесь.
 */
export const useOnboardingDraft = create<OnboardingDraft>()((set) => ({
  level: "A0",
  motivation: "life",
  dailyGoal: 10,
  setLevel: (level) => set({ level }),
  setMotivation: (motivation) => set({ motivation }),
  setDailyGoal: (dailyGoal) => set({ dailyGoal }),
}));
