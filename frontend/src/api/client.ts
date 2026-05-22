import {
  clearAccessToken,
  getAccessToken,
  setAccessToken,
} from "@/auth/storage"

export const signOutEventName = "auth:signOut"

type RequestOptions = {
  skipAuthRefresh?: boolean
}

type TokenResponse = {
  access_token: string
  expires_in: number
}

export class ApiError extends Error {
  status: number
  payload: unknown

  constructor(status: number, message: string, payload: unknown) {
    super(message)
    this.name = "ApiError"
    this.status = status
    this.payload = payload
  }
}

let refreshPromise: Promise<string> | null = null

function apiUrl(path: string) {
  return `/api${path.startsWith("/") ? path : `/${path}`}`
}

function dispatchSignOut() {
  window.dispatchEvent(new Event(signOutEventName))
}

async function parseResponse(response: Response) {
  const contentType = response.headers.get("content-type") ?? ""
  if (contentType.includes("application/json")) {
    return response.json()
  }

  const text = await response.text()
  return text.length > 0 ? text : null
}

function errorMessage(payload: unknown, fallback: string) {
  if (
    payload &&
    typeof payload === "object" &&
    "error" in payload &&
    typeof payload.error === "string"
  ) {
    return payload.error
  }

  return fallback
}

async function fetchJson<T>(path: string, init: RequestInit = {}) {
  const headers = new Headers(init.headers)
  const token = getAccessToken()

  if (token) {
    headers.set("Authorization", `Bearer ${token}`)
  }

  const response = await fetch(apiUrl(path), {
    ...init,
    credentials: "include",
    headers,
  })
  const payload = await parseResponse(response)

  if (!response.ok) {
    throw new ApiError(
      response.status,
      errorMessage(payload, response.statusText),
      payload
    )
  }

  return payload as T
}

async function refreshAccessToken() {
  if (!refreshPromise) {
    refreshPromise = fetch(apiUrl("/auth/refresh"), {
      method: "POST",
      credentials: "include",
    })
      .then(async (response) => {
        const payload = await parseResponse(response)

        if (!response.ok) {
          throw new ApiError(
            response.status,
            errorMessage(payload, response.statusText),
            payload
          )
        }

        const tokenPayload = payload as TokenResponse
        setAccessToken(tokenPayload.access_token)
        return tokenPayload.access_token
      })
      .finally(() => {
        refreshPromise = null
      })
  }

  return refreshPromise
}

export async function request<T>(
  path: string,
  init: RequestInit = {},
  options: RequestOptions = {}
) {
  try {
    return await fetchJson<T>(path, init)
  } catch (error) {
    if (
      options.skipAuthRefresh ||
      !(error instanceof ApiError) ||
      error.status !== 401
    ) {
      throw error
    }

    try {
      await refreshAccessToken()
      return await fetchJson<T>(path, init)
    } catch (refreshError) {
      clearAccessToken()
      dispatchSignOut()
      throw refreshError
    }
  }
}
