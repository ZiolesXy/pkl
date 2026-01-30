// lib/axios.ts
import axios from "axios"

const api = axios.create({
  baseURL: "http://172.16.17.67:8080",
  timeout: 5000,
})

export default api
