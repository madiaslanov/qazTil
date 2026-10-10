"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { House, Trophy, User } from "lucide-react";

import { cn } from "@/shared/lib/cn";

const items = [
  { href: "/learn", label: "Путь", Icon: House },
  // Отдельного рейтинга в API нет — пока вкладка ведёт на прогресс.
  { href: "/progress", label: "Рейтинг", Icon: Trophy },
  { href: "/profile", label: "Профиль", Icon: User },
] as const;

/** Нижняя навигация на три раздела. Активный — коралловый. */
export function BottomNav() {
  const pathname = usePathname();

  return (
    <nav className="flex h-nav shrink-0 items-center justify-between bg-surface px-8.5">
      {items.map(({ href, label, Icon }) => {
        const active = pathname === href || pathname.startsWith(`${href}/`);
        return (
          <Link
            key={href}
            href={href}
            aria-current={active ? "page" : undefined}
            className={cn(
              "flex min-w-12 flex-col items-center gap-1 text-nano",
              active ? "font-bold text-accent" : "font-medium text-muted",
            )}
          >
            <Icon
              className={cn(
                "size-5.25",
                active ? "text-accent" : "text-foreground",
              )}
            />
            {label}
          </Link>
        );
      })}
    </nav>
  );
}
