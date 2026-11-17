import { queryOptions } from "@tanstack/react-query";

import type {
  Exercise,
  ExerciseAnswer,
  ExerciseResult,
} from "@/entities/exercise";
import { request } from "@/shared/api";

import type { Lesson } from "../model/types";

/** Упражнение как оно приходит: старый выбор может нести legacy_options. */
type WireExercise = Exercise & { legacy_options?: string[] };

/**
 * Выборы, засеянные до payload, приходят плоским списком строк.
 * Приводим их к обычным вариантам: текст служит и id, и подписью.
 */
function normalize(exercise: WireExercise): Exercise {
  if (exercise.kind !== "choice") return exercise;
  const options =
    exercise.options ??
    (exercise.legacy_options ?? []).map((text) => ({ id: text, text }));
  return { ...exercise, options };
}

export const lessonQueries = {
  detail: (id: number) =>
    queryOptions({
      queryKey: ["lessons", id],
      queryFn: async () => {
        const lesson = await request<
          Omit<Lesson, "exercises"> & { exercises: WireExercise[] }
        >(`/lessons/${id}`);
        return { ...lesson, exercises: lesson.exercises.map(normalize) };
      },
      // Точку продолжения читаем один раз на вход: пока идёт урок,
      // повторный запрос сдвинул бы упражнения под ногами.
      staleTime: Infinity,
      gcTime: 0,
      refetchOnWindowFocus: false,
      retry: false,
    }),
};

export const lessonApi = {
  /** Проверка одного ответа. Завершение урока приходит в ответе на последнее упражнение. */
  answer: (exerciseId: number, answer: ExerciseAnswer) =>
    request<ExerciseResult>(`/exercises/${exerciseId}/answers`, {
      method: "POST",
      body: JSON.stringify(answer),
    }),
};
