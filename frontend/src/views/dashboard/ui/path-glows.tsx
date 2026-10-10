import Image from "next/image";

/*
 * Размытые пятна за путём. Координаты и размеры — из Figma
 * (область под шапкой, 390×678), поэтому позиционирование абсолютное.
 */
const glows = [
  { src: "/path/glow-4.svg", width: 371, height: 372, className: "-left-15.5 -top-10.5" },
  { src: "/path/glow-3.svg", width: 386, height: 389, className: "left-8 top-23.25" },
  { src: "/path/glow-1.svg", width: 411, height: 414, className: "-left-27.25 top-48.5" },
  { src: "/path/glow-2.svg", width: 427, height: 431, className: "left-15.75 top-81.25" },
] as const;

export function PathGlows() {
  return (
    <div aria-hidden className="pointer-events-none absolute inset-0 overflow-hidden">
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
  );
}
