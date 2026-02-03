import api from "../axios";

export const putData = async <T>(
  endpoint: string,
  id: number | string,
  payload: T
) => {
  const res = await api.put(`${endpoint}/${id}`, payload);
  return res.data;
};