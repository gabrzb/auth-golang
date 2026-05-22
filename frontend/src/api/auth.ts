import { request } from "@/api/client"

export type User = {
  id: number
  email: string
  created_at: string
  updated_at: string
}

export type AccessTokenResponse = {
  access_token: string
  expires_in: number
}

export type Credentials = {
  email: string
  password: string
}

function jsonInit(method: string, body?: unknown): RequestInit {
  return {
    method,
    headers: {
      "Content-Type": "application/json",
    },
    body: body ? JSON.stringify(body) : undefined,
  }
}

export function register(credentials: Credentials) {
  return request<User>("/auth/register", jsonInit("POST", credentials), {
    skipAuthRefresh: true,
  })
}

export function login(credentials: Credentials) {
  return request<AccessTokenResponse>(
    "/auth/login",
    jsonInit("POST", credentials),
    { skipAuthRefresh: true }
  )
}

export function refresh() {
  return request<AccessTokenResponse>(
    "/auth/refresh",
    { method: "POST" },
    { skipAuthRefresh: true }
  )
}

export function logout() {
  return request<{ message: string }>("/auth/logout", { method: "POST" })
}

export function me() {
  return request<User>("/me")
}
