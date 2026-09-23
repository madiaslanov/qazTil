import * as React from "react";

/** Единый вид для «грузим» и «не вышло» внутри экрана. */
export function StateNote({
  text,
  action,
}: {
  text: string;
  action?: React.ReactNode;
}) {
  return (
    <div className="flex flex-1 flex-col items-center justify-center gap-4 px-6 text-center">
      <p className="text-[17px] text-muted">{text}</p>
      {action}
    </div>
  );
}
