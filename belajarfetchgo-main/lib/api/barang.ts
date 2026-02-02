import api from "../axios"

export async function getBarangs() {
  const res = await api.get("/barangs")
  const payload = res.data

  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.data)) return payload.data
  if (Array.isArray(payload?.barangs)) return payload.barangs
  if (Array.isArray(payload?.result)) return payload.result
  if (Array.isArray(payload?.payload)) return payload.payload

  throw new Error(
    `Unexpected barangs response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
  
}
