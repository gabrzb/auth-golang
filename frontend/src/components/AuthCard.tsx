import { cn } from "@/lib/utils"

type AuthCardProps = {
  children: React.ReactNode
  className?: string
}

export function AuthCard({ children, className }: AuthCardProps) {
  return (
    <div className="min-h-screen flex items-center justify-center px-4 py-10 bg-paper-2">
      <div
        className={cn(
          "w-full max-w-[26rem] bg-paper rounded-md px-9 py-12 flex flex-col gap-4",
          className
        )}
      >
        {children}
      </div>
    </div>
  )
}
