import { notFound } from "next/navigation";

import { LessonPage } from "@/views/lesson";

export default async function Page({ params }: PageProps<"/lesson/[categoryId]">) {
  const { categoryId } = await params;
  const id = Number(categoryId);
  if (!Number.isInteger(id) || id <= 0) {
    notFound();
  }
  return <LessonPage categoryId={id} />;
}
