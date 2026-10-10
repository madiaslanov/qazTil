"use client";

import { Flame, Heart } from "lucide-react";

import { useLearnerStore, useLives } from "@/entities/learner";
import { Logo } from "@/shared/ui";

/** Шапка вкладок: логотип слева, страйк и жизни справа. */
export function TopBar() {
  const learner = useLearnerStore((state) => state.learner);
  const { lives } = useLives();

  return (
    <header className="flex h-bar shrink-0 items-center justify-between bg-surface px-gutter">
      <Logo />
      <div className="flex items-center gap-4">
        <Counter
          label="Дней подряд"
          value={learner?.streak ?? 0}
          icon={<Flame className="size-5.5 text-streak" />}
        />
        <Counter
          label="Жизни"
          value={lives}
          icon={<Heart className="size-5.5 text-lives" />}
        />
      </div>
    </header>
  );
}

function Counter({
  label,
  value,
  icon,
}: {
  label: string;
  value: number;
  icon: React.ReactNode;
}) {
  return (
    <span className="flex items-center gap-1.5" title={label}>
      {icon}
      <b className="text-body font-extrabold text-foreground">
        <span className="sr-only">{label}: </span>
        {value}
      </b>
    </span>
  );
}
