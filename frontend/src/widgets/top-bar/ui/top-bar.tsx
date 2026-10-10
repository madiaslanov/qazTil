"use client";

import { Flame, Heart } from "lucide-react";

import { useLearnerStore } from "@/entities/learner";

/** Шапка макета: логотип слева, страйк и жизни справа. */
export function TopBar() {
  const learner = useLearnerStore((state) => state.learner);

  return (
    <header className="flex h-bar shrink-0 items-center justify-between border-b-3 border-primary px-6">
      <p className="text-title text-foreground">QazTil</p>
      <div className="flex items-center gap-4">
        <span className="flex items-center gap-1.5">
          <Flame className="size-5.5 text-accent" strokeWidth={2.2} />
          <b className="text-body font-extrabold">{learner?.streak ?? 0}</b>
        </span>
        <span className="flex items-center gap-1.5">
          <Heart className="size-5.5 text-lives" strokeWidth={2.2} />
          <b className="text-body font-extrabold">{learner?.lives ?? 0}</b>
        </span>
      </div>
    </header>
  );
}
