import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/shared/lib/cn";

const fieldHintVariants = cva("text-eyebrow font-semibold", {
  variants: {
    tone: {
      muted: "text-muted",
      error: "text-danger",
      success: "text-success",
    },
  },
  defaultVariants: { tone: "muted" },
});

/** Подсказка под полем: ошибка проверки или подтверждение. */
function FieldHint({
  className,
  tone,
  ...props
}: React.ComponentProps<"p"> & VariantProps<typeof fieldHintVariants>) {
  return (
    <p
      data-slot="field-hint"
      role={tone === "error" ? "alert" : undefined}
      className={cn(fieldHintVariants({ tone }), className)}
      {...props}
    />
  );
}

export { FieldHint };
