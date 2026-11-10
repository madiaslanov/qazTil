export type StartLevel = "A0" | "A1" | "A2" | "B1" | "B2";
export type ProfileMotivation = "work" | "study" | "life";

/** GET /me/profile */
export type Profile = {
  start_level: StartLevel;
  motivation: ProfileMotivation;
  daily_goal_xp: number;
  timezone: string;
  onboarded: boolean;
};

export type TodayProgress = {
  xp_earned: number;
  goal_xp: number;
  goal_met: boolean;
};

/** GET /me/stats. today.date — местная дата ученика, считает сервер. */
export type Stats = {
  xp_total: number;
  streak_current: number;
  streak_best: number;
  last_lesson_date: string | null;
  today: TodayProgress & { date: string };
};
