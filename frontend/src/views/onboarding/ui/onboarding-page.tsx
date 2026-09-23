"use client";

import { useEffect } from "react";
import Image from "next/image";
import { useRouter } from "next/navigation";

import { useLearner } from "@/entities/learner";
import { OnboardingForm } from "@/features/onboarding";
import { Screen } from "@/shared/ui";

/** Первый экран: приветствие, поля входа и выбор дневной цели. */
export function OnboardingPage() {
  const router = useRouter();
  const learner = useLearner();

  useEffect(() => {
    if (learner) router.replace("/learn");
  }, [learner, router]);

  return (
    <Screen>
      <Image
        src="/bg/onboarding.png"
        alt=""
        fill
        priority
        sizes="430px"
        className="pointer-events-none object-cover opacity-48 blur-[2px]"
      />

      <div className="relative flex flex-1 flex-col gap-6.5 px-6 pt-[34px] pb-7">
        <div className="flex items-center justify-between">
          <p className="text-[24px]">QazTil</p>
          <span className="rounded-full border-2 border-ink bg-orange px-2.5 py-1.5 text-[11px] font-extrabold">
            BETA
          </span>
        </div>

        <div className="flex flex-col gap-2">
          <h1 className="text-[36px] leading-[1.02]">
            Приобрети новую привычку.
            <br />
            Выучи казахский язык.
          </h1>
          <p className="text-[15px] leading-[1.45] text-muted">
            Сохраняй ежедневный прогресс и увеличивай свои знания.
          </p>
        </div>

        <OnboardingForm />
      </div>
    </Screen>
  );
}
