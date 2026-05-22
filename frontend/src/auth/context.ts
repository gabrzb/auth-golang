import { createContext } from "react"

import type { Credentials, User } from "@/api/auth"

export type AuthStatus = "loading" | "authed" | "anon"

export type AuthContextValue = {
  user: User | null
  status: AuthStatus
  signIn: (credentials: Credentials) => Promise<void>
  signUp: (credentials: Credentials) => Promise<void>
  signOut: () => Promise<void>
}

export const AuthContext = createContext<AuthContextValue | null>(null)
