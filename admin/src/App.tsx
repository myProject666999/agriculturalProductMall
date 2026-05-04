import React from 'react'
import { Routes, Route, Navigate } from 'react-router-dom'
import { ConfigProvider } from 'antd'
import zhCN from 'antd/locale/zh_CN'
import Login from './pages/Login'
import Layout from './components/Layout'
import Dashboard from './pages/Dashboard'
import PlaceholderPage from './components/common/PlaceholderPage'
import PrivateRoute from './components/common/PrivateRoute'

const App: React.FC = () => {
  return (
    <ConfigProvider
      locale={zhCN}
      theme={{
        token: {
          colorPrimary: '#52c41a',
        },
      }}
    >
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <PrivateRoute>
              <Layout />
            </PrivateRoute>
          }
        >
          <Route index element={<Navigate to="/dashboard" replace />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="products" element={<PlaceholderPage title="商品列表" />} />
          <Route path="categories" element={<PlaceholderPage title="商品分类" />} />
          <Route path="orders" element={<PlaceholderPage title="订单列表" />} />
          <Route path="orders/:id" element={<PlaceholderPage title="订单详情" />} />
          <Route path="users" element={<PlaceholderPage title="用户列表" />} />
          <Route path="merchants" element={<PlaceholderPage title="商户列表" />} />
          <Route path="inventory" element={<PlaceholderPage title="库存列表" />} />
          <Route path="banners" element={<PlaceholderPage title="轮播图管理" />} />
          <Route path="news" element={<PlaceholderPage title="资讯管理" />} />
          <Route path="announcements" element={<PlaceholderPage title="公告管理" />} />
          <Route path="reviews" element={<PlaceholderPage title="评价管理" />} />
          <Route path="logistics" element={<PlaceholderPage title="物流列表" />} />
          <Route path="carts" element={<PlaceholderPage title="购物车管理" />} />
          <Route path="roles" element={<PlaceholderPage title="角色管理" />} />
          <Route path="menus" element={<PlaceholderPage title="菜单管理" />} />
          <Route path="permissions" element={<PlaceholderPage title="权限管理" />} />
        </Route>
      </Routes>
    </ConfigProvider>
  )
}

export default App
