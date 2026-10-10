"use client";

import * as React from "react";
import { Dialog } from "radix-ui";

import { cn } from "@/shared/lib/cn";

/**
 * Шторка снизу экрана: превью урока, разбор ответа.
 * Ширина — как у телефонного холста, чтобы на десктопе не растягивалась.
 */
function BottomSheet({
  open,
  onOpenChange,
  title,
  description,
  className,
  children,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Для скринридеров: заголовок шторки, если он не виден в контенте. */
  title: string;
  description?: string;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-scrim animate-fade-in" />
        <Dialog.Content
          // Без описания Radix ругается в консоли — явно говорим, что его нет.
          {...(!description && { "aria-describedby": undefined })}
          className={cn(
            "fixed inset-x-0 bottom-0 z-50 mx-auto flex w-full max-w-phone flex-col gap-5 rounded-t-sheet bg-surface px-gutter pt-2.5 pb-7 shadow-sheet outline-none animate-sheet-in",
            className,
          )}
        >
          <Dialog.Title className="sr-only">{title}</Dialog.Title>
          {description && (
            <Dialog.Description className="sr-only">
              {description}
            </Dialog.Description>
          )}
          <div aria-hidden className="flex justify-center">
            <span className="h-1.25 w-11 rounded-full bg-line" />
          </div>
          {children}
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

export { BottomSheet };
