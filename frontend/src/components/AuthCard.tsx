import { cn } from "@/lib/utils"

type AuthCardProps = {
  children: React.ReactNode
  className?: string
}

export function AuthCard({ children, className }: AuthCardProps) {
  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-12 bg-paper-2">
      <div
        className={cn(
          "w-full max-w-[34rem] min-h-[40rem] bg-paper rounded-lg px-14 py-20 flex flex-col justify-center gap-6 shadow-ink",
          className
        )}
      >
        {children}
      </div>
    </div>
  )
}
