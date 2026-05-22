import * as React from "react"
import { cva, type VariantProps } from "class-variance-authority"
import { Slot } from "radix-ui"

import { cn } from "@/lib/utils"

const buttonVariants = cva(
  "inline-flex items-center justify-center rounded-lg border-2 font-hand font-bold lowercase select-none transition-[box-shadow,transform] shadow-ink hover:-translate-y-px active:translate-x-[2px] active:translate-y-[2px] active:[box-shadow:1px_1px_0_var(--color-ink)] disabled:opacity-50 disabled:pointer-events-none focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-ink-soft cursor-pointer",
  {
    variants: {
      variant: {
        default: "bg-ink text-paper border-ink",
        accent: "bg-brand text-white border-brand",
        ghost: "bg-transparent text-ink border-ink",
      },
      size: {
        default: "px-4 py-3 text-lg",
        sm: "px-3 py-2 text-base",
        lg: "px-5 py-3.5 text-xl",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
    },
  }
)

function Button({
  className,
  variant = "default",
  size = "default",
  asChild = false,
  ...props
}: React.ComponentProps<"button"> &
  VariantProps<typeof buttonVariants> & {
    asChild?: boolean
  }) {
  const Comp = asChild ? Slot.Root : "button"

  return (
    <Comp
      data-slot="button"
      data-variant={variant}
      data-size={size}
      className={cn(buttonVariants({ variant, size, className }))}
      {...props}
    />
  )
}

export { Button, buttonVariants }
