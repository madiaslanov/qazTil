import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";

import { cn } from "@/shared/lib/cn";

const fieldVariants = cva(
  "flex h-field w-full items-center gap-2 rounded-control bg-surface px-4 transition-colors has-disabled:opacity-50 [&_svg]:size-5 [&_svg]:shrink-0",
  {
    variants: {
      status: {
        default: "border border-line focus-within:border-primary",
        error: "border-2 border-danger [&_svg]:text-danger",
        success: "border-2 border-success [&_svg]:text-success",
      },
    },
    defaultVariants: { status: "default" },
  },
);

/** Поле ввода макета. Иконка справа подсвечивает статус проверки. */
function Input({
  className,
  type,
  status,
  icon,
  ...props
}: React.ComponentProps<"input"> &
  VariantProps<typeof fieldVariants> & { icon?: React.ReactNode }) {
  return (
    <div data-slot="field" className={cn(fieldVariants({ status }), className)}>
      <input
        type={type}
        data-slot="input"
        aria-invalid={status === "error" || undefined}
        className="h-full min-w-0 flex-1 bg-transparent text-body text-foreground outline-none placeholder:text-muted disabled:pointer-events-none"
        {...props}
      />
      {icon}
    </div>
  );
}

export { Input };
