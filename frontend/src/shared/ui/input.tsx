import * as React from "react";

import { cn } from "@/shared/lib/cn";

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "h-[54px] w-full min-w-0 rounded-field border-[3px] border-ink bg-card px-4 text-base text-ink outline-none",
        "placeholder:text-subtle disabled:pointer-events-none disabled:opacity-50",
        "aria-invalid:border-orange",
        className,
      )}
      {...props}
    />
  );
}

export { Input };
