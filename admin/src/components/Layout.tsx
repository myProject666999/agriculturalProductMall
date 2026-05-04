import React, { useState } from 'react'
import { Outlet, useNavigate, useLocation } from 'react-router-dom'
import { Layout as AntLayout, Menu, Avatar, Dropdown, Button, message } from 'antd'
import type { MenuProps } from 'antd'
import {
  DashboardOutlined,
  ProductOutlined,
  ShoppingOutlined,
  UserOutlined,
  ShopOutlined,
  StockOutlined,
  PictureOutlined,
  FileTextOutlined,
  NotificationOutlined,
  StarOutlined,
  TruckOutlined,
  ShoppingCartOutlined,
  SafetyOutlined,
  LogoutOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined,
} from '@ant-design/icons'
import { useAppSelector, useAppDispatch } from '@/store'
import { logout } from '@/store/userSlice'
import { authApi } from '@/services/api'

const { Header, Sider, Content } = AntLayout

const Layout: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const dispatch = useAppDispatch()
  const [collapsed, setCollapsed] = useState(false)

  const user = useAppSelector((state) => state.user.user)

  const getSelectedKeys = () => {
    const path = location.pathname
    if (path === '/dashboard') return ['dashboard']
    if (path.startsWith('/products')) return ['products']
    if (path.startsWith('/categories')) return ['categories']
    if (path.startsWith('/orders')) return ['orders']
    if (path.startsWith('/users')) return ['users']
    if (path.startsWith('/merchants')) return ['merchants']
    if (path.startsWith('/inventory')) return ['inventory']
    if (path.startsWith('/banners')) return ['banners']
    if (path.startsWith('/news')) return ['news']
    if (path.startsWith('/announcements')) return ['announcements']
    if (path.startsWith('/reviews')) return ['reviews']
    if (path.startsWith('/logistics')) return ['logistics']
    if (path.startsWith('/carts')) return ['carts']
    if (path.startsWith('/roles')) return ['roles']
    if (path.startsWith('/menus')) return ['menus']
    if (path.startsWith('/permissions')) return ['permissions']
    return ['dashboard']
  }

  const handleLogout = () => {
    authApi.logout().catch(() => {})
    dispatch(logout())
    message.success('已退出登录')
    navigate('/login')
  }

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'profile',
      icon: <UserOutlined />,
      label: '个人信息',
      onClick: () => navigate('/dashboard'),
    },
    {
      key: 'divider',
      type: 'divider',
    },
    {
      key: 'logout',
      icon: <LogoutOutlined />,
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  const menuItems: MenuProps['items'] = [
    {
      key: 'dashboard',
      icon: <DashboardOutlined />,
      label: '数据概览',
      onClick: () => navigate('/dashboard'),
    },
    {
      key: 'product',
      icon: <ProductOutlined />,
      label: '商品管理',
      children: [
        { key: 'products', label: '商品列表', onClick: () => navigate('/products') },
        { key: 'categories', label: '商品分类', onClick: () => navigate('/categories') },
      ],
    },
    {
      key: 'order',
      icon: <ShoppingOutlined />,
      label: '订单管理',
      children: [
        { key: 'orders', label: '订单列表', onClick: () => navigate('/orders') },
        { key: 'reviews', label: '评价管理', onClick: () => navigate('/reviews') },
      ],
    },
    {
      key: 'user',
      icon: <UserOutlined />,
      label: '用户管理',
      children: [
        { key: 'users', label: '用户列表', onClick: () => navigate('/users') },
        { key: 'carts', label: '购物车管理', onClick: () => navigate('/carts') },
      ],
    },
    {
      key: 'merchant',
      icon: <ShopOutlined />,
      label: '商户管理',
      children: [
        { key: 'merchants', label: '商户列表', onClick: () => navigate('/merchants') },
      ],
    },
    {
      key: 'inventory',
      icon: <StockOutlined />,
      label: '库存管理',
      children: [
        { key: 'inventory-list', label: '库存列表', onClick: () => navigate('/inventory') },
      ],
    },
    {
      key: 'content',
      icon: <PictureOutlined />,
      label: '内容管理',
      children: [
        { key: 'banners', label: '轮播图管理', onClick: () => navigate('/banners') },
        { key: 'news', label: '资讯管理', onClick: () => navigate('/news') },
        { key: 'announcements', label: '公告管理', onClick: () => navigate('/announcements') },
      ],
    },
    {
      key: 'logistics',
      icon: <TruckOutlined />,
      label: '物流管理',
      children: [
        { key: 'logistics-list', label: '物流列表', onClick: () => navigate('/logistics') },
      ],
    },
    {
      key: 'system',
      icon: <SafetyOutlined />,
      label: '系统管理',
      children: [
        { key: 'roles', label: '角色管理', onClick: () => navigate('/roles') },
        { key: 'menus', label: '菜单管理', onClick: () => navigate('/menus') },
        { key: 'permissions', label: '权限管理', onClick: () => navigate('/permissions') },
      ],
    },
  ]

  return (
    <AntLayout style={{ minHeight: '100vh' }}>
      <Sider trigger={null} collapsible collapsed={collapsed} theme="dark">
        <div className="sidebar-logo">{collapsed ? '🌿' : '🌿 农产品商城'}</div>
        <Menu
          theme="dark"
          mode="inline"
          selectedKeys={getSelectedKeys()}
          defaultOpenKeys={['product', 'order', 'user', 'content', 'system']}
          items={menuItems}
        />
      </Sider>
      <AntLayout>
        <Header
          style={{
            padding: '0 24px',
            background: '#fff',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
            style={{ fontSize: '16px', width: 64, height: 64 }}
          />
          <Dropdown menu={{ items: userMenuItems }} placement="bottomRight">
            <div className="header-user">
              <Avatar size={32} icon={<UserOutlined />} src={user?.avatar} />
              <span>{user?.username || '管理员'}</span>
            </div>
          </Dropdown>
        </Header>
        <Content
          style={{
            margin: '24px',
            padding: 24,
            background: '#fff',
            borderRadius: 8,
            minHeight: 280,
          }}
        >
          <Outlet />
        </Content>
      </AntLayout>
    </AntLayout>
  )
}

export default Layout
