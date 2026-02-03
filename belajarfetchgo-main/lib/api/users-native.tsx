function extractUsers(payload: any) {
  if (Array.isArray(payload)) return payload
  if (Array.isArray(payload?.data)) return payload.data
  if (Array.isArray(payload?.users)) return payload.users
  if (Array.isArray(payload?.result)) return payload.result
  if (Array.isArray(payload?.payload)) return payload.payload

  throw new Error(
    `Unexpected users response shape: ${typeof payload === "string" ? payload : JSON.stringify(payload)}`
  )
}

async function getUsers() {
  const res = await fetch("http://172.16.17.123:8080/users/barangs", {
    cache: "no-store",
  })

  if (!res.ok) {
    throw new Error(`Failed to fetch users: ${res.status} ${res.statusText}`)
  }

  const json = await res.json()
  return extractUsers(json)
}

export default getUsers