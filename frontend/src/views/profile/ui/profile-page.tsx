"use client";

import Image from "next/image";
import Link from "next/link";
import { BookOpen, Flame, Heart, Target } from "lucide-react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import {
  profileApi,
  profileQueries,
  useLives,
  type DailyGoal,
} from "@/entities/learner";
import { useSignOut } from "@/features/auth";
import { SWAGGER_URL } from "@/shared/config/env";
import { cn } from "@/shared/lib/cn";
import { Button, Card, IconBadge } from "@/shared/ui";

const goals: DailyGoal[] = [5, 10, 15];

/** Профиль ученика: аккаунт, счётчики и цель с сервера. */
export function ProfilePage() {
  const queryClient = useQueryClient();
  const profile = useQuery(profileQueries.current());
  const { lives, streak, xp, pending } = useLives();
  const signOut = useSignOut();

  const updateGoal = useMutation({
    mutationFn: (daily_goal: DailyGoal) => profileApi.update({ daily_goal }),
    onSuccess: (data) => {
      queryClient.setQueryData(profileQueries.current().queryKey, data);
    },
  });

  const dailyGoal = profile.data?.daily_goal;
  const busy = profile.isPending || updateGoal.isPending;

  return (
    <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
      <div>
        <p className="text-eyebrow font-extrabold uppercase text-muted">
          Профиль
        </p>
        <h1 className="mt-1 text-h3">
          {profile.data?.display_name ?? "…"}
        </h1>
        <p className="mt-0.5 text-body-sm text-muted">
          {profile.data?.email ?? "…"}
        </p>
      </div>

      <Card className="flex flex-col gap-4.5 p-5">
        <Stat
          icon={<Flame strokeWidth={2.2} />}
          tone="bg-accent"
          label="Страйк"
          value={pending ? "…" : `${streak}`}
        />
        <div className="h-0.5 w-full bg-primary" />
        <Stat
          icon={<Target strokeWidth={2.2} />}
          tone="bg-success"
          label="Опыт"
          value={pending ? "…" : `${xp} XP`}
        />
        <div className="h-0.5 w-full bg-primary" />
        <Stat
          icon={<Heart strokeWidth={2.2} className="text-lives" />}
          tone="bg-surface"
          label="Жизни"
          value={pending || lives === undefined ? "…" : `${lives}`}
          labelSticker={
            <Image
              src="/stickers/heart.jpg"
              alt=""
              width={22}
              height={22}
              className="size-5.5 object-contain drop-shadow-sm"
              aria-hidden
            />
          }
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
              disabled={busy}
              onClick={() => updateGoal.mutate(goal)}
              aria-pressed={dailyGoal === goal}
              className={cn(
                "h-13 flex-1 rounded-option border-3 border-primary text-body font-extrabold",
                dailyGoal === goal
                  ? "bg-accent shadow-button"
                  : "bg-background",
              )}
            >
              {goal} мин
            </button>
          ))}
        </div>
      </Card>

      <Button asChild variant="secondary">
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
  );
}

function Stat({
  icon,
  tone,
  label,
  value,
  labelSticker,
}: {
  icon: React.ReactNode;
  tone: string;
  label: string;
  value: string;
  labelSticker?: React.ReactNode;
}) {
  return (
    <div className="flex items-center gap-3">
      <IconBadge className={tone}>{icon}</IconBadge>
      <div className="flex flex-col gap-0.5">
        <p className="flex items-center gap-1.5 text-eyebrow font-bold uppercase text-muted">
          {label}
          {labelSticker}
        </p>
        <p className="text-stat">{value}</p>
      </div>
    </div>
  );
}
