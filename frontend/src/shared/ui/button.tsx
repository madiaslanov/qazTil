import * as React from "react";
import { cva, type VariantProps } from "class-variance-authority";
import { Slot } from "radix-ui";

import { cn } from "@/shared/lib/cn";

const buttonVariants = cva(
  "inline-flex shrink-0 items-center justify-center gap-2 whitespace-nowrap font-bold transition-[transform,box-shadow,background-color] outline-none disabled:pointer-events-none disabled:opacity-50 [&_svg]:pointer-events-none [&_svg]:shrink-0",
  {
    variants: {
      variant: {
        // Основное действие экрана: тёмная кнопка с «подставкой» снизу.
        primary:
          "border-2 border-primary bg-primary text-primary-foreground shadow-button active:translate-y-0.5 active:shadow-none",
        // Альтернатива рядом с основной: «Продолжить как гость».
        secondary:
          "border border-line bg-surface text-primary hover:bg-surface-muted",
        // Круглые кнопки в шапке: «назад», «закрыть».
        surface: "bg-surface text-foreground hover:bg-surface-muted",
        // Текстовые ссылки-действия: «Забыли пароль?», «Вернуться на путь».
        ghost: "bg-transparent text-muted hover:text-foreground",
      },
      size: {
        lg: "h-control w-full rounded-control px-6 text-lead",
        md: "h-13.5 rounded-control px-5 text-body",
        sm: "h-9 rounded-option px-3 text-caption",
        icon: "size-9 rounded-full [&_svg]:size-5",
      },
    },
    compoundVariants: [
      { variant: "ghost", className: "h-auto px-0" },
    ],
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
