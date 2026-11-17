import { notFound } from "next/navigation";

import { LessonPage } from "@/views/lesson";

export default async function Page({ params }: PageProps<"/lesson/[lessonId]">) {
  const { lessonId } = await params;
  const id = Number(lessonId);
  if (!Number.isInteger(id) || id <= 0) {
    notFound();
  }
  return <LessonPage lessonId={id} />;
}
