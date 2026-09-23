"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { BookOpen, Flame, Heart, Target } from "lucide-react";

import { learnerStore, useLearner, type DailyGoal } from "@/entities/learner";
import { SWAGGER_URL } from "@/shared/config/env";
import { cn } from "@/shared/lib/cn";
import { Button, Card, IconBadge, Screen } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { TopBar } from "@/widgets/top-bar";

const goals: DailyGoal[] = [5, 10, 15];

/** Профиль ученика: локальные счётчики и настройки цели. */
export function ProfilePage() {
  const router = useRouter();
  const learner = useLearner();

  function signOut() {
    learnerStore.reset();
    router.replace("/");
  }

  return (
    <Screen>
      <TopBar />

      <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
        <div>
          <p className="text-[12px] font-extrabold uppercase text-muted">
            Профиль
          </p>
          <h1 className="mt-1 text-[26px]">{learner?.email ?? "Гость"}</h1>
        </div>

        <Card className="flex flex-col gap-4.5 p-5">
          <Stat
            icon={<Flame strokeWidth={2.2} />}
            tone="bg-orange"
            label="Страйк"
            value={`${learner?.streak ?? 0}`}
          />
          <div className="h-0.5 w-full bg-ink" />
          <Stat
            icon={<Target strokeWidth={2.2} />}
            tone="bg-green"
            label="Опыт"
            value={`${learner?.xp ?? 0} XP`}
          />
          <div className="h-0.5 w-full bg-ink" />
          <Stat
            icon={<Heart strokeWidth={2.2} />}
            tone="bg-card"
            label="Жизни"
            value={`${learner?.lives ?? 0}`}
          />
        </Card>

        <Card className="flex flex-col gap-3 p-5">
          <p className="text-[12px] font-bold uppercase text-subtle">
            Ежедневная цель
          </p>
          <div className="flex gap-2.5">
            {goals.map((goal) => (
              <button
                key={goal}
                type="button"
                onClick={() => learnerStore.setDailyGoal(goal)}
                aria-pressed={learner?.dailyGoal === goal}
                className={cn(
                  "h-[52px] flex-1 rounded-field border-[3px] border-ink text-[15px] font-extrabold",
                  learner?.dailyGoal === goal
                    ? "bg-orange shadow-hard-sm"
                    : "bg-paper",
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
          Сбросить профиль
        </Button>

        <a
          href={SWAGGER_URL}
          target="_blank"
          rel="noreferrer"
          className="text-center text-[13px] text-subtle underline"
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
        <p className="text-[12px] font-bold uppercase text-subtle">{label}</p>
        <p className="text-[19px]">{value}</p>
      </div>
    </div>
  );
}
