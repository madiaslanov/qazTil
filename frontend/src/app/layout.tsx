import type { Metadata, Viewport } from "next";
import { Geologica } from "next/font/google";

import { BRAND } from "@/shared/config/brand";

import { Providers } from "./providers";
import "@/shared/styles/globals.css";

// cyrillic-ext нужен для казахских букв: Қ, Ә, Ү, Ұ, Ғ, Ң, Ө, І.
const geologica = Geologica({
  variable: "--font-geologica",
  subsets: ["latin", "cyrillic", "cyrillic-ext"],
});

export const metadata: Metadata = {
  title: `${BRAND.name} — казахский язык`,
  description:
    "Учи казахский по урокам: словарь, категории, квиз и ежедневный прогресс.",
};

export const viewport: Viewport = {
  themeColor: BRAND.themeColor,
  viewportFit: "cover",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru" className={geologica.variable}>
      <body>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
