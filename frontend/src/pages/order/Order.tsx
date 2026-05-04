import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Table, Card, Button, Tag, Space, Empty, Spin, message, Popconfirm, Tabs, TabsProps } from 'antd'
import { EyeOutlined, PayCircleOutlined, CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons'
import { orderApi } from '@/services/api'
import dayjs from 'dayjs'

const Order: React.FC = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [orders, setOrders] = useState<any[]>([])
  const [activeTab, setActiveTab] = useState<string>('all')

  useEffect(() => {
    loadOrders()
  }, [activeTab])

  const loadOrders = async () => {
    setLoading(true)
    try {
      const params: any = {}
      if (activeTab !== 'all') {
        params.status = Number(activeTab)
      }
      const res: any = await orderApi.getList(activeTab !== 'all' ? Number(activeTab) : undefined)
      if (res.code === 0) {
        setOrders(res.data?.list || res.data || [])
      }
    } catch (error) {
      console.error('Failed to load orders:', error)
    } finally {
      setLoading(false)
    }
  }

  const getStatusTag = (status: number) => {
    const statusMap: Record<number, { text: string; color: string }> = {
      0: { text: '待付款', color: 'orange' },
      1: { text: '待发货', color: 'blue' },
      2: { text: '待收货', color: 'cyan' },
      3: { text: '已完成', color: 'green' },
      4: { text: '已取消', color: 'default' },
      5: { text: '退款中', color: 'purple' },
      6: { text: '已退款', color: 'default' },
    }
    const info = statusMap[status] || { text: '未知', color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  const handlePay = async (id: number) => {
    try {
      const res: any = await orderApi.pay(id)
      if (res.code === 0) {
        message.success('支付成功')
        loadOrders()
      } else {
        message.error(res.message || '支付失败')
      }
    } catch (error) {
      message.error('支付失败，请稍后重试')
    }
  }

  const handleCancel = async (id: number) => {
    try {
      const res: any = await orderApi.cancel(id, '用户取消')
      if (res.code === 0) {
        message.success('订单已取消')
        loadOrders()
      } else {
        message.error(res.message || '取消失败')
      }
    } catch (error) {
      message.error('取消失败，请稍后重试')
    }
  }

  const handleConfirm = async (id: number) => {
    try {
      const res: any = await orderApi.confirm(id)
      if (res.code === 0) {
        message.success('确认收货成功')
        loadOrders()
      } else {
        message.error(res.message || '操作失败')
      }
    } catch (error) {
      message.error('操作失败，请稍后重试')
    }
  }

  const tabItems: TabsProps['items'] = [
    { key: 'all', label: '全部订单' },
    { key: '0', label: '待付款' },
    { key: '1', label: '待发货' },
    { key: '2', label: '待收货' },
    { key: '3', label: '已完成' },
  ]

  const columns = [
    {
      title: '订单信息',
      key: 'orderInfo',
      render: (_: any, record: any) => (
        <div>
          <p className="text-gray-500 text-sm">订单号: {record.orderNo}</p>
          <p className="text-gray-400 text-xs">
            {dayjs(record.createdAt).format('YYYY-MM-DD HH:mm:ss')}
          </p>
        </div>
      ),
    },
    {
      title: '商品信息',
      key: 'products',
      render: (_: any, record: any) => (
        <div>
          {(record.items || []).slice(0, 2).map((item: any, idx: number) => (
            <div key={idx} className="flex items-center gap-2 mb-1">
              <img
                src={item.image || 'https://picsum.photos/40/40'}
                alt={item.productName}
                className="w-10 h-10 object-cover rounded"
              />
              <div>
                <p className="text-sm truncate max-w-xs">{item.productName}</p>
                <p className="text-gray-400 text-xs">x{item.quantity}</p>
              </div>
            </div>
          ))}
          {(record.items?.length > 2) && (
            <p className="text-gray-400 text-xs">还有 {record.items.length - 2} 件商品</p>
          )}
        </div>
      ),
    },
    {
      title: '订单金额',
      key: 'amount',
      render: (_: any, record: any) => (
        <div className="text-right">
          <p className="text-red-500 font-bold">¥{record.totalAmount?.toFixed(2) || '0.00'}</p>
          {record.shippingFee > 0 && (
            <p className="text-gray-400 text-xs">含运费 ¥{record.shippingFee}</p>
          )}
        </div>
      ),
    },
    {
      title: '订单状态',
      key: 'status',
      render: (_: any, record: any) => getStatusTag(record.status),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" size="small" onClick={() => navigate(`/orders/${record.id}`)}>
            <EyeOutlined /> 详情
          </Button>
          {record.status === 0 && (
            <>
              <Button type="primary" size="small" icon={<PayCircleOutlined />} onClick={() => handlePay(record.id)}>
                付款
              </Button>
              <Popconfirm
                title="确定要取消订单吗？"
                onConfirm={() => handleCancel(record.id)}
                okText="确定"
                cancelText="取消"
              >
                <Button type="text" danger size="small">
                  取消
                </Button>
              </Popconfirm>
            </>
          )}
          {record.status === 2 && (
            <Button type="primary" size="small" icon={<CheckCircleOutlined />} onClick={() => handleConfirm(record.id)}>
              确认收货
            </Button>
          )}
        </Space>
      ),
    },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">我的订单</h1>

      <Card>
        <Tabs
          activeKey={activeTab}
          onChange={setActiveTab}
          items={tabItems}
        />

        <Spin spinning={loading}>
          {orders.length > 0 ? (
            <Table
              columns={columns}
              dataSource={orders}
              rowKey="id"
              pagination={{
                pageSize: 10,
                showTotal: (total) => `共 ${total} 条记录`,
              }}
            />
          ) : (
            <div className="py-12">
              <Empty description="暂无订单" />
            </div>
          )}
        </Spin>
      </Card>
    </div>
  )
}

export default Order
