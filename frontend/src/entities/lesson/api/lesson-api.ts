import { queryOptions } from "@tanstack/react-query";

import type { ExerciseResult, LessonAnswer } from "@/entities/exercise/@x/lesson";
import { request } from "@/shared/api";
import { USE_MOCKS } from "@/shared/config/env";
import type { Lesson, LessonComplete, LessonStart } from "../model/types";
import {
  mockAnswer,
  mockComplete,
  mockLesson,
  mockStart,
  resetMockHearts,
} from "./mock";

export const lessonQueries = {
  detail: (id: string) =>
    queryOptions({
      queryKey: ["lessons", id],
      queryFn: () =>
        USE_MOCKS
          ? Promise.resolve({ ...mockLesson, id })
          : request<Lesson>(`/lessons/${id}`),
    }),
};

/** Вход в урок: именно этот запрос создаёт запись прогресса на сервере. */
export function startLesson(id: string): Promise<LessonStart> {
  if (USE_MOCKS) {
    resetMockHearts();
    return Promise.resolve(mockStart);
  }
  return request<LessonStart>(`/lessons/${id}/start`, { method: "POST" });
}

/** Проверка одного ответа. correct_answer приходит всегда. */
export function submitAnswer(
  exerciseId: string,
  answer: LessonAnswer,
): Promise<ExerciseResult> {
  if (USE_MOCKS) return Promise.resolve(mockAnswer(exerciseId, answer));
  return request<ExerciseResult>(`/exercises/${exerciseId}/answers`, {
    method: "POST",
    body: JSON.stringify(answer),
  });
}

/** Завершение урока — после последнего упражнения. */
export function completeLesson(id: string): Promise<LessonComplete> {
  if (USE_MOCKS) return Promise.resolve(mockComplete);
  return request<LessonComplete>(`/lessons/${id}/complete`, { method: "POST" });
}
