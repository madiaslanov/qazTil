"use client";

import Link from "next/link";
import { BookOpen, Flame, Heart, Target } from "lucide-react";

import { useLearnerStore, type DailyGoal } from "@/entities/learner";
import { useSessionStore } from "@/entities/session";
import { useSignOut } from "@/features/auth";
import { SWAGGER_URL } from "@/shared/config/env";
import { cn } from "@/shared/lib/cn";
import { Button, Card, IconBadge, Screen } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { TopBar } from "@/widgets/top-bar";

const goals: DailyGoal[] = [5, 10, 15];

/** Профиль ученика: аккаунт из API, локальные счётчики и настройки цели. */
export function ProfilePage() {
  const user = useSessionStore((state) => state.session?.user);
  const learner = useLearnerStore((state) => state.learner);
  const setDailyGoal = useLearnerStore((state) => state.setDailyGoal);
  const signOut = useSignOut();

  return (
    <Screen>
      <TopBar />

      <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
        <div>
          <p className="text-eyebrow font-extrabold uppercase text-muted">
            Профиль
          </p>
          <h1 className="mt-1 text-h3">{user?.display_name}</h1>
          <p className="mt-0.5 text-body-sm text-muted">{user?.email}</p>
        </div>

        <Card className="flex flex-col gap-4.5 p-5">
          <Stat
            icon={<Flame strokeWidth={2.2} />}
            tone="bg-accent"
            label="Страйк"
            value={`${learner?.streak ?? 0}`}
          />
          <div className="h-0.5 w-full bg-primary" />
          <Stat
            icon={<Target strokeWidth={2.2} />}
            tone="bg-success"
            label="Опыт"
            value={`${learner?.xp ?? 0} XP`}
          />
          <div className="h-0.5 w-full bg-primary" />
          <Stat
            icon={<Heart strokeWidth={2.2} />}
            tone="bg-surface"
            label="Жизни"
            value={`${learner?.lives ?? 0}`}
          />
        </Card>

        <Card className="flex flex-col gap-3 p-5">
          <p className="text-eyebrow font-bold uppercase text-muted">
            Ежедневная цель
          </p>
          <div className="flex gap-2.5">
            {goals.map((goal) => (
              <button
                key={goal}
                type="button"
                onClick={() => setDailyGoal(goal)}
                aria-pressed={learner?.dailyGoal === goal}
                className={cn(
                  "h-13 flex-1 rounded-option border-3 border-primary text-body font-extrabold",
                  learner?.dailyGoal === goal
                    ? "bg-accent shadow-button"
                    : "bg-background",
                )}
              >
                {goal} мин
              </button>
            ))}
          </div>
        </Card>

        <Button asChild variant="quiet">
          <Link href="/words">
            <BookOpen className="size-5" strokeWidth={2.2} />
            Словарь
          </Link>
        </Button>

        <Button variant="ghost" size="md" onClick={signOut} className="w-full">
          Выйти
        </Button>

        <a
          href={SWAGGER_URL}
          target="_blank"
          rel="noreferrer"
          className="text-center text-caption text-muted underline"
        >
          Swagger API
        </a>
      </div>

      <BottomNav />
    </Screen>
  );
}

function Stat({
  icon,
  tone,
  label,
  value,
}: {
  icon: React.ReactNode;
  tone: string;
  label: string;
  value: string;
}) {
  return (
    <div className="flex items-center gap-3">
      <IconBadge className={tone}>{icon}</IconBadge>
      <div className="flex flex-col gap-0.5">
        <p className="text-eyebrow font-bold uppercase text-muted">{label}</p>
        <p className="text-stat">{value}</p>
      </div>
    </div>
  );
}
