"use client";

import Image from "next/image";

import { GuestGuard } from "@/features/auth";
import { OnboardingForm } from "@/features/onboarding";
import { Screen } from "@/shared/ui";

/** Регистрация и вход: поля и выбор дневной цели. */
export function OnboardingPage() {
  return (
    <GuestGuard>
      <Screen>
        <Image
          src="/bg/onboarding.png"
          alt=""
          fill
          priority
          sizes="430px"
          className="pointer-events-none object-cover opacity-48 blur-[2px]"
        />

        <div className="relative flex flex-1 flex-col gap-6.5 px-6 pt-8.5 pb-7">
          <div className="flex items-center justify-between">
            <p className="text-title">QazTil</p>
            <span className="rounded-full border-2 border-primary bg-accent px-2.5 py-1.5 text-micro font-extrabold">
              BETA
            </span>
          </div>

          <div className="flex flex-col gap-2">
            <h1 className="text-display">
              Приобрети новую привычку.
              <br />
              Выучи казахский язык.
            </h1>
            <p className="text-body text-muted">
              Сохраняй ежедневный прогресс и увеличивай свои знания.
            </p>
          </div>

          <OnboardingForm />
        </div>
      </Screen>
    </GuestGuard>
  );
}
