"use client";

import { useRouter } from "next/navigation";
import { useQueryClient } from "@tanstack/react-query";

import { useSessionStore } from "@/entities/session";

/**
 * Выход: забываем токен и кэш пользователя. Локальный прогресс остаётся,
 * чтобы тот же ученик после входа продолжил с того же места.
 */
export function useSignOut() {
  const router = useRouter();
  const queryClient = useQueryClient();
  const clear = useSessionStore((state) => state.clear);

  return () => {
    clear();
    queryClient.removeQueries({ queryKey: ["session"] });
    router.replace("/");
  };
}
