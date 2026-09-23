import { request } from "@/shared/api";
import type { AnswerResult, Quiz } from "../model/types";

export const quizApi = {
  start: (categoryId: number, size: number) =>
    request<Quiz>("/quizzes", {
      method: "POST",
      body: JSON.stringify({ category_id: categoryId, size }),
    }),
  get: (quizId: number) => request<Quiz>(`/quizzes/${quizId}`),
  answer: (quizId: number, questionId: number, selectedIndex: number) =>
    request<AnswerResult>(`/quizzes/${quizId}/answers`, {
      method: "POST",
      body: JSON.stringify({
        question_id: questionId,
        selected_index: selectedIndex,
      }),
    }),
};
