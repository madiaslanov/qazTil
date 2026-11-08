import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";

import type { DailyGoal, Level, Motivation } from "../model/types";

/** Ответ GET/PUT /me/profile. */
export type Profile = {
  id: number;
  email: string;
  display_name: string;
  created_at: string;
  daily_goal: DailyGoal;
  level: Level;
  motivation: Motivation;
  xp: number;
  lives: number;
  max_lives: number;
  next_life_at: string | null;
  streak: number;
  last_lesson_on: string | null;
};

export type ProfileUpdate = {
  display_name?: string;
  daily_goal?: DailyGoal;
  level?: Level;
  motivation?: Motivation;
};

export const profileQueries = {
  current: () =>
    queryOptions({
      queryKey: ["me", "profile"],
      queryFn: () => request<Profile>("/me/profile"),
    }),
};

export const profileApi = {
  update: (patch: ProfileUpdate) =>
    request<Profile>("/me/profile", {
      method: "PUT",
      body: JSON.stringify(patch),
    }),
};
