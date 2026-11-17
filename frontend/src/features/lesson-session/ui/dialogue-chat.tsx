import { useEffect, useRef, useState } from "react";
import { Lightbulb } from "lucide-react";

import type { DialogueExercise } from "@/entities/exercise";
import { cn } from "@/shared/lib/cn";
import { capitalize } from "@/shared/lib/text";

import type { DialogueTurn } from "../model/use-lesson-session";

function Bubble({
  speaker,
  line,
  hint,
  mine,
  dim,
}: {
  speaker: string;
  line: string;
  hint?: string;
  mine?: boolean;
  dim?: boolean;
}) {
  return (
    <div
      className={cn(
        "flex max-w-[85%] flex-col gap-1 rounded-card px-4 py-3",
        mine
          ? "self-end rounded-br-chip bg-primary text-primary-foreground"
          : "self-start rounded-bl-chip border border-line bg-surface text-foreground",
        dim && "opacity-70",
      )}
    >
      <span
        className={cn(
          "text-micro font-extrabold uppercase",
          mine ? "text-primary-foreground/70" : "text-muted",
        )}
      >
        {speaker}
      </span>
      <p lang="kk" className="text-lead">
        {line}
      </p>
      {hint && <p className="text-caption text-muted">{hint}</p>}
    </div>
  );
}

/**
 * Лента сцены: уже сыгранные реплики сверху, под ними текущая реплика
 * собеседника. Ответ ученик выбирает ниже, вариантами.
 */
export function DialogueChat({
  turns,
  step,
}: {
  turns: DialogueTurn[];
  step: DialogueExercise;
}) {
  const end = useRef<HTMLDivElement>(null);

  // Новая реплика подтягивает ленту вниз, чтобы её было видно.
  useEffect(() => {
    end.current?.scrollIntoView({ block: "nearest", behavior: "smooth" });
  }, [step.id]);

  return (
    <div className="flex flex-col gap-3" aria-label="Диалог">
      {turns.map((turn, index) => (
        <div key={index} className="flex flex-col gap-3">
          <Bubble speaker={turn.speaker} line={turn.line} dim />
          <Bubble speaker="Вы" line={capitalize(turn.reply)} mine />
        </div>
      ))}
      <Bubble speaker={step.speaker} line={step.prompt} hint={step.hint} />
      <div ref={end} />
    </div>
  );
}

/** Подсказка к сложной реплике: прячется за кнопкой, пока ученик не попросит. */
export function DialogueHelp({ text }: { text: string }) {
  const [open, setOpen] = useState(false);

  return (
    <div className="flex flex-col gap-2">
      <button
        type="button"
        onClick={() => setOpen((value) => !value)}
        aria-expanded={open}
        className="flex items-center gap-1.5 self-start text-caption font-bold text-primary outline-none"
      >
        <Lightbulb className="size-4" />
        {open ? "Скрыть подсказку" : "Подсказка"}
      </button>
      {open && (
        <p className="rounded-control bg-highlight-soft px-3.5 py-2.5 text-caption text-foreground">
          {text}
        </p>
      )}
    </div>
  );
}
