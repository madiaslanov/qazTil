import type { Metadata, Viewport } from "next";
import { Inter } from "next/font/google";

import "@/shared/config/globals.css";

const inter = Inter({
  variable: "--font-inter",
  subsets: ["latin", "cyrillic"],
  weight: ["400", "700", "800", "900"],
});

export const metadata: Metadata = {
  title: "QazTil — казахский язык",
  description:
    "Учи казахский по урокам: словарь, категории, квиз и ежедневный прогресс.",
};

export const viewport: Viewport = {
  themeColor: "#f2f0e9",
  viewportFit: "cover",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru" className={inter.variable}>
      <body>{children}</body>
    </html>
  );
}
