"use client";

import { useEffect, useRef, useState } from "react";

import { ApiError } from "@/shared/api";

type Result<T> = {
  key: string;
  attempt: number;
  data: T | null;
  error: string | null;
};

/**
 * Загрузка данных из Go API. `key` описывает запрос: как только он меняется,
 * прошлый ответ считается устаревшим и экран снова показывает загрузку.
 */
export function useRequest<T>(run: () => Promise<T>, key: string) {
  const [attempt, setAttempt] = useState(0);
  const [result, setResult] = useState<Result<T> | null>(null);

  const runRef = useRef(run);
  useEffect(() => {
    runRef.current = run;
  });

  useEffect(() => {
    let cancelled = false;
    runRef
      .current()
      .then((data) => {
        if (!cancelled) setResult({ key, attempt, data, error: null });
      })
      .catch((cause: unknown) => {
        if (cancelled) return;
        setResult({
          key,
          attempt,
          data: null,
          error:
            cause instanceof ApiError ? cause.message : "не удалось загрузить",
        });
      });
    return () => {
      cancelled = true;
    };
  }, [key, attempt]);

  const fresh =
    result && result.key === key && result.attempt === attempt ? result : null;

  return {
    data: fresh?.data ?? null,
    error: fresh?.error ?? null,
    loading: fresh === null,
    retry: () => setAttempt((value) => value + 1),
  };
}
