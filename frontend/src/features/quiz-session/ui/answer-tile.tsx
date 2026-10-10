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

const tones: Record<Tone, string> = {
  idle: "border-line bg-surface hover:border-primary/30",
  selected: "border-transparent bg-highlight-soft",
  correct: "border-transparent bg-success-soft",
  wrong: "border-transparent bg-danger-soft",
};

/** Вариант ответа: буква в плашке и текст. Цвет — по результату проверки. */
export function AnswerTile({
  index,
  option,
  selected,
  correctIndex,
  disabled,
  onSelect,
}: {
  index: number;
  option: string;
  selected: number | null;
  correctIndex: number | null;
  disabled: boolean;
  onSelect: (index: number) => void;
}) {
  const tone = toneOf({ index, selected, correctIndex });

  return (
    <button
      type="button"
      onClick={() => onSelect(index)}
      disabled={disabled}
      aria-pressed={selected === index}
      className={cn(
        "flex h-answer w-full items-center gap-3 rounded-control border px-4 text-left transition-colors outline-none",
        tones[tone],
      )}
    >
      <span className="flex size-7 shrink-0 items-center justify-center rounded-option bg-line text-eyebrow font-extrabold text-foreground">
        {letters[index] ?? index + 1}
      </span>
      <span className="text-lead text-foreground">{option}</span>
    </button>
  );
}
