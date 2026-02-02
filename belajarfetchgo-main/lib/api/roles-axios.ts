import api from "../axios"

export async function getRoles() {
  const res = await api.get("/roles")
  const payload = res.data

  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.data)) return payload.data
  if (Array.isArray(payload?.roles)) return payload.roles
  if (Array.isArray(payload?.result)) return payload.result
  if (Array.isArray(payload?.payload)) return payload.payload

  throw new Error(
    `Unexpected roles response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
  
}
