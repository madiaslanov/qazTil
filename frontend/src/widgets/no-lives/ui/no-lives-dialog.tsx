"use client";

import { Clock3, HeartCrack } from "lucide-react";

import { useLives } from "@/entities/learner";
import { Button, Modal } from "@/shared/ui";

function formatCountdown(ms: number) {
  const totalSeconds = Math.ceil(ms / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;
  return `${minutes}:${String(seconds).padStart(2, "0")}`;
}

/** «Жизни закончились»: таймер до новой жизни и выход на путь. */
export function NoLivesDialog({
  open,
  onClose,
}: {
  open: boolean;
  onClose: () => void;
}) {
  const { msToNextLife } = useLives();

  return (
    <Modal
      open={open}
      onOpenChange={(next) => !next && onClose()}
      title="Жизни закончились"
    >
      <span className="flex size-29 items-center justify-center rounded-full bg-danger-soft">
        <HeartCrack className="size-15.5 text-danger" strokeWidth={1.5} />
      </span>

      <div className="flex flex-col gap-2 text-center">
        <p className="text-h1 font-bold text-heading">Жизни закончились</p>
        <p className="text-body-sm text-muted">
          Сделай паузу или восстанови запас, чтобы продолжить урок сейчас.
        </p>
      </div>

      {msToNextLife !== null && (
        <div className="flex w-full items-center justify-between rounded-control bg-surface px-4 py-3.5">
          <span className="flex items-center gap-2 text-caption font-bold text-foreground">
            <Clock3 className="size-5 text-primary" />
            Новая жизнь через
          </span>
          <span className="text-stat font-extrabold text-danger tabular-nums">
            {formatCountdown(msToNextLife)}
          </span>
        </div>
      )}

      {/* Оплаты в приложении пока нет — кнопка из макета ждёт её. */}
      <Button disabled title="Покупки скоро появятся">
        Купить 5 жизней · 120 ₸
      </Button>

      <Button variant="ghost" size="sm" onClick={onClose}>
        Вернуться на путь
      </Button>
    </Modal>
  );
}
