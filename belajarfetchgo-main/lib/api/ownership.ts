import api from "../axios"

export interface Ownership {
  barang_id: number
  user_id: number
}

export async function getOwnership(): Promise<Ownership[]> {
  const res = await api.get("/user/barang")
  const payload = res.data

  if (Array.isArray(payload)) return payload as Ownership[]
  if (Array.isArray(payload?.data)) return payload.data as Ownership[]
  if (Array.isArray(payload?.ownership)) return payload.ownership as Ownership[]
  if (Array.isArray(payload?.result)) return payload.result as Ownership[]
  if (Array.isArray(payload?.payload)) return payload.payload as Ownership[]

  throw new Error(
    `Unexpected ownership response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
  
}
