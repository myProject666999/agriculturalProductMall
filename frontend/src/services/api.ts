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
    const token = localStorage.getItem('token')
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
      localStorage.removeItem('token')
      localStorage.removeItem('user')
      window.location.href = '/login'
    }
    return Promise.reject(error.response?.data || error)
  }
)

export default api

export const authApi = {
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  register: (data: {
    username: string
    email: string
    password: string
    confirmPassword: string
    code: string
  }) => api.post('/auth/register', data),
  sendCode: (email: string) => api.post('/auth/send-code', { email }),
  getInfo: () => api.get('/auth/info'),
  logout: () => api.post('/auth/logout'),
}

export const productApi = {
  getList: (params?: {
    page?: number
    pageSize?: number
    categoryId?: number
    keyword?: string
  }) => api.get('/products', { params }),
  getDetail: (id: number) => api.get(`/products/${id}`),
  getCategories: () => api.get('/categories'),
  getRecommend: () => api.get('/products/recommend'),
  search: (keyword: string) => api.get('/products/search', { params: { keyword } }),
}

export const cartApi = {
  getList: () => api.get('/cart'),
  add: (data: { productId: number; quantity: number }) =>
    api.post('/cart/add', data),
  update: (data: { id: number; quantity: number }) =>
    api.put('/cart/update', data),
  remove: (id: number) => api.delete(`/cart/remove/${id}`),
  clear: () => api.delete('/cart/clear'),
}

export const orderApi = {
  getList: (status?: number) => api.get('/orders', { params: { status } }),
  getDetail: (id: number) => api.get(`/orders/${id}`),
  create: (data: {
    cartIds: number[]
    addressId: number
    remark?: string
  }) => api.post('/orders/create', data),
  pay: (id: number) => api.post(`/orders/pay/${id}`),
  cancel: (id: number, reason: string) =>
    api.post(`/orders/cancel/${id}`, { reason }),
  confirm: (id: number) => api.post(`/orders/confirm/${id}`),
  refund: (id: number, data: { reason: string; amount: number }) =>
    api.post(`/orders/refund/${id}`, data),
  evaluate: (id: number, data: { rating: number; content: string; images?: string[] }) =>
    api.post(`/orders/evaluate/${id}`, data),
}

export const userApi = {
  getAddresses: () => api.get('/user/addresses'),
  addAddress: (data: any) => api.post('/user/addresses', data),
  updateAddress: (id: number, data: any) =>
    api.put(`/user/addresses/${id}`, data),
  deleteAddress: (id: number) => api.delete(`/user/addresses/${id}`),
  getFavorites: () => api.get('/user/favorites'),
  addFavorite: (productId: number) =>
    api.post('/user/favorites', { productId }),
  removeFavorite: (productId: number) =>
    api.delete(`/user/favorites/${productId}`),
  updateInfo: (data: any) => api.put('/user/info', data),
  updatePassword: (data: { oldPassword: string; newPassword: string }) =>
    api.put('/user/password', data),
}

export const bannerApi = {
  getList: () => api.get('/banners'),
}

export const newsApi = {
  getList: (params?: { page?: number; pageSize?: number }) =>
    api.get('/news', { params }),
  getDetail: (id: number) => api.get(`/news/${id}`),
}
