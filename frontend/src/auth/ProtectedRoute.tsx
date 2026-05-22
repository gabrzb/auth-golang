import type { ReactNode } from "react"
import { Navigate } from "react-router-dom"

import { useAuth } from "@/auth/useAuth"

export function ProtectedRoute({ children }: { children: ReactNode }) {
  const { status } = useAuth()

  if (status === "loading") {
    return null
  }

  if (status === "anon") {
    return <Navigate to="/sign-in" replace />
  }

  return children
}
