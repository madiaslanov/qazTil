export type LessonProgress = {
  lesson_id: number;
  completed: boolean;
  correct: number;
  total: number;
};

export type UnitProgress = {
  unit_id: number;
  completed_lessons: number;
  total_lessons: number;
};

export type CourseProgress = {
  lessons: LessonProgress[];
  units: UnitProgress[];
};
