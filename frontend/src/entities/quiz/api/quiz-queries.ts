import { queryOptions } from "@tanstack/react-query";

import { request } from "@/shared/api";
import type { AnswerResult, Quiz } from "../model/types";

/**
 * Квиз заводится POST-запросом, поэтому запрос держится в кэше навсегда:
 * повторный вход на экран урока не должен плодить новые квизы на сервере.
 */
export const quizQueries = {
  session: (categoryId: number, size: number) =>
    queryOptions({
      queryKey: ["quiz", "session", categoryId, size],
      queryFn: () =>
        request<Quiz>("/quizzes", {
          method: "POST",
          body: JSON.stringify({ category_id: categoryId, size }),
        }),
      staleTime: Infinity,
      gcTime: 0,
      retry: false,
      refetchOnMount: false,
      refetchOnReconnect: false,
      refetchOnWindowFocus: false,
    }),
};

export const quizApi = {
  answer: (quizId: number, questionId: number, selectedIndex: number) =>
    request<AnswerResult>(`/quizzes/${quizId}/answers`, {
      method: "POST",
      body: JSON.stringify({
        question_id: questionId,
        selected_index: selectedIndex,
      }),
    }),
};
