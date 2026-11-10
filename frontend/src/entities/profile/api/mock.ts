import type { Profile, Stats } from "../model/types";

export const mockProfile: Profile = {
  start_level: "A0",
  motivation: "life",
  daily_goal_xp: 10,
  timezone: "Asia/Almaty",
  onboarded: true,
};

export const mockStats: Stats = {
  xp_total: 120,
  streak_current: 4,
  streak_best: 11,
  last_lesson_date: "2026-10-27",
  today: { date: "2026-10-28", xp_earned: 5, goal_xp: 10, goal_met: false },
};
