import Link from "next/link";
import { ChevronLeft } from "lucide-react";

import { Button } from "./button";

/** Круглая кнопка «назад» в шапке экранов онбординга. */
function BackLink({ href }: { href: string }) {
  return (
    <Button asChild variant="surface" size="icon">
      <Link href={href} aria-label="Назад">
        <ChevronLeft />
      </Link>
    </Button>
  );
}

export { BackLink };
