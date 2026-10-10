"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";

import { useSessionStore } from "@/entities/session";
import { Screen } from "@/shared/ui";

/** Экраны для гостей: с сессией отсюда уводит сразу на путь обучения. */
export function GuestGuard({ children }: { children: React.ReactNode }) {
  const router = useRouter();
  const session = useSessionStore((state) => state.session);
  const hydrated = useSessionStore((state) => state.hydrated);

  useEffect(() => {
    if (session) router.replace("/learn");
  }, [session, router]);

  // Пока persist не поднялся, не мигаем гостевым экраном тому, кто уже вошёл.
  if (!hydrated || session) {
    return <Screen />;
  }
  return children;
}
