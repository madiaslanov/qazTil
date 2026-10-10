"use client";

import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { Progress as ProgressPrimitive } from "radix-ui";

import { cn } from "@/shared/lib/cn";

const indicatorVariants = cva(
  "h-full w-full flex-1 rounded-full transition-transform duration-300",
  {
    variants: {
      tone: {
        // Прогресс урока в шапке.
        primary: "bg-primary",
        // Прогресс темы и прочая статистика.
        highlight: "bg-highlight",
      },
    },
    defaultVariants: { tone: "primary" },
  },
);

function Progress({
  className,
  value,
  tone,
  ...props
}: React.ComponentProps<typeof ProgressPrimitive.Root> &
  VariantProps<typeof indicatorVariants>) {
  return (
    <ProgressPrimitive.Root
      data-slot="progress"
      value={value}
      className={cn(
        "relative h-1.75 w-full overflow-hidden rounded-full bg-line",
        className,
      )}
      {...props}
    >
      <ProgressPrimitive.Indicator
        data-slot="progress-indicator"
        className={indicatorVariants({ tone })}
        style={{ transform: `translateX(-${100 - (value || 0)}%)` }}
      />
    </ProgressPrimitive.Root>
  );
}

export { Progress };
