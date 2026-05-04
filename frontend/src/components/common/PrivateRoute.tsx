import React from 'react'
import { Navigate } from 'react-router-dom'
import { useAppSelector } from '@/store'

interface PrivateRouteProps {
  children: React.ReactNode
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const isLoggedIn = useAppSelector((state) => state.user.isLoggedIn)
  const token = localStorage.getItem('token')

  if (!isLoggedIn && !token) {
    return <Navigate to="/login" replace />
  }

  return <>{children}</>
}

export default PrivateRoute
