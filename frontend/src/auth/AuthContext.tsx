import {
  useCallback,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from "react"

import * as authApi from "@/api/auth"
import { signOutEventName } from "@/api/client"
import {
  AuthContext,
  type AuthStatus,
} from "@/auth/context"
import { clearAccessToken, setAccessToken } from "@/auth/storage"

let bootPromise: Promise<authApi.User> | null = null

function bootstrapSession() {
  if (!bootPromise) {
    bootPromise = authApi
      .refresh()
      .then(async (tokenPayload) => {
        setAccessToken(tokenPayload.access_token)
        return authApi.me()
      })
      .finally(() => {
        bootPromise = null
      })
  }

  return bootPromise
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<authApi.User | null>(null)
  const [status, setStatus] = useState<AuthStatus>("loading")

  const becomeAnon = useCallback(() => {
    clearAccessToken()
    setUser(null)
    setStatus("anon")
  }, [])

  const hydrateUser = useCallback(async () => {
    const userPayload = await authApi.me()
    setUser(userPayload)
    setStatus("authed")
  }, [])

  useEffect(() => {
    let active = true

    async function boot() {
      try {
        const userPayload = await bootstrapSession()
        if (!active) return

        setUser(userPayload)
        setStatus("authed")
      } catch {
        if (active) {
          becomeAnon()
        }
      }
    }

    boot()

    return () => {
      active = false
    }
  }, [becomeAnon])

  useEffect(() => {
    window.addEventListener(signOutEventName, becomeAnon)
    return () => window.removeEventListener(signOutEventName, becomeAnon)
  }, [becomeAnon])

  const signIn = useCallback(
    async (credentials: authApi.Credentials) => {
      const tokenPayload = await authApi.login(credentials)
      setAccessToken(tokenPayload.access_token)
      await hydrateUser()
    },
    [hydrateUser]
  )

  const signUp = useCallback(
    async (credentials: authApi.Credentials) => {
      await authApi.register(credentials)
      await signIn(credentials)
    },
    [signIn]
  )

  const signOut = useCallback(async () => {
    try {
      await authApi.logout()
    } finally {
      becomeAnon()
    }
  }, [becomeAnon])

  const value = useMemo(
    () => ({ user, status, signIn, signUp, signOut }),
    [signIn, signOut, signUp, status, user]
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
