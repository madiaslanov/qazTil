/** Шкала курса: после стартового уровня модули идут по порядку. */
export const COURSE_LEVELS = ["A0", "A1", "A2", "B1", "B2"] as const;

export type CourseLevel = (typeof COURSE_LEVELS)[number];

const titles: Record<CourseLevel, string> = {
  A0: "с нуля",
  A1: "начальный",
  A2: "элементарный",
  B1: "средний",
  B2: "выше среднего",
};

/** «A1 · начальный» */
export function levelTitle(level: CourseLevel) {
  return `${level} · ${titles[level]}`;
}

/** Уровень модуля по его номеру от стартового уровня ученика. */
export function levelAt(start: CourseLevel, offset: number): CourseLevel {
  const index = COURSE_LEVELS.indexOf(start) + offset;
  return COURSE_LEVELS[Math.min(index, COURSE_LEVELS.length - 1)];
}
