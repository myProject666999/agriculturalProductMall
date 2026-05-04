import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('admin_token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => {
    return response.data
  },
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('admin_token')
      localStorage.removeItem('admin_user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

export default api

export const authApi = {
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  logout: () => api.post('/auth/logout'),
  getInfo: () => api.get('/auth/info'),
}

export const statsApi = {
  getSalesStats: () => api.get('/admin/stats/sales'),
  getUserStats: () => api.get('/admin/stats/users'),
  getCategoryStats: () => api.get('/admin/stats/categories'),
  getProductStats: () => api.get('/admin/stats/products'),
}

export const productApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string; categoryId?: number }) =>
    api.get('/admin/products', { params }),
  getDetail: (id: number) => api.get(`/admin/products/${id}`),
  create: (data: any) => api.post('/admin/products', data),
  update: (id: number, data: any) => api.put(`/admin/products/${id}`, data),
  delete: (id: number) => api.delete(`/admin/products/${id}`),
  updateStatus: (id: number, status: number) => api.put(`/admin/products/${id}/status`, { status }),
}

export const categoryApi = {
  getList: () => api.get('/admin/categories'),
  getDetail: (id: number) => api.get(`/admin/categories/${id}`),
  create: (data: any) => api.post('/admin/categories', data),
  update: (id: number, data: any) => api.put(`/admin/categories/${id}`, data),
  delete: (id: number) => api.delete(`/admin/categories/${id}`),
}

export const orderApi = {
  getList: (params?: { page?: number; pageSize?: number; status?: number; keyword?: string }) =>
    api.get('/admin/orders', { params }),
  getDetail: (id: number) => api.get(`/admin/orders/${id}`),
  updateStatus: (id: number, status: number) => api.put(`/admin/orders/${id}/status`, { status }),
  ship: (id: number, data: { logisticsCompany: string; trackingNumber: string }) =>
    api.post(`/admin/orders/${id}/ship`, data),
}

export const userApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string; status?: number }) =>
    api.get('/admin/users', { params }),
  getDetail: (id: number) => api.get(`/admin/users/${id}`),
  updateStatus: (id: number, status: number) => api.put(`/admin/users/${id}/status`, { status }),
}

export const merchantApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string; status?: number }) =>
    api.get('/admin/merchants', { params }),
  getDetail: (id: number) => api.get(`/admin/merchants/${id}`),
  create: (data: any) => api.post('/admin/merchants', data),
  update: (id: number, data: any) => api.put(`/admin/merchants/${id}`, data),
  updateStatus: (id: number, status: number) =>
    api.put(`/admin/merchants/${id}/status`, { status }),
}

export const inventoryApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string; productId?: number }) =>
    api.get('/admin/inventory', { params }),
  getDetail: (id: number) => api.get(`/admin/inventory/${id}`),
  inStock: (data: { productId: number; quantity: number; remark?: string }) =>
    api.post('/admin/inventory/in', data),
  outStock: (data: { productId: number; quantity: number; remark?: string }) =>
    api.post('/admin/inventory/out', data),
}

export const bannerApi = {
  getList: (params?: { page?: number; pageSize?: number }) =>
    api.get('/admin/banners', { params }),
  getDetail: (id: number) => api.get(`/admin/banners/${id}`),
  create: (data: any) => api.post('/admin/banners', data),
  update: (id: number, data: any) => api.put(`/admin/banners/${id}`, data),
  delete: (id: number) => api.delete(`/admin/banners/${id}`),
  updateStatus: (id: number, status: number) =>
    api.put(`/admin/banners/${id}/status`, { status }),
}

export const newsApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string }) =>
    api.get('/admin/news', { params }),
  getDetail: (id: number) => api.get(`/admin/news/${id}`),
  create: (data: any) => api.post('/admin/news', data),
  update: (id: number, data: any) => api.put(`/admin/news/${id}`, data),
  delete: (id: number) => api.delete(`/admin/news/${id}`),
}

export const announcementApi = {
  getList: (params?: { page?: number; pageSize?: number }) =>
    api.get('/admin/announcements', { params }),
  getDetail: (id: number) => api.get(`/admin/announcements/${id}`),
  create: (data: any) => api.post('/admin/announcements', data),
  update: (id: number, data: any) => api.put(`/admin/announcements/${id}`, data),
  delete: (id: number) => api.delete(`/admin/announcements/${id}`),
}

export const reviewApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string }) =>
    api.get('/admin/reviews', { params }),
  getDetail: (id: number) => api.get(`/admin/reviews/${id}`),
  updateStatus: (id: number, status: number) =>
    api.put(`/admin/reviews/${id}/status`, { status }),
  delete: (id: number) => api.delete(`/admin/reviews/${id}`),
}

export const logisticsApi = {
  getList: (params?: { page?: number; pageSize?: number; keyword?: string }) =>
    api.get('/admin/logistics', { params }),
  getDetail: (id: number) => api.get(`/admin/logistics/${id}`),
  create: (data: any) => api.post('/admin/logistics', data),
  update: (id: number, data: any) => api.put(`/admin/logistics/${id}`, data),
  delete: (id: number) => api.delete(`/admin/logistics/${id}`),
}

export const cartApi = {
  getList: (params?: { page?: number; pageSize?: number; userId?: number }) =>
    api.get('/admin/carts', { params }),
  delete: (id: number) => api.delete(`/admin/carts/${id}`),
}

export const roleApi = {
  getList: () => api.get('/admin/roles'),
  getDetail: (id: number) => api.get(`/admin/roles/${id}`),
  create: (data: any) => api.post('/admin/roles', data),
  update: (id: number, data: any) => api.put(`/admin/roles/${id}`, data),
  delete: (id: number) => api.delete(`/admin/roles/${id}`),
}

export const menuApi = {
  getList: () => api.get('/admin/menus'),
  getDetail: (id: number) => api.get(`/admin/menus/${id}`),
  create: (data: any) => api.post('/admin/menus', data),
  update: (id: number, data: any) => api.put(`/admin/menus/${id}`, data),
  delete: (id: number) => api.delete(`/admin/menus/${id}`),
}

export const permissionApi = {
  getList: () => api.get('/admin/permissions'),
  getDetail: (id: number) => api.get(`/admin/permissions/${id}`),
  create: (data: any) => api.post('/admin/permissions', data),
  update: (id: number, data: any) => api.put(`/admin/permissions/${id}`, data),
  delete: (id: number) => api.delete(`/admin/permissions/${id}`),
}
