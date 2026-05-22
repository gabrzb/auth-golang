import { cn } from "@/lib/utils"

type PasswordStrengthProps = {
  score: 0 | 1 | 2 | 3 | 4
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

export function passwordScore(value: string): 0 | 1 | 2 | 3 | 4 {
  if (!value) return 0
  let score = 0
  if (value.length >= 8) score++
  if (value.length >= 12) score++
  if (/[A-Z]/.test(value) && /[a-z]/.test(value)) score++
  if (/\d/.test(value) && /[^A-Za-z0-9]/.test(value)) score++
  return Math.min(score, 4) as 0 | 1 | 2 | 3 | 4
}
