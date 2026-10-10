import Image from "next/image";

import { cn } from "@/shared/lib/cn";

import type { WelcomeSlide } from "../model/slides";

/*
 * Композиция 342×300 из макета: свечение, тёмный шар и цветная «плитка».
 * Координаты — из Figma, поэтому позиционирование абсолютное.
 */
const LAYOUT = {
  mint: {
    glow: { src: "/welcome/glow-lime.svg", className: "left-19 -top-6.5" },
    control: "left-60.75 top-46.5",
    volume: "left-17 top-15.5 rounded-sheet bg-highlight",
    glyph: "text-glyph",
  },
  coral: {
    glow: { src: "/welcome/glow-lime.svg", className: "-left-7.25 top-2.75" },
    control: "left-60.75 top-46.5",
    volume: "left-24 top-13.25 -rotate-5 rounded-hero bg-accent",
    glyph: "text-glyph",
  },
  sky: {
    glow: { src: "/welcome/glow-coral.svg", className: "left-19 -top-6.5" },
    control: "left-6.5 top-46.5",
    volume: "left-22.5 top-15.5 rounded-sheet bg-sky",
    glyph: "text-display",
  },
} as const;

export function WelcomeIllustration({
  glyph,
  tone,
}: Pick<WelcomeSlide, "glyph" | "tone">) {
  const layout = LAYOUT[tone];

  return (
    <div aria-hidden className="relative mx-auto h-75 w-85.5 shrink-0">
      <Image
        src="/welcome/control.svg"
        alt=""
        width={96}
        height={96}
        className={cn("absolute", layout.control)}
      />
      <Image
        src={layout.glow.src}
        alt=""
        width={298}
        height={298}
        className={cn("absolute max-w-none", layout.glow.className)}
      />
      <div
        className={cn(
          "absolute flex h-38.5 w-50.5 items-center justify-center shadow-float",
          layout.volume,
        )}
      >
        <div className="flex h-21.5 w-26 items-center justify-center rounded-tile border-2 border-primary bg-surface">
          <span className={cn("font-bold text-primary", layout.glyph)}>{glyph}</span>
        </div>
      </div>
    </div>
  );
}
