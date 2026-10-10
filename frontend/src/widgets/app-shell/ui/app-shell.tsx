import { Screen } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { TopBar } from "@/widgets/top-bar";

/** Каркас вкладок: шапка и навигация на месте, прокручивается только контент. */
export function AppShell({ children }: { children: React.ReactNode }) {
  return (
    <Screen className="h-dvh">
      <TopBar />
      <main className="relative flex min-h-0 flex-1 flex-col overflow-y-auto">
        {children}
      </main>
      <BottomNav />
    </Screen>
  );
}
