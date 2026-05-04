import React, { useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Layout, Input, Button, Badge, Dropdown, Avatar, Space, MenuProps, message } from 'antd'
import { ShoppingCartOutlined, UserOutlined, SearchOutlined } from '@ant-design/icons'
import { useAppSelector, useAppDispatch } from '@/store'
import { logout } from '@/store/userSlice'
import { authApi, cartApi } from '@/services/api'

const { Header: AntHeader } = Layout

const Header: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [searchParams] = useSearchParams()
  const [searchValue, setSearchValue] = React.useState(searchParams.get('keyword') || '')
  const [cartCount, setCartCount] = React.useState(0)

  const user = useAppSelector((state) => state.user.user)
  const isLoggedIn = useAppSelector((state) => state.user.isLoggedIn)
  const cartItems = useAppSelector((state) => state.cart.items)

  useEffect(() => {
    const storedUser = localStorage.getItem('user')
    if (storedUser && !user) {
      // User is stored in localStorage but not in Redux, we could dispatch a login action here
    }
  }, [user])

  useEffect(() => {
    if (isLoggedIn) {
      cartApi.getList().then((res: any) => {
        if (res.code === 0) {
          setCartCount(res.data?.length || 0)
        }
      }).catch(() => {})
    } else {
      const localCart = localStorage.getItem('cart')
      if (localCart) {
        setCartCount(JSON.parse(localCart).length)
      } else {
        setCartCount(cartItems.length)
      }
    }
  }, [isLoggedIn, cartItems.length])

  const handleSearch = () => {
    if (searchValue.trim()) {
      navigate(`/products?keyword=${encodeURIComponent(searchValue.trim())}`)
    }
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter') {
      handleSearch()
    }
  }

  const handleLogout = () => {
    authApi.logout().then(() => {
      dispatch(logout())
      message.success('已退出登录')
      navigate('/')
    }).catch(() => {
      dispatch(logout())
      navigate('/')
    })
  }

  const userMenuItems: MenuProps['items'] = [
    {
      key: 'center',
      label: '个人中心',
      onClick: () => navigate('/user'),
    },
    {
      key: 'orders',
      label: '我的订单',
      onClick: () => navigate('/orders'),
    },
    {
      key: 'favorites',
      label: '我的收藏',
      onClick: () => navigate('/user?tab=favorites'),
    },
    {
      key: 'divider',
      type: 'divider',
    },
    {
      key: 'logout',
      label: '退出登录',
      onClick: handleLogout,
    },
  ]

  const guestMenuItems: MenuProps['items'] = [
    {
      key: 'login',
      label: '登录',
      onClick: () => navigate('/login'),
    },
    {
      key: 'register',
      label: '注册',
      onClick: () => navigate('/register'),
    },
  ]

  return (
    <AntHeader className="header-container">
      <div className="flex items-center justify-between h-full max-w-7xl mx-auto">
        <div className="flex items-center gap-8">
          <div className="logo" onClick={() => navigate('/')}>
            🌿 农产品商城
          </div>
          <Space className="hidden md:flex">
            <Button type="link" className="menu-link" onClick={() => navigate('/')}>
              首页
            </Button>
            <Button type="link" className="menu-link" onClick={() => navigate('/products')}>
              商品中心
            </Button>
            <Button type="link" className="menu-link" onClick={() => navigate('/news')}>
              资讯中心
            </Button>
          </Space>
        </div>

        <div className="flex items-center gap-4">
          <Input
            placeholder="搜索商品..."
            prefix={<SearchOutlined />}
            value={searchValue}
            onChange={(e) => setSearchValue(e.target.value)}
            onKeyPress={handleKeyPress}
            style={{ width: 250 }}
          />
          <Button type="primary" className="btn-primary" onClick={handleSearch}>
            搜索
          </Button>

          <Badge count={cartCount} className="cart-badge">
            <Button
              type="text"
              icon={<ShoppingCartOutlined style={{ fontSize: 20, color: 'white' }} />}
              onClick={() => navigate('/cart')}
            />
          </Badge>

          <Dropdown
            menu={{ items: isLoggedIn ? userMenuItems : guestMenuItems }}
            placement="bottomRight"
          >
            <Space className="cursor-pointer">
              <Avatar
                size={32}
                icon={<UserOutlined />}
                src={user?.avatar}
                className="user-avatar"
              />
              <span className="text-white">{isLoggedIn ? user?.username : '登录/注册'}</span>
            </Space>
          </Dropdown>
        </div>
      </div>
    </AntHeader>
  )
}

export default Header
