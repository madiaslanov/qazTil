import { BRAND } from "@/shared/config/brand";
import { cn } from "@/shared/lib/cn";

/** Словесный знак QazTil в шапках. */
function Logo({ className }: { className?: string }) {
  return (
    <p className={cn("text-title font-semibold text-foreground", className)}>
      {BRAND.name}
    </p>
  );
}

export { Logo };
