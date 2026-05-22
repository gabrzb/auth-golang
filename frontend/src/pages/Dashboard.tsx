import { useNavigate } from "react-router-dom"
import { toast } from "sonner"

import { useAuth } from "@/auth/useAuth"
import { AuthCard } from "@/components/AuthCard"
import { CircleIcon } from "@/components/CircleIcon"
import { Button } from "@/components/ui/button"

export function Dashboard() {
  const navigate = useNavigate()
  const { signOut, user } = useAuth()

  const createdAt = user
    ? new Intl.DateTimeFormat("en", {
        dateStyle: "medium",
        timeStyle: "short",
      }).format(new Date(user.created_at))
    : null

  async function handleSignOut() {
    try {
      await signOut()
    } catch {
      toast.error("Signed out locally, but the server did not respond")
    } finally {
      navigate("/sign-in", { replace: true })
    }
  }

  return (
    <AuthCard className="items-center text-center gap-5 py-14">
      <CircleIcon size={96} strokeWidth={2.5} fontFamily="display">
        ✓
      </CircleIcon>

      <h1 className="text-4xl leading-none">you're in!</h1>

      <p className="text-[15px] text-ink-mute font-hand max-w-[16rem]">
        signed in as <span className="text-ink">{user?.email}</span>
      </p>

      {createdAt && (
        <p className="text-xs text-ink-mute font-mono uppercase tracking-wider">
          created {createdAt}
        </p>
      )}

      <Button variant="accent" className="mt-2 w-full" onClick={() => navigate("/")}>
        take me home →
      </Button>

      <button
        type="button"
        onClick={handleSignOut}
        className="link-dashed text-sm text-ink-mute font-hand cursor-pointer"
      >
        sign out
      </button>
    </AuthCard>
  )
}
