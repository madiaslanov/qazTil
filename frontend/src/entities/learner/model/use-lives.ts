"use client";

import { useEffect, useRef, useState } from "react";
import { useQuery } from "@tanstack/react-query";

import { profileQueries } from "../api/profile-queries";

/** Жизни с сервера. Таймер идёт по next_life_at, восстановление считает API. */
export function useLives() {
  const profile = useQuery(profileQueries.current());
  const [now, setNow] = useState(() => Date.now());
  const nextLifeAt = profile.data?.next_life_at
    ? Date.parse(profile.data.next_life_at)
    : null;
  const refreshed = useRef<number | null>(null);

  useEffect(() => {
    if (nextLifeAt === null) return;
    const timer = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(timer);
  }, [nextLifeAt]);

  useEffect(() => {
    if (nextLifeAt === null || now < nextLifeAt) return;
    if (refreshed.current === nextLifeAt) return;
    refreshed.current = nextLifeAt;
    void profile.refetch();
  }, [nextLifeAt, now, profile]);

  return {
    lives: profile.data?.lives,
    pending: profile.isPending,
    msToNextLife:
      nextLifeAt === null ? null : Math.max(0, nextLifeAt - now),
    xp: profile.data?.xp ?? 0,
    streak: profile.data?.streak ?? 0,
  };
}
