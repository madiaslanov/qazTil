export type LessonState = "completed" | "available" | "locked";

export type PathLesson = {
  id: string;
  title_kk: string;
  title_ru: string;
  position: number;
  /** Доступность считает сервер, выводить её из порядка не нужно. */
  state: LessonState;
  /** Номер упражнения для продолжения. Есть только у начатого урока. */
  resume_position?: number;
  correct_answers: number;
  total_answers: number;
};

export type PathUnit = {
  id: string;
  title_kk: string;
  title_ru: string;
  description: string;
  level: string;
  topic: string | null;
  position: number;
  lessons: PathLesson[];
};

/** GET /path. Юниты и уроки уже отсортированы по position. */
export type Path = { units: PathUnit[] };
