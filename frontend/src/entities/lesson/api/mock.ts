import type { ExerciseResult, LessonAnswer } from "@/entities/exercise/@x/lesson";
import type { Lesson, LessonComplete, LessonStart } from "../model/types";

export const mockLesson: Lesson = {
  id: "2",
  unit_id: "1",
  title_kk: "Тапсырыс",
  title_ru: "Заказ",
  state: "available",
  resume_position: 0,
  hearts_left: 3,
  exercises: [
    {
      id: "12",
      kind: "choice",
      position: 1,
      prompt: "Как сказать «спасибо»?",
      options: [
        { id: "o1", text: "Рақмет" },
        { id: "o2", text: "Кешіріңіз" },
        { id: "o3", text: "Сәлем" },
        { id: "o4", text: "Сау бол" },
      ],
    },
    {
      id: "13",
      kind: "match",
      position: 2,
      pairs_left: [
        { id: "l1", text: "Рақмет" },
        { id: "l2", text: "Сәлем" },
      ],
      pairs_right: [
        { id: "r1", text: "Привет" },
        { id: "r2", text: "Спасибо" },
      ],
    },
    {
      id: "14",
      kind: "assemble",
      position: 3,
      prompt: "Соберите: «Я пью кофе»",
      tiles: [
        { id: "t1", text: "кофе" },
        { id: "t2", text: "Мен" },
        { id: "t3", text: "ішемін" },
        { id: "t4", text: "шай" },
      ],
    },
  ],
};

export const mockStart: LessonStart = { resume_position: 0, hearts_left: 3, attempt: 1 };

export const mockComplete: LessonComplete = {
  xp_earned: 10,
  xp_total: 130,
  streak_current: 5,
  streak_increased: true,
  today: { xp_earned: 10, goal_xp: 10, goal_met: true },
  next_lesson_id: "3",
};

// Ключи ответов — только в заглушке, настоящий сервер их не отдаёт.
let hearts = mockStart.hearts_left;

export function mockAnswer(exerciseId: string, answer: LessonAnswer): ExerciseResult {
  let correct = false;
  let correct_answer = "";
  if (exerciseId === "12" && answer.kind === "choice") {
    correct = answer.option_id === "o1";
    correct_answer = "Рақмет";
  } else if (exerciseId === "13" && answer.kind === "match") {
    const expected = new Map([["l1", "r2"], ["l2", "r1"]]);
    correct =
      answer.pairs.length === expected.size &&
      answer.pairs.every(([l, r]) => expected.get(l) === r);
    correct_answer = "Рақмет — Спасибо · Сәлем — Привет";
  } else if (exerciseId === "14" && answer.kind === "assemble") {
    correct = answer.tile_ids.join(",") === "t2,t1,t3";
    correct_answer = "Мен кофе ішемін";
  }
  if (!correct) hearts = Math.max(0, hearts - 1);
  return { correct, correct_answer, hearts_left: hearts, lesson_failed: hearts === 0 };
}

export function resetMockHearts() {
  hearts = mockStart.hearts_left;
}
