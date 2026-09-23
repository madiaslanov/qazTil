"use client";

import { useSyncExternalStore } from "react";

import { learnerStore } from "./store";
import type { Learner } from "./types";

export function useLearner(): Learner | null {
  return useSyncExternalStore(
    learnerStore.subscribe,
    learnerStore.get,
    learnerStore.getServerSnapshot,
  );
}
