import * as React from "react";

import { cn } from "@/shared/lib/cn";

/** Плотная карточка с чёрной обводкой и жёсткой тенью — базовая поверхность макета. */
function Card({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card"
      className={cn(
        "rounded-card border-[3px] border-ink bg-card shadow-hard",
        className,
      )}
      {...props}
    />
  );
}

export { Card };
