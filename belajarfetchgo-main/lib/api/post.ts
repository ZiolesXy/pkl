
import api from "../axios"

export async function postJson<TResponse = any, TBody = any>(
  path: string,
  body: TBody
): Promise<TResponse> {
  const res = await api.post(path, body)
  return res.data as TResponse
}

