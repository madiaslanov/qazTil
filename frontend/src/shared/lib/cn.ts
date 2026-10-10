import { clsx, type ClassValue } from "clsx";
import { extendTailwindMerge } from "tailwind-merge";

/**
 * Кастомные шкалы из shared/styles/tokens.css. Без них tailwind-merge
 * принимает text-body за цвет и выкидывает его рядом с text-muted.
 */
const twMerge = extendTailwindMerge({
  extend: {
    theme: {
      text: [
        "hero",
        "display",
        "h1",
        "h2",
        "h3",
        "title",
        "stat",
        "lead",
        "body",
        "body-sm",
        "caption",
        "eyebrow",
        "micro",
        "nano",
        "glyph",
      ],
      radius: ["option", "chip", "control", "card", "tile", "sheet", "hero"],
      shadow: [
        "button",
        "glow",
        "selected",
        "float",
        "sheet",
        "dialog",
        "success",
        "danger",
        "module",
      ],
      spacing: ["gutter", "bar", "nav", "control", "field", "answer"],
      container: ["phone"],
    },
  },
});

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}
