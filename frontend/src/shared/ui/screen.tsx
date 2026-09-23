import * as React from "react";

import { cn } from "@/shared/lib/cn";

/**
 * Телефонный холст макета: 390px по центру, чёрные границы по бокам
 * на широких экранах, чтобы композиция не растягивалась.
 */
function Screen({ className, children, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "relative mx-auto flex min-h-dvh w-full max-w-[430px] flex-col overflow-hidden bg-paper sm:border-x-[3px] sm:border-ink",
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export { Screen };
