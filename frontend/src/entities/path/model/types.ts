/** Состояние считает сервер, выводить его из порядка на клиенте не нужно. */
export type LessonState = "completed" | "current" | "locked";

export type PathLesson = {
  id: number;
  unit_id: number;
  title_kk: string;
  title_ru: string;
  description: string;
  position: number;
  state: LessonState;
};

export type PathUnit = {
  id: number;
  situation_id: number;
  title_kk: string;
  title_ru: string;
  description: string;
  position: number;
  lesson_count: number;
  lessons: PathLesson[];
};

/** GET /path. current_lesson_id — null, когда путь пройден. */
export type Path = {
  units: PathUnit[];
  current_lesson_id: number | null;
};
