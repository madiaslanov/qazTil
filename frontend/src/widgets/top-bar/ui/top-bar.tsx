"use client";

import { Flame, Heart } from "lucide-react";

import { useLearnerStore } from "@/entities/learner";

/** Шапка макета: логотип слева, страйк и жизни справа. */
export function TopBar() {
  const learner = useLearnerStore((state) => state.learner);

  return (
    <header className="flex h-[82px] shrink-0 items-center justify-between border-b-[3px] border-ink px-6">
      <p className="text-[23px] text-ink">QazTil</p>
      <div className="flex items-center gap-4">
        <span className="flex items-center gap-1.5">
          <Flame className="size-[22px] text-orange" strokeWidth={2.2} />
          <b className="text-[15px] font-extrabold">{learner?.streak ?? 0}</b>
        </span>
        <span className="flex items-center gap-1.5">
          <Heart className="size-[22px] text-[#e05a5a]" strokeWidth={2.2} />
          <b className="text-[15px] font-extrabold">{learner?.lives ?? 0}</b>
        </span>
      </div>
    </header>
  );
}
