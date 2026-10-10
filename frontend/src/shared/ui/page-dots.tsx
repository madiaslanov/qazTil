import { cn } from "@/shared/lib/cn";

/** Индикатор страниц: активная точка вытянута в капсулу. */
function PageDots({
  count,
  active,
  onSelect,
  className,
}: {
  count: number;
  active: number;
  onSelect?: (index: number) => void;
  className?: string;
}) {
  return (
    <div className={cn("flex items-center gap-2", className)}>
      {Array.from({ length: count }, (_, index) => (
        <button
          key={index}
          type="button"
          aria-label={`Страница ${index + 1} из ${count}`}
          aria-current={index === active || undefined}
          onClick={() => onSelect?.(index)}
          className={cn(
            "h-2 rounded-full transition-[width,background-color] duration-300",
            index === active ? "w-7 bg-primary" : "w-2 bg-line",
          )}
        />
      ))}
    </div>
  );
}

export { PageDots };
