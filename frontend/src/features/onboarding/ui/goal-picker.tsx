import type { DailyGoal } from "@/entities/learner";
import { cn } from "@/shared/lib/cn";

const goals: DailyGoal[] = [5, 10, 15];

/** Выбор дневной цели: три карточки, активная — оранжевая. */
export function GoalPicker({
  value,
  onChange,
}: {
  value: DailyGoal;
  onChange: (goal: DailyGoal) => void;
}) {
  return (
    <fieldset className="flex flex-col gap-3">
      <legend className="mb-3 text-stat font-extrabold">
        Установи ежедневную цель
      </legend>
      <div className="flex gap-2.5">
        {goals.map((goal) => (
          <button
            key={goal}
            type="button"
            onClick={() => onChange(goal)}
            aria-pressed={value === goal}
            className={cn(
              "h-18 flex-1 rounded-option border-3 border-primary text-body font-extrabold transition-[transform,box-shadow]",
              value === goal
                ? "bg-accent shadow-button"
                : "bg-surface active:translate-x-0.5 active:translate-y-0.5",
            )}
          >
            {goal} мин
          </button>
        ))}
      </div>
    </fieldset>
  );
}
