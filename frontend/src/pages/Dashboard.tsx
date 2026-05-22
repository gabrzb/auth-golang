import { useNavigate } from "react-router-dom"

import { AuthCard } from "@/components/AuthCard"
import { CircleIcon } from "@/components/CircleIcon"
import { Button } from "@/components/ui/button"

export function Dashboard() {
  const navigate = useNavigate()

  return (
    <AuthCard className="items-center text-center gap-5 py-14">
      <CircleIcon size={96} strokeWidth={2.5} fontFamily="display">
        ✓
      </CircleIcon>

      <h1 className="text-4xl leading-none">you're in!</h1>

      <p className="text-[15px] text-ink-mute font-hand max-w-[16rem]">
        signed in as <span className="text-ink">@placeholder</span>
      </p>

      <Button
        variant="accent"
        className="mt-2 w-full"
        disabled
        title="enabled in Phase 10 once the main app route exists"
      >
        take me home →
      </Button>

      <button
        type="button"
        onClick={() => {
          // Real sign-out call lands in Phase 10; for now just bounce back.
          navigate("/sign-in")
        }}
        className="link-dashed text-sm text-ink-mute font-hand cursor-pointer"
      >
        sign out
      </button>
    </AuthCard>
  )
}
