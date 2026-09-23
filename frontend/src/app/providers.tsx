"use client";

import { useEffect } from "react";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { useLearnerStore } from "@/entities/learner";

function makeQueryClient() {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Словарь и категории меняются редко: держим ответ свежим минуту,
        // чтобы переходы между экранами не дёргали Go лишний раз.
        staleTime: 60_000,
        gcTime: 5 * 60_000,
        retry: 1,
        refetchOnWindowFocus: false,
      },
    },
  });
}

let browserQueryClient: QueryClient | undefined;

function getQueryClient() {
  if (typeof window === "undefined") {
    return makeQueryClient();
  }
  browserQueryClient ??= makeQueryClient();
  return browserQueryClient;
}

export function Providers({ children }: { children: React.ReactNode }) {
  const queryClient = getQueryClient();

  // persist поднимаем вручную, чтобы разметка сервера и клиента совпадала.
  useEffect(() => {
    void useLearnerStore.persist.rehydrate();
  }, []);

  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
