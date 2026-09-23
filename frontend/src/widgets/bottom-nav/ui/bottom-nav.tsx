"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { House, Trophy, User } from "lucide-react";

import { cn } from "@/shared/lib/cn";

const items = [
  { href: "/learn", label: "Обучение", Icon: House },
  { href: "/progress", label: "Прогресс", Icon: Trophy },
  { href: "/profile", label: "Профиль", Icon: User },
] as const;

/** Нижняя навигация на три раздела, как в макете. */
export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="flex h-[84px] shrink-0 items-center justify-between border-t-[3px] border-ink bg-card px-[34px]">
      {items.map(({ href, label, Icon }) => {
        const active = pathname === href || pathname.startsWith(`${href}/`);
        return (
          <Link
            key={href}
            href={href}
            aria-label={label}
            aria-current={active ? "page" : undefined}
            className="flex flex-col items-center gap-1"
          >
            <Icon
              className={cn("size-[21px]", active ? "text-orange" : "text-ink")}
              strokeWidth={2.2}
            />
            <span className="sr-only">{label}</span>
          </Link>
        );
      })}
    </nav>
  );
}
