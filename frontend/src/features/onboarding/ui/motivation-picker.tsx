import { Briefcase, Coffee, GraduationCap, type LucideIcon } from "lucide-react";

import type { Motivation } from "@/entities/learner";
import { OptionCard } from "@/shared/ui";

const motivations: { value: Motivation; title: string; icon: LucideIcon }[] = [
  { value: "work", title: "Работа", icon: Briefcase },
  { value: "study", title: "Учеба", icon: GraduationCap },
  { value: "life", title: "Быт", icon: Coffee },
];

/** Зачем ученику язык: три плитки с иконками. */
export function MotivationPicker({
  value,
  onChange,
}: {
  value: Motivation;
  onChange: (motivation: Motivation) => void;
}) {
  return (
    <fieldset className="flex flex-col gap-2.5">
      <legend className="mb-2.5 text-caption font-extrabold text-primary">
        Для чего тебе язык?
      </legend>
      <div className="flex gap-2">
        {motivations.map(({ value: option, title, icon: Icon }) => (
          <OptionCard
            key={option}
            selected={option === value}
            onClick={() => onChange(option)}
            className="flex h-20.5 flex-1 flex-col items-center justify-center gap-1.75"
          >
            <Icon className="size-5" />
            <span className="text-eyebrow font-bold">{title}</span>
          </OptionCard>
        ))}
      </div>
    </fieldset>
  );
}
