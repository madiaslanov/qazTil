import type { Metadata } from "next";
import "./globals.css";
// import { Nav } from "@/components/nav";

export const metadata: Metadata = {
  title: "qazTil — казахский язык",
  description: "Словарь, категории, квиз и прогресс для изучения казахского языка.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="ru">
      <body>
        <header className="site-header">
          <p className="brand-mark">qazTil</p>
          <h1>Қазақ тілі</h1>
          <p className="lede">
            Словарь, категории и квиз. Прогресс хранится на сервере для одного локального ученика.
          </p>
        </header>
        {/* <Nav /> */}
        <main>{children}</main>
      </body>
    </html>
  );
}
