import * as React from "react";

import { cn } from "@/shared/lib/cn";

/** Плотная карточка с чёрной обводкой и жёсткой тенью — базовая поверхность макета. */
function Card({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card"
      className={cn(
        "rounded-card border-3 border-primary bg-surface shadow-button",
        className,
      )}
      {...props}
    />
  );
}

export { Card };
