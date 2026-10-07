import { AuthGuard } from "@/features/auth";

/** Экраны обучения: без сессии отсюда уводит на онбординг. */
export default function Layout({ children }: LayoutProps<"/">) {
  return <AuthGuard>{children}</AuthGuard>;
}
