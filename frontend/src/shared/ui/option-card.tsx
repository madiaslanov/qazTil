import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/shared/lib/cn";

const optionCardVariants = cva(
  "text-foreground transition-[background-color,border-color,box-shadow] outline-none disabled:pointer-events-none disabled:opacity-50",
  {
    variants: {
      tone: {
        // Белая карточка: уровень, мотивация.
        surface:
          "rounded-control border border-line bg-surface hover:border-primary/40 data-selected:border-2 data-selected:border-primary data-selected:bg-highlight-soft data-selected:shadow-selected",
        // Мятная плитка: ежедневная цель.
        soft: "rounded-chip border-2 border-transparent bg-glow hover:border-primary/20 data-selected:border-primary data-selected:bg-highlight-strong data-selected:shadow-glow",
      },
    },
    defaultVariants: { tone: "surface" },
  },
);

/** Выбираемая карточка: уровень, цель, ежедневная норма. */
function OptionCard({
  className,
  selected = false,
  tone,
  type = "button",
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof optionCardVariants> & { selected?: boolean }) {
  return (
    <button
      type={type}
      data-slot="option-card"
      aria-pressed={selected}
      data-selected={selected || undefined}
      className={cn(optionCardVariants({ tone }), className)}
      {...props}
    />
  );
}

export { OptionCard };
