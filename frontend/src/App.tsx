import React from 'react'
import { Routes, Route } from 'react-router-dom'
import { Layout, ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import Header from './components/layout/Header'
import Footer from './components/layout/Footer'
import Home from './pages/Home'
import ProductList from './pages/product/ProductList'
import ProductDetail from './pages/product/ProductDetail'
import Cart from './pages/cart/Cart'
import Order from './pages/order/Order'
import OrderDetail from './pages/order/OrderDetail'
import UserCenter from './pages/user/UserCenter'
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import News from './pages/News'
import PrivateRoute from './components/common/PrivateRoute'

const { Content } = Layout

const App: React.FC = () => {
  return (
    <ConfigProvider locale={zhCN}>
      <Layout className="min-h-screen">
        <Header />
        <Content className="bg-gray-50">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/register" element={<Register />} />
            <Route path="/products" element={<ProductList />} />
            <Route path="/products/:id" element={<ProductDetail />} />
            <Route path="/news" element={<News />} />
            <Route path="/cart" element={<PrivateRoute><Cart /></PrivateRoute>} />
            <Route path="/orders" element={<PrivateRoute><Order /></PrivateRoute>} />
            <Route path="/orders/:id" element={<PrivateRoute><OrderDetail /></PrivateRoute>} />
            <Route path="/user/*" element={<PrivateRoute><UserCenter /></PrivateRoute>} />
          </Routes>
        </Content>
        <Footer />
      </Layout>
    </ConfigProvider>
  )
}

export default App
