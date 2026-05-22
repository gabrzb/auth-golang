import * as React from "react"

import { cn } from "@/lib/utils"

function Input({ className, type, ...props }: React.ComponentProps<"input">) {
  return (
    <input
      type={type}
      data-slot="input"
      className={cn(
        "w-full rounded-md border-[1.5px] border-dashed border-ink bg-paper px-3 py-2.5 font-hand text-base text-ink-soft placeholder:text-ink-mute placeholder:italic outline-none transition-colors focus:border-solid focus:border-ink-soft disabled:opacity-50 disabled:cursor-not-allowed aria-invalid:border-destructive aria-invalid:text-destructive",
        className
      )}
      {...props}
    />
  )
}

export { Input }
