import { AppShell } from "@/widgets/app-shell";

/** Вкладки с общей шапкой и нижней навигацией. */
export default function Layout({ children }: LayoutProps<"/">) {
  return <AppShell>{children}</AppShell>;
}
