import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { Slot } from "radix-ui";

import { cn } from "@/shared/lib/cn";

const buttonVariants = cva(
  "inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap border-[3px] border-ink transition-[transform,box-shadow,background-color] outline-none disabled:pointer-events-none disabled:opacity-50 active:translate-x-[2px] active:translate-y-[2px] active:shadow-none [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        primary: "bg-ink text-white shadow-hard",
        accent: "bg-orange text-ink shadow-hard",
        success: "bg-green text-ink shadow-hard",
        quiet: "bg-card text-ink hover:bg-paper",
        ghost: "border-transparent bg-transparent text-ink shadow-none active:translate-none",
      },
      size: {
        lg: "h-[58px] w-full rounded-tile px-6 text-base",
        md: "h-[46px] rounded-tile px-5 text-[15px]",
        sm: "h-[36px] rounded-field px-3 text-[13px] font-bold",
        icon: "size-11 rounded-full",
      },
    },
    defaultVariants: {
      variant: "primary",
      size: "lg",
    },
  },
);

function Button({
  className,
  variant,
  size,
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean;
  }) {
  const Comp = asChild ? Slot.Root : "button";

  return (
    <Comp
      data-slot="button"
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  );
}

export { Button, buttonVariants };
