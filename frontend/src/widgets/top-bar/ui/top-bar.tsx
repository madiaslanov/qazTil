"use client";

import { Flame, Heart } from "lucide-react";

import { useLives } from "@/entities/learner";
import { Logo } from "@/shared/ui";

/** Шапка вкладок: логотип слева, страйк и жизни справа. */
export function TopBar() {
  const { lives, streak, pending } = useLives();

  return (
    <header className="flex h-bar shrink-0 items-center justify-between bg-surface px-gutter">
      <Logo />
      <div className="flex items-center gap-4">
        <Counter
          label="Дней подряд"
          value={pending ? "…" : streak}
          icon={<Flame className="size-5.5 text-streak" />}
        />
        <Counter
          label="Жизни"
          value={pending || lives === undefined ? "…" : lives}
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
  value: number | string;
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
