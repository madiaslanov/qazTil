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
        "flex h-[76px] w-full items-center gap-3.5 rounded-tile border-[3px] border-ink px-[18px] text-left transition-[transform,box-shadow]",
        tone === "idle" && "bg-card",
        tone === "selected" && "bg-orange shadow-hard-sm",
        tone === "correct" && "bg-green shadow-hard-sm",
        tone === "wrong" && "bg-[#e8705f] shadow-hard-sm",
        correctIndex === null && "active:translate-x-[2px] active:translate-y-[2px] active:shadow-none",
      )}
    >
      <span className="flex size-[30px] shrink-0 items-center justify-center rounded-field border-2 border-ink bg-paper text-[13px] font-extrabold">
        {letters[index] ?? index + 1}
      </span>
      <span className="text-[17px] text-ink">{option}</span>
    </button>
  );
}
