import { cn } from "@/shared/lib/cn";

const letters = ["A", "B", "C", "D", "E", "F"] as const;

type Tone = "idle" | "selected" | "correct" | "wrong";

function toneOf({
  index,
  selected,
  correctIndex,
}: {
  index: number;
  selected: number | null;
  correctIndex: number | null;
}): Tone {
  if (correctIndex === null) return selected === index ? "selected" : "idle";
  if (index === correctIndex) return "correct";
  if (index === selected) return "wrong";
  return "idle";
}

/** Вариант ответа: буква в квадрате плюс текст, цвет зависит от проверки. */
export function AnswerTile({
  index,
  option,
  selected,
  correctIndex,
  onSelect,
}: {
  index: number;
  option: string;
  selected: number | null;
  correctIndex: number | null;
  onSelect: (index: number) => void;
}) {
  const tone = toneOf({ index, selected, correctIndex });

  return (
    <button
      type="button"
      onClick={() => onSelect(index)}
      disabled={correctIndex !== null}
      aria-pressed={selected === index}
      className={cn(
        "flex h-19 w-full items-center gap-3.5 rounded-control border-3 border-primary px-4.5 text-left transition-[transform,box-shadow]",
        tone === "idle" && "bg-surface",
        tone === "selected" && "bg-accent shadow-button",
        tone === "correct" && "bg-success shadow-button",
        tone === "wrong" && "bg-danger shadow-button",
        correctIndex === null && "active:translate-x-0.5 active:translate-y-0.5 active:shadow-none",
      )}
    >
      <span className="flex size-7.5 shrink-0 items-center justify-center rounded-option border-2 border-primary bg-background text-caption font-extrabold">
        {letters[index] ?? index + 1}
      </span>
      <span className="text-lead text-foreground">{option}</span>
    </button>
  );
}
