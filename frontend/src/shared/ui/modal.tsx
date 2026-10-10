"use client";

import * as React from "react";
import { Dialog } from "radix-ui";

import { cn } from "@/shared/lib/cn";

/** Карточка-диалог по центру экрана: «Жизни закончились» и подобное. */
function Modal({
  open,
  onOpenChange,
  title,
  className,
  children,
}: {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  /** Для скринридеров: заголовок, если он не виден в контенте. */
  title: string;
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <Dialog.Root open={open} onOpenChange={onOpenChange}>
      <Dialog.Portal>
        <Dialog.Overlay className="fixed inset-0 z-40 bg-background/80 backdrop-blur-sm animate-fade-in" />
        <div className="pointer-events-none fixed inset-0 z-50 flex items-center justify-center px-gutter">
          <Dialog.Content
            aria-describedby={undefined}
            className={cn(
              "pointer-events-auto flex w-full max-w-85.5 flex-col items-center gap-5.5 rounded-sheet bg-background px-6 py-7 shadow-dialog outline-none animate-pop-in",
              className,
            )}
          >
            <Dialog.Title className="sr-only">{title}</Dialog.Title>
            {children}
          </Dialog.Content>
        </div>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

export { Modal };
