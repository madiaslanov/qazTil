"use client";

import { useState } from "react";
import Link from "next/link";
import { ChevronLeft } from "lucide-react";

import { useQuery } from "@tanstack/react-query";

import { categoryQueries } from "@/entities/category";
import { wordQueries } from "@/entities/word";
import { cn } from "@/shared/lib/cn";
import { Button, Card, Screen, StateNote } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";

/** Словарь: слова из API с фильтром по категории. */
export function WordsPage() {
  const [categoryId, setCategoryId] = useState<number | null>(null);
  const categories = useQuery(categoryQueries.all());
  const words = useQuery(wordQueries.list(categoryId));

  return (
    <Screen>
      <header className="flex h-bar shrink-0 items-center gap-3 border-b-3 border-primary px-6">
        <Link href="/profile" aria-label="Назад в профиль">
          <ChevronLeft className="size-5.5" strokeWidth={2.4} />
        </Link>
        <p className="text-title">Словарь</p>
      </header>

      <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-5 pb-8">
        <div className="flex flex-wrap gap-2">
          <Chip active={categoryId === null} onClick={() => setCategoryId(null)}>
            Все
          </Chip>
          {categories.data?.map((category) => (
            <Chip
              key={category.id}
              active={categoryId === category.id}
              onClick={() => setCategoryId(category.id)}
            >
              {category.name_ru}
            </Chip>
          ))}
        </div>

        {words.isPending && <StateNote text="Открываем словарь…" />}
        {words.isError && (
          <StateNote
            text={words.error.message}
            action={
              <Button
                size="md"
                variant="quiet"
                className="w-auto"
                onClick={() => words.refetch()}
              >
                Повторить
              </Button>
            }
          />
        )}
        {words.data?.length === 0 && <StateNote text="В этой категории пусто." />}

        {words.data?.map((word) => (
          <Card key={word.id} className="flex flex-col gap-1 p-4">
            <div className="flex items-baseline justify-between gap-3">
              <p className="text-stat">{word.kazakh}</p>
              <p className="text-caption text-muted">{word.transcription}</p>
            </div>
            <p className="text-body text-muted">{word.russian}</p>
            {word.example_kk && (
              <p className="mt-1 text-caption text-muted">
                {word.example_kk} — {word.example_ru}
              </p>
            )}
          </Card>
        ))}
      </div>

      <BottomNav />
    </Screen>
  );
}

function Chip({
  active,
  onClick,
  children,
}: {
  active: boolean;
  onClick: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      type="button"
      onClick={onClick}
      aria-pressed={active}
      className={cn(
        "rounded-full border-2 border-primary px-3 py-1.5 text-caption font-extrabold",
        active ? "bg-accent shadow-button" : "bg-surface",
      )}
    >
      {children}
    </button>
  );
}
