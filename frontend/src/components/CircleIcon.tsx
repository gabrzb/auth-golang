import { cn } from "@/lib/utils"

type CircleIconProps = {
  children: React.ReactNode
  size?: number
  strokeWidth?: number
  fontFamily?: "display" | "hand"
  className?: string
}

export function CircleIcon({
  children,
  size = 48,
  strokeWidth = 2,
  fontFamily = "display",
  className,
}: CircleIconProps) {
  return (
    <div
      className={cn(
        "rounded-full border-ink flex items-center justify-center leading-none",
        fontFamily === "display" ? "font-display" : "font-hand",
        className
      )}
      style={{
        width: size,
        height: size,
        borderWidth: strokeWidth,
        fontSize: Math.round(size * 0.5),
      }}
    >
      {children}
    </div>
  )
}
