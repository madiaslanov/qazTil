import * as React from "react";

import { cn } from "@/shared/lib/cn";

/** Круглый значок с иконкой: страйк, точность и прочая статистика. */
function IconBadge({
  className,
  children,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "flex size-11 shrink-0 items-center justify-center rounded-full border-2 border-ink [&_svg]:size-5",
        className,
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export { IconBadge };
