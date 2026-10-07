"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useQuery } from "@tanstack/react-query";

import { useLearnerStore } from "@/entities/learner";
import { sessionQueries, useSessionStore } from "@/entities/session";
import { Screen } from "@/shared/ui";

/**
 * Пускает на экраны обучения только с сессией. Токен лежит в localStorage,
 * поэтому проверка живёт на клиенте, а не в middleware.
 */
export function AuthGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const session = useSessionStore((state) => state.session);
  const sessionHydrated = useSessionStore((state) => state.hydrated);
  const setUser = useSessionStore((state) => state.setUser);
  const learner = useLearnerStore((state) => state.learner);
  const learnerHydrated = useLearnerStore((state) => state.hydrated);
  const createLearner = useLearnerStore((state) => state.create);
  const hydrated = sessionHydrated && learnerHydrated;

  // На free-тарифе Render база стирается при деплое: /auth/me отдаст 401,
  // клиент сбросит сессию, и гард отправит на вход.
  const me = useQuery({
    ...sessionQueries.me(),
    enabled: hydrated && session !== null,
  });

  useEffect(() => {
    if (me.data) setUser(me.data);
  }, [me.data, setUser]);

  useEffect(() => {
    if (!hydrated) return;
    if (!session) {
      router.replace("/");
      return;
    }
    // Локальный профиль могли стереть отдельно от сессии — заводим заново.
    if (learner?.email !== session.user.email) {
      createLearner(session.user.email, learner?.dailyGoal ?? 10);
    }
  }, [hydrated, session, learner, createLearner, router]);

  if (!hydrated || !session || learner?.email !== session.user.email) {
    return <Screen />;
  }
  return children;
}
