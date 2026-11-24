"use client";

import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";

import { useSessionStore } from "@/entities/session";

/**
 * Выход: забываем токен и весь кэш запросов. Путь, прогресс и уроки
 * у каждого ученика свои, и следующий аккаунт не должен увидеть чужие.
 * Локальные категории остаются: при входе другим email они пересоздаются.
 */
export function useSignOut() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const clear = useSessionStore((state) => state.clear);

  return () => {
    clear();
    queryClient.clear();
    router.replace("/");
  };
}
