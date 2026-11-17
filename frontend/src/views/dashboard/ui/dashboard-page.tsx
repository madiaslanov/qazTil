"use client";

import { useState } from "react";
import { useQuery } from "@tanstack/react-query";

import {
  levelTitle,
  profileQueries,
  useLives,
} from "@/entities/learner";
import { pathQueries, type PathLesson, type PathUnit } from "@/entities/path";
import { Button, Eyebrow, StateNote } from "@/shared/ui";
import { LessonPreview } from "@/widgets/lesson-preview";
import { NoLivesDialog } from "@/widgets/no-lives";
import { SkillTree } from "@/widgets/skill-tree";

import { PathGlows } from "./path-glows";

/** Путь обучения: юниты и уроки из GET /path, состояние считает сервер. */
export function DashboardPage() {
  const profile = useQuery(profileQueries.current());
  const { lives, pending: livesPending } = useLives();
  const path = useQuery(pathQueries.all());
  const [preview, setPreview] = useState<{
    lesson: PathLesson;
    unit: PathUnit;
  } | null>(null);
  const [noLivesOpen, setNoLivesOpen] = useState(false);

  function select(lesson: PathLesson, unit: PathUnit) {
    if (!livesPending && lives === 0) setNoLivesOpen(true);
    else setPreview({ lesson, unit });
  }

  return (
    <div className="relative flex flex-1 flex-col bg-background">
      <PathGlows />

      <div className="relative flex flex-1 flex-col px-gutter pt-6 pb-8">
        <div className="flex flex-col gap-1.25">
          {profile.data && <Eyebrow>{levelTitle(profile.data.level)}</Eyebrow>}
          <h1 className="text-h3 text-foreground">Твой путь обучения</h1>
        </div>

        <div className="mt-10 flex flex-1 flex-col">
          {path.isPending && <StateNote text="Загружаем уроки…" />}
          {path.isError && (
            <StateNote
              text={path.error.message}
              action={
                <Button
                  size="md"
                  variant="secondary"
                  className="w-auto"
                  onClick={() => path.refetch()}
                >
                  Повторить
                </Button>
              }
            />
          )}
          {path.data && path.data.units.length === 0 && (
            <StateNote text="Уроков пока нет" />
          )}
          {path.data && path.data.units.length > 0 && (
            <SkillTree units={path.data.units} onSelect={select} />
          )}
        </div>
      </div>

      <LessonPreview selected={preview} onClose={() => setPreview(null)} />
      <NoLivesDialog open={noLivesOpen} onClose={() => setNoLivesOpen(false)} />
    </div>
  );
}
