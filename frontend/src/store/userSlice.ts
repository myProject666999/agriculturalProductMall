import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface User {
  id: number
  username: string
  email: string
  avatar?: string
  role: string
}

interface UserState {
  user: User | null
  token: string | null
  isLoggedIn: boolean
}

const initialState: UserState = {
  user: null,
  token: localStorage.getItem('token'),
  isLoggedIn: !!localStorage.getItem('token'),
}

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    login: (state, action: PayloadAction<{ user: User; token: string }>) => {
      state.user = action.payload.user
      state.token = action.payload.token
      state.isLoggedIn = true
      localStorage.setItem('token', action.payload.token)
      localStorage.setItem('user', JSON.stringify(action.payload.user))
    },
    logout: (state) => {
      state.user = null
      state.token = null
      state.isLoggedIn = false
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    },
    updateUser: (state, action: PayloadAction<User>) => {
      state.user = action.payload.user
      localStorage.setItem('user', JSON.stringify(action.payload))
    },
  },
})

export const { login, logout, updateUser } = userSlice.actions

export default userSlice.reducer
