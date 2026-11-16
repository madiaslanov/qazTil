import type { DailyGoal } from "@/entities/learner";
import { OptionCard } from "@/shared/ui";

const goals: DailyGoal[] = [5, 10, 15];

/** Выбор дневной цели: три мятные плитки. */
export function GoalPicker({
  value,
  onChange,
}: {
  value: DailyGoal;
  onChange: (goal: DailyGoal) => void;
}) {
  return (
    <fieldset className="flex flex-col gap-3">
      <legend className="mb-1 text-stat font-extrabold text-heading">
        Установи ежедневную цель
      </legend>
      <p className="mb-2 text-body-sm text-muted">
        Опыт за день. Каждый верный ответ — 5 XP.
      </p>
      <div className="flex gap-2.5">
        {goals.map((goal) => (
          <OptionCard
            key={goal}
            tone="soft"
            selected={value === goal}
            onClick={() => onChange(goal)}
            className="h-18 flex-1 text-body font-extrabold text-heading"
          >
            {goal} XP
          </OptionCard>
        ))}
      </div>
    </fieldset>
  );
}
