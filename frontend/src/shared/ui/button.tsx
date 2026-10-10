import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { Slot } from "radix-ui";

import { cn } from "@/shared/lib/cn";

const buttonVariants = cva(
  "inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap border-3 border-primary transition-[transform,box-shadow,background-color] outline-none disabled:pointer-events-none disabled:opacity-50 active:translate-x-0.5 active:translate-y-0.5 active:shadow-none [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        primary: "bg-primary text-primary-foreground shadow-button",
        accent: "bg-accent text-foreground shadow-button",
        success: "bg-success text-foreground shadow-button",
        quiet: "bg-surface text-foreground hover:bg-background",
        ghost: "border-transparent bg-transparent text-foreground shadow-none active:translate-none",
      },
      size: {
        lg: "h-control w-full rounded-control px-6 text-lead",
        md: "h-11.5 rounded-control px-5 text-body",
        sm: "h-9 rounded-option px-3 text-caption font-bold",
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
