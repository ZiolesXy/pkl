import axios from "axios"

const api = axios.create({
  baseURL: "http://172.16.17.123:8080", 
  headers: {
    "Content-Type": "application/json",
  },
})

export default api