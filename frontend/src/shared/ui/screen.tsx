import * as React from "react";

import { cn } from "@/shared/lib/cn";

/**
 * Телефонный холст макета: по центру и не шире телефона,
 * чтобы на десктопе композиция не растягивалась.
 */
function Screen({ className, children, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "relative mx-auto flex min-h-dvh w-full max-w-phone flex-col overflow-hidden bg-background",
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export { Screen };
