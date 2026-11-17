import { useEffect, useRef } from "react";

import type { DialogueExercise } from "@/entities/exercise";
import { cn } from "@/shared/lib/cn";

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
          <Bubble speaker="Ты" line={turn.reply} mine />
        </div>
      ))}
      <Bubble speaker={step.speaker} line={step.prompt} hint={step.hint} />
      <div ref={end} />
    </div>
  );
}
