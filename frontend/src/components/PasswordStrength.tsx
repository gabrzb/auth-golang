import { cn } from "@/lib/utils"
import type { PasswordScore } from "@/lib/password"

type PasswordStrengthProps = {
  score: PasswordScore
  className?: string
}

export function PasswordStrength({ score, className }: PasswordStrengthProps) {
  return (
    <div className={cn("flex gap-1", className)} aria-hidden="true">
      {[0, 1, 2, 3].map((i) => (
        <div
          key={i}
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
