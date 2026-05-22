import { cn } from "@/lib/utils"
import type { PasswordScore } from "@/lib/password"

type PasswordStrengthProps = {
  score: PasswordScore
  className?: string
}

const STRENGTH_LABELS: Record<PasswordScore, string> = {
  0: "empty",
  1: "weak",
  2: "fair",
  3: "good",
  4: "strong",
}

export function PasswordStrength({ score, className }: PasswordStrengthProps) {
  return (
    <div
      role="status"
      aria-live="polite"
      className={cn("flex gap-1", className)}
    >
      <span className="sr-only">password strength: {STRENGTH_LABELS[score]}</span>
      {[0, 1, 2, 3].map((i) => (
        <div
          key={i}
          aria-hidden="true"
          className={cn(
            "flex-1 h-[5px] rounded-sm",
            i < score
              ? "bg-ink"
              : "bg-paper-2 border border-dashed border-ink-mute"
          )}
        />
      ))}
    </div>
  )
}
