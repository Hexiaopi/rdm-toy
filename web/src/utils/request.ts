import axios from 'axios'
import { ElMessage } from 'element-plus'

const service = axios.create({
    baseURL: import.meta.env.VITE_APP_BASE_API,
    timeout: 5000,
    headers: { 'Content-Type': 'application/json;charset=utf-8' },
})

service.interceptors.request.use(
    (config) => {
        return config
    },
    (error) => {
        return Promise.reject(error);
    }
)

service.interceptors.response.use(
    (response) => {
        if (response.status === 200) {
            if (response.data.code !== "000000") {
                ElMessage.error(response.data.desc)
                return Promise.reject(new Error(response.data.desc || 'Error'))
            } else {
                return response.data
            }
        } else {
            Promise.reject()
        }
    },
    (error) => {
        return Promise.reject(error)
    }
)

export default service;