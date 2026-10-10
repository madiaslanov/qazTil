"use client";

import { useRef, useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";

import { GuestGuard } from "@/features/auth";
import {
  BetaBadge,
  Button,
  Eyebrow,
  Logo,
  PageDots,
  Screen,
} from "@/shared/ui";

import { WELCOME_SLIDES } from "../model/slides";
import { WelcomeIllustration } from "./welcome-illustration";

/** После приветствия — настройка маршрута, потом регистрация. */
const NEXT_ROUTE = "/setup";
/** Минимальный горизонтальный свайп, чтобы перелистнуть слайд. */
const SWIPE_THRESHOLD = 50;

/** Три слайда приветствия для нового гостя. */
export function WelcomePage() {
  return (
    <GuestGuard>
      <WelcomeSlides />
    </GuestGuard>
  );
}

function WelcomeSlides() {
  const router = useRouter();
  const [index, setIndex] = useState(0);
  const touchStartX = useRef<number | null>(null);
  const slide = WELCOME_SLIDES[index];
  const isLast = index === WELCOME_SLIDES.length - 1;

  function goTo(next: number) {
    setIndex(Math.min(Math.max(next, 0), WELCOME_SLIDES.length - 1));
  }

  function advance() {
    if (isLast) router.push(NEXT_ROUTE);
    else goTo(index + 1);
  }

  return (
    <Screen
      className="px-gutter pt-6 pb-11"
      onTouchStart={(event) => {
        touchStartX.current = event.touches[0].clientX;
      }}
      onTouchEnd={(event) => {
        if (touchStartX.current === null) return;
        const delta = event.changedTouches[0].clientX - touchStartX.current;
        touchStartX.current = null;
        if (Math.abs(delta) < SWIPE_THRESHOLD) return;
        goTo(delta < 0 ? index + 1 : index - 1);
      }}
    >
      <header className="flex h-9 items-center justify-between">
        <Logo />
        {isLast ? (
          <BetaBadge />
        ) : (
          <Link
            href={NEXT_ROUTE}
            className="text-caption font-bold text-accent"
          >
            Пропустить
          </Link>
        )}
      </header>

      <main className="mt-10.5 flex flex-col gap-5.5">
        <div key={index} className="flex flex-col gap-2.5">
          <Eyebrow>Казахский без барьеров</Eyebrow>
          <h1 className="text-display text-heading">{slide.title}</h1>
          <p className="text-body text-muted">{slide.text}</p>
        </div>
        <WelcomeIllustration glyph={slide.glyph} tone={slide.tone} />
      </main>

      <footer className="mt-auto flex flex-col items-center gap-4.5 pt-6">
        <PageDots
          count={WELCOME_SLIDES.length}
          active={index}
          onSelect={goTo}
        />
        <Button onClick={advance}>{slide.action}</Button>
      </footer>
    </Screen>
  );
}
