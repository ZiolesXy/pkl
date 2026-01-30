import api from "../axios"

export async function getUsers() {
  const res = await api.get("/users")
  const payload = res.data

  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.data)) return payload.data
  if (Array.isArray(payload?.users)) return payload.users
  if (Array.isArray(payload?.result)) return payload.result
  if (Array.isArray(payload?.payload)) return payload.payload

  throw new Error(
    `Unexpected users response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
}
