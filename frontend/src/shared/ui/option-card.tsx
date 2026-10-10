import * as React from "react";

import { cn } from "@/shared/lib/cn";

/**
 * Выбираемая карточка: уровень, цель, ежедневная норма.
 * Выбранная — мятная, с тёмной обводкой и свечением.
 */
function OptionCard({
  className,
  selected = false,
  type = "button",
  ...props
}: React.ComponentProps<"button"> & { selected?: boolean }) {
  return (
    <button
      type={type}
      data-slot="option-card"
      aria-pressed={selected}
      data-selected={selected || undefined}
      className={cn(
        "rounded-control border border-line bg-surface text-foreground transition-[background-color,border-color,box-shadow] outline-none",
        "hover:border-primary/40 disabled:pointer-events-none disabled:opacity-50",
        "data-selected:border-2 data-selected:border-primary data-selected:bg-highlight-soft data-selected:shadow-selected",
        className,
      )}
      {...props}
    />
  );
}

export { OptionCard };
