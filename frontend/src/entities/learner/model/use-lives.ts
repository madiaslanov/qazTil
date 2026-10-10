"use client";

import { useEffect, useState } from "react";

import { useLearnerStore } from "./store";

/** Жизни с восстановлением по таймеру и отсчёт до следующей. */
export function useLives() {
  const lives = useLearnerStore((state) => state.learner?.lives ?? 0);
  const nextLifeAt = useLearnerStore((state) => state.learner?.nextLifeAt ?? null);
  const regenerate = useLearnerStore((state) => state.regenerateLives);
  const [now, setNow] = useState(() => Date.now());

  useEffect(() => {
    if (nextLifeAt === null) return;
    const tick = () => {
      const current = Date.now();
      setNow(current);
      regenerate(current);
    };
    tick();
    const timer = setInterval(tick, 1000);
    return () => clearInterval(timer);
  }, [nextLifeAt, regenerate]);

  return {
    lives,
    msToNextLife: nextLifeAt === null ? null : Math.max(0, nextLifeAt - now),
  };
}
