"use client";

import { progressApi } from "@/entities/progress";
import { useRequest } from "@/shared/lib/use-request";
import { Button, Card, Progress, Screen, StateNote } from "@/shared/ui";
import { BottomNav } from "@/widgets/bottom-nav";
import { TopBar } from "@/widgets/top-bar";

/** Прогресс по категориям: то, что считает сам сервер. */
export function ProgressPage() {
  const progress = useRequest(() => progressApi.list(), "progress");

  return (
    <Screen>
      <TopBar />

      <div className="flex min-h-0 flex-1 flex-col gap-4 overflow-y-auto px-6 pt-6 pb-8">
        <div>
          <p className="text-[12px] font-extrabold uppercase text-muted">
            Статистика
          </p>
          <h1 className="mt-1 text-[26px]">Твой прогресс</h1>
        </div>

        {progress.loading && <StateNote text="Считаем ответы…" />}
        {progress.error && (
          <StateNote
            text={progress.error}
            action={
              <Button
                size="md"
                variant="quiet"
                className="w-auto"
                onClick={progress.retry}
              >
                Повторить
              </Button>
            }
          />
        )}

        {progress.data?.length === 0 && (
          <StateNote text="Пока ни одного ответа — начни первый урок." />
        )}

        {progress.data?.map((item) => {
          const percent =
            item.total_answers === 0
              ? 0
              : Math.round((item.correct_answers / item.total_answers) * 100);
          return (
            <Card key={item.category_id} className="flex flex-col gap-3 p-5">
              <div className="flex items-baseline justify-between gap-3">
                <div>
                  <p className="text-[19px]">{item.category_name_ru}</p>
                  <p className="text-[13px] text-subtle">
                    {item.category_name_kk}
                  </p>
                </div>
                <b className="text-[19px] font-extrabold">{percent}%</b>
              </div>
              <Progress value={percent} />
              <p className="text-[13px] text-muted">
                Верно {item.correct_answers} из {item.total_answers}
              </p>
            </Card>
          );
        })}
      </div>

      <BottomNav />
    </Screen>
  );
}
