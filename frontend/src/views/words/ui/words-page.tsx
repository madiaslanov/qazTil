"use client";

import { useState } from "react";

import { useQuery } from "@tanstack/react-query";

import { categoryQueries } from "@/entities/category";
import { wordQueries } from "@/entities/word";
import { cn } from "@/shared/lib/cn";
import { BackLink, Button, Card, StateNote } from "@/shared/ui";

/** Словарь: слова из API с фильтром по категории. */
export function WordsPage() {
  const [categoryId, setCategoryId] = useState<number | null>(null);
  const categories = useQuery(categoryQueries.all());
  const words = useQuery(wordQueries.list(categoryId));

  return (
    <div className="flex flex-1 flex-col gap-4 px-gutter pt-6 pb-8">
      <div className="flex items-center gap-3">
        <BackLink href="/profile" />
        <h1 className="text-h3">Словарь</h1>
      </div>

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
              variant="secondary"
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
