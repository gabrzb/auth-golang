import type { ReactNode } from "react"
import { Navigate } from "react-router-dom"

import { useAuth } from "@/auth/useAuth"

export function PublicOnlyRoute({ children }: { children: ReactNode }) {
  const { status } = useAuth()

  if (status === "loading") {
    return null
  }

  if (status === "authed") {
    return <Navigate to="/" replace />
  }

  return children
}
