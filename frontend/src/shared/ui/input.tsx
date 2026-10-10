import * as React from "react";

import { cn } from "@/shared/lib/cn";

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "h-field w-full min-w-0 rounded-option border-3 border-primary bg-surface px-4 text-body text-foreground outline-none",
        "placeholder:text-muted disabled:pointer-events-none disabled:opacity-50",
        "aria-invalid:border-accent",
        className,
      )}
      {...props}
    />
  );
}

export { Input };
