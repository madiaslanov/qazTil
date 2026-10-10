import * as React from "react";

import { cn } from "@/shared/lib/cn";

/** Короткая коралловая надпись капсом над заголовком: «A1 · начальный». */
function Eyebrow({ className, ...props }: React.ComponentProps<"p">) {
  return (
    <p
      data-slot="eyebrow"
      className={cn(
        "text-eyebrow font-extrabold text-accent uppercase",
        className,
      )}
      {...props}
    />
  );
}

export { Eyebrow };
