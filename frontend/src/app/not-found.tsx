import Link from "next/link";

import { Button, Screen } from "@/shared/ui";

export default function NotFound() {
  return (
    <Screen>
      <div className="flex flex-1 flex-col items-center justify-center gap-5 px-6 text-center">
        <p className="text-[38px] font-black leading-[1.02]">Страница потерялась</p>
        <p className="text-[15px] text-muted">
          Такого урока нет. Вернись на путь обучения.
        </p>
        <Button asChild variant="primary" className="w-auto px-8">
          <Link href="/learn">К обучению</Link>
        </Button>
      </div>
    </Screen>
  );
}
