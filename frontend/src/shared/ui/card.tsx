import * as React from "react";

import { cn } from "@/shared/lib/cn";

/** Белая карточка с тонкой обводкой — базовая поверхность макета. */
function Card({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      data-slot="card"
      className={cn(
        "rounded-card border border-line bg-surface",
        className,
      )}
      {...props}
    />
  );
}

export { Card };
