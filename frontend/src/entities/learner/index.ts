export type {
  DailyGoal,
  Learner,
  LearnerPrefs,
  Level,
  Motivation,
} from "./model/types";
export { LIFE_REGEN_MS, MAX_LIVES, XP_PER_CORRECT_ANSWER } from "./model/types";
export {
  COURSE_LEVELS,
  levelAt,
  levelTitle,
  type CourseLevel,
} from "./model/levels";
export { useLearnerStore } from "./model/store";
export { useLives } from "./model/use-lives";
