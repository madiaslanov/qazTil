import { cn } from "@/shared/lib/cn";

/** Плитка слова в сборке фразы: в банке, в ответе или уже использованная. */
export function WordChip({
  text,
  placed = false,
  used = false,
  disabled,
  onClick,
}: {
  text: string;
  placed?: boolean;
  used?: boolean;
  disabled?: boolean;
  onClick: () => void;
}) {
  return (
    <button
      type="button"
      lang="kk"
      onClick={onClick}
      disabled={disabled || used}
      aria-hidden={used || undefined}
      className={cn(
        "rounded-control px-4 py-3 text-body text-primary transition-[transform,box-shadow,opacity] outline-none active:translate-y-0.5 active:shadow-none",
        placed
          ? "border-2 border-primary bg-highlight-soft shadow-chip"
          : "border border-line bg-surface shadow-chip",
        used && "border-line bg-line opacity-55 shadow-none",
      )}
    >
      {text}
    </button>
  );
}
