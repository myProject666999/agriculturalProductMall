import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Button, Checkbox, InputNumber, Empty, message, Card, Row, Col, Tag, Popconfirm, Space } from 'antd'
import { DeleteOutlined, ShoppingCartOutlined } from '@ant-design/icons'
import { useAppSelector, useAppDispatch } from '@/store'
import { removeItem, updateQuantity, clearCart, addItem } from '@/store/cartSlice'
import { cartApi } from '@/services/api'

const Cart: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [loading, setLoading] = useState(false)
  const [selectedRowKeys, setSelectedRowKeys] = useState<React.Key[]>([])
  const [cartItems, setCartItems] = useState<any[]>([])
  
  const reduxCartItems = useAppSelector((state) => state.cart.items)
  const isLoggedIn = useAppSelector((state) => state.user.isLoggedIn)

  useEffect(() => {
    loadCart()
  }, [isLoggedIn])

  const loadCart = async () => {
    if (isLoggedIn) {
      setLoading(true)
      try {
        const res: any = await cartApi.getList()
        if (res.code === 0) {
          setCartItems(res.data || [])
        }
      } catch (error) {
        console.error('Failed to load cart:', error)
        setCartItems(reduxCartItems)
      } finally {
        setLoading(false)
      }
    } else {
      setCartItems(reduxCartItems)
    }
  }

  const handleQuantityChange = async (id: number, quantity: number) => {
    dispatch(updateQuantity({ id, quantity }))
    setCartItems((prev) =>
      prev.map((item) => (item.id === id ? { ...item, quantity } : item))
    )
  }

  const handleRemove = async (id: number) => {
    if (isLoggedIn) {
      try {
        await cartApi.remove(id)
        message.success('已移除')
      } catch (error) {
        console.error('Failed to remove item:', error)
      }
    }
    dispatch(removeItem(id))
    setCartItems((prev) => prev.filter((item) => item.id !== id))
    setSelectedRowKeys((prev) => prev.filter((key) => key !== id))
  }

  const handleClearCart = async () => {
    if (isLoggedIn) {
      try {
        await cartApi.clear()
        message.success('购物车已清空')
      } catch (error) {
        console.error('Failed to clear cart:', error)
      }
    }
    dispatch(clearCart())
    setCartItems([])
    setSelectedRowKeys([])
  }

  const handleCheckout = () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要结算的商品')
      return
    }
    if (!isLoggedIn) {
      message.warning('请先登录')
      navigate('/login')
      return
    }
    navigate('/order/checkout')
  }

  const handleSelectAll = (e: any) => {
    if (e.target.checked) {
      setSelectedRowKeys(cartItems.map((item) => item.id))
    } else {
      setSelectedRowKeys([])
    }
  }

  const handleRowSelect = (record: any, checked: boolean) => {
    if (checked) {
      setSelectedRowKeys((prev) => [...prev, record.id])
    } else {
      setSelectedRowKeys((prev) => prev.filter((key) => key !== record.id))
    }
  }

  const selectedItems = cartItems.filter((item) => selectedRowKeys.includes(item.id))
  const totalPrice = selectedItems.reduce((sum, item) => sum + item.price * item.quantity, 0)
  const totalCount = selectedItems.reduce((sum, item) => sum + item.quantity, 0)

  const columns = [
    {
      title: '选择',
      key: 'select',
      width: 50,
      render: (_: any, record: any) => (
        <Checkbox
          checked={selectedRowKeys.includes(record.id)}
          onChange={(e) => handleRowSelect(record, e.target.checked)}
        />
      ),
    },
    {
      title: '商品信息',
      key: 'product',
      render: (_: any, record: any) => (
        <div className="flex items-center gap-4">
          <img
            src={record.image || record.productImage || 'https://picsum.photos/80/80'}
            alt={record.productName || record.name}
            className="w-20 h-20 object-cover rounded"
            onClick={() => navigate(`/products/${record.productId}`)}
            style={{ cursor: 'pointer' }}
          />
          <div>
            <p className="font-medium cursor-pointer hover:text-green-600" onClick={() => navigate(`/products/${record.productId}`)}>
              {record.productName || record.name}
            </p>
            <p className="text-gray-500 text-sm">单价: ¥{record.price}</p>
          </div>
        </div>
      ),
    },
    {
      title: '单价',
      dataIndex: 'price',
      key: 'price',
      render: (price: number) => <span className="text-red-500 font-bold">¥{price}</span>,
    },
    {
      title: '数量',
      key: 'quantity',
      render: (_: any, record: any) => (
        <InputNumber
          min={1}
          max={record.stock || 999}
          value={record.quantity}
          onChange={(val) => handleQuantityChange(record.id, val || 1)}
        />
      ),
    },
    {
      title: '小计',
      key: 'subtotal',
      render: (_: any, record: any) => (
        <span className="text-red-500 font-bold">¥{(record.price * record.quantity).toFixed(2)}</span>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Popconfirm
          title="确定要移除该商品吗？"
          onConfirm={() => handleRemove(record.id)}
          okText="确定"
          cancelText="取消"
        >
          <Button type="text" danger icon={<DeleteOutlined />}>
            删除
          </Button>
        </Popconfirm>
      ),
    },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6 flex items-center gap-2">
        <ShoppingCartOutlined /> 我的购物车
      </h1>

      {cartItems.length > 0 ? (
        <>
          <Card>
            <Table
              columns={columns}
              dataSource={cartItems}
              rowKey="id"
              loading={loading}
              pagination={false}
            />
          </Card>

          <Card className="mt-6 sticky bottom-0">
            <Row justify="space-between" align="middle">
              <Col>
                <Space>
                  <Checkbox
                    checked={selectedRowKeys.length === cartItems.length && cartItems.length > 0}
                    indeterminate={selectedRowKeys.length > 0 && selectedRowKeys.length < cartItems.length}
                    onChange={handleSelectAll}
                  >
                    全选
                  </Checkbox>
                  <Popconfirm
                    title="确定要清空购物车吗？"
                    onConfirm={handleClearCart}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button type="text" danger>
                      清空购物车
                    </Button>
                  </Popconfirm>
                </Space>
              </Col>
              <Col>
                <Space size="large">
                  <span>
                    已选 <span className="text-green-600 font-bold">{selectedRowKeys.length}</span> 件商品
                  </span>
                  <span>
                    共 <span className="text-green-600 font-bold">{totalCount}</span> 件
                  </span>
                  <span className="text-lg">
                    合计: <span className="text-red-500 text-2xl font-bold">¥{totalPrice.toFixed(2)}</span>
                  </span>
                  <Button
                    type="primary"
                    size="large"
                    onClick={handleCheckout}
                    disabled={selectedRowKeys.length === 0}
                    style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
                  >
                    去结算
                  </Button>
                </Space>
              </Col>
            </Row>
          </Card>
        </>
      ) : (
        <Card>
          <Empty
            description="购物车是空的"
            image={Empty.PRESENTED_IMAGE_SIMPLE}
          >
            <Button type="primary" onClick={() => navigate('/products')} style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}>
              去逛逛
            </Button>
          </Empty>
        </Card>
      )}
    </div>
  )
}

export default Cart
