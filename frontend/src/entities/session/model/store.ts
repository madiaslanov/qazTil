import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import type { Session, User } from "./types";

type SessionState = {
  session: Session | null;
  /** persist поднимается вручную в провайдере, до этого состояние пустое. */
  hydrated: boolean;
  setSession: (session: Session) => void;
  setUser: (user: User) => void;
  clear: () => void;
  markHydrated: () => void;
};

function isExpired(session: Session): boolean {
  return new Date(session.expires_at).getTime() <= Date.now();
}

export const useSessionStore = create<SessionState>()(
  persist(
    (set) => ({
      session: null,
      hydrated: false,

      setSession: (session) => set({ session }),

      setUser: (user) =>
        set(({ session }) => ({
          session: session ? { ...session, user } : null,
        })),

      clear: () => set({ session: null }),

      markHydrated: () => set({ hydrated: true }),
    }),
    {
      name: "qaztil.session",
      storage: createJSONStorage(() => localStorage),
      partialize: ({ session }) => ({ session }),
      skipHydration: true,
      onRehydrateStorage: () => (state) => {
        // Просроченный токен Go всё равно отобьёт, не держим его зря.
        if (state?.session && isExpired(state.session)) {
          state.clear();
        }
        state?.markHydrated();
      },
    },
  ),
);

/** Токен для запросов вне React, например в shared/api. */
export function getSessionToken(): string | null {
  const { session } = useSessionStore.getState();
  if (!session || isExpired(session)) return null;
  return session.token;
}
