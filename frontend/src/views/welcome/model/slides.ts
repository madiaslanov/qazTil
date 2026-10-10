export type WelcomeSlide = {
  title: string;
  text: string;
  action: string;
  /** Буква или цифра на «плитке» иллюстрации. */
  glyph: string;
  tone: "mint" | "coral" | "sky";
};

export const WELCOME_SLIDES: WelcomeSlide[] = [
  {
    title: "Говори по-казахски каждый день.",
    text: "Короткие уроки помогают встроить язык в ритм дня и сохранять уверенный прогресс.",
    action: "Продолжить",
    glyph: "Қ",
    tone: "mint",
  },
  {
    title: "Живой язык для реальных ситуаций.",
    text: "Работа, учеба, кафе и повседневные диалоги — практика сразу становится полезной.",
    action: "Дальше",
    glyph: "10′",
    tone: "coral",
  },
  {
    title: "Твой маршрут — в твоем темпе.",
    text: "QazTil подстраивает упражнения под уровень, цель и время, которое есть сегодня.",
    action: "Настроить обучение",
    glyph: "A1",
    tone: "sky",
  },
];
