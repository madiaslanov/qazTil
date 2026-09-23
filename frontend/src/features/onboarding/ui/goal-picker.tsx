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
      <legend className="mb-3 text-[18px] font-extrabold">
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
              "h-[72px] flex-1 rounded-field border-[3px] border-ink text-[15px] font-extrabold transition-[transform,box-shadow]",
              value === goal
                ? "bg-orange shadow-hard-sm"
                : "bg-card active:translate-x-[2px] active:translate-y-[2px]",
            )}
          >
            {goal} мин
          </button>
        ))}
      </div>
    </fieldset>
  );
}
