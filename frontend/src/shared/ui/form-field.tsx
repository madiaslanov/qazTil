import * as React from "react";

import { cn } from "@/shared/lib/cn";

import { FieldHint } from "./field-hint";
import { Label } from "./label";

/** Подпись, поле и подсказка под ним — одной колонкой, как в макете. */
function FormField({
  id,
  label,
  hint,
  hintTone,
  className,
  children,
}: {
  id: string;
  label: string;
  hint?: React.ReactNode;
  hintTone?: React.ComponentProps<typeof FieldHint>["tone"];
  className?: string;
  children: React.ReactNode;
}) {
  return (
    <div className={cn("flex flex-col gap-1.75", className)}>
      <Label htmlFor={id}>{label}</Label>
      {children}
      {hint && (
        <FieldHint id={`${id}-hint`} tone={hintTone}>
          {hint}
        </FieldHint>
      )}
    </div>
  );
}

export { FormField };
