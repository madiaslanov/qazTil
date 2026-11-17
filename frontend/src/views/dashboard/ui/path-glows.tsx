"use client";

import { useEffect, useRef, useState } from "react";
import Image from "next/image";

/*
 * Размытые пятна за путём. Координаты и размеры — из Figma
 * (область под шапкой, 390×678): это один «блок» фона. Путь длиннее экрана,
 * поэтому блок повторяется по высоте, а каждый второй зеркалится,
 * чтобы пятна не шли одним и тем же рисунком.
 */
const BAND_HEIGHT = 678;

const glows = [
  { src: "/path/glow-4.svg", width: 371, height: 372, className: "-left-15.5 -top-10.5" },
  { src: "/path/glow-3.svg", width: 386, height: 389, className: "left-8 top-23.25" },
  { src: "/path/glow-1.svg", width: 411, height: 414, className: "-left-27.25 top-48.5" },
  { src: "/path/glow-2.svg", width: 427, height: 431, className: "left-15.75 top-81.25" },
] as const;

export function PathGlows() {
  const root = useRef<HTMLDivElement>(null);
  const [bands, setBands] = useState(1);

  // Блоков столько, сколько нужно, чтобы закрыть всю высоту страницы пути.
  useEffect(() => {
    const element = root.current;
    if (!element) return;
    const update = () =>
      setBands(Math.max(1, Math.ceil(element.clientHeight / BAND_HEIGHT)));
    update();
    const observer = new ResizeObserver(update);
    observer.observe(element);
    return () => observer.disconnect();
  }, []);

  return (
    <div
      ref={root}
      aria-hidden
      className="pointer-events-none absolute inset-0 overflow-hidden"
    >
      {Array.from({ length: bands }, (_, band) => (
        <div
          key={band}
          className={`absolute inset-x-0 ${band % 2 === 1 ? "-scale-x-100" : ""}`}
          style={{ top: band * BAND_HEIGHT, height: BAND_HEIGHT }}
        >
          {glows.map((glow) => (
            <Image
              key={glow.src}
              src={glow.src}
              alt=""
              width={glow.width}
              height={glow.height}
              className={`absolute max-w-none ${glow.className}`}
            />
          ))}
        </div>
      ))}
    </div>
  );
}
