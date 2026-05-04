import { createSlice, PayloadAction } from '@reduxjs/toolkit'

interface CartItem {
  id: number
  productId: number
  productName: string
  price: number
  quantity: number
  image: string
  stock: number
}

interface CartState {
  items: CartItem[]
  total: number
}

const initialState: CartState = {
  items: [],
  total: 0,
}

const calculateTotal = (items: CartItem[]): number => {
  return items.reduce((sum, item) => sum + item.price * item.quantity, 0)
}

export const cartSlice = createSlice({
  name: 'cart',
  initialState,
  reducers: {
    addItem: (state, action: PayloadAction<CartItem>) => {
      const existingItem = state.items.find(
        (item) => item.productId === action.payload.productId
      )
      if (existingItem) {
        existingItem.quantity += action.payload.quantity
      } else {
        state.items.push(action.payload)
      }
      state.total = calculateTotal(state.items)
      localStorage.setItem('cart', JSON.stringify(state.items))
    },
    removeItem: (state, action: PayloadAction<number>) => {
      state.items = state.items.filter((item) => item.id !== action.payload)
      state.total = calculateTotal(state.items)
      localStorage.setItem('cart', JSON.stringify(state.items))
    },
    updateQuantity: (
      state,
      action: PayloadAction<{ id: number; quantity: number }>
    ) => {
      const item = state.items.find((item) => item.id === action.payload.id)
      if (item) {
        item.quantity = action.payload.quantity
        state.total = calculateTotal(state.items)
        localStorage.setItem('cart', JSON.stringify(state.items))
      }
    },
    clearCart: (state) => {
      state.items = []
      state.total = 0
      localStorage.removeItem('cart')
    },
    loadCart: (state, action: PayloadAction<CartItem[]>) => {
      state.items = action.payload
      state.total = calculateTotal(action.payload)
    },
  },
})

export const { addItem, removeItem, updateQuantity, clearCart, loadCart } =
  cartSlice.actions

export default cartSlice.reducer
