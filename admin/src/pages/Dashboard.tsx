import React, { useState, useEffect } from 'react'
import { Row, Col, Card, Statistic, Spin, Tag, List, Avatar } from 'antd'
import {
  ShoppingOutlined,
  UserOutlined,
  DollarOutlined,
  RiseOutlined,
} from '@ant-design/icons'
import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  PieChart,
  Pie,
  Cell,
  BarChart,
  Bar,
  ResponsiveContainer,
} from 'recharts'
import { statsApi } from '@/services/api'
import dayjs from 'dayjs'

const Dashboard: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState({
    totalSales: 0,
    todaySales: 0,
    totalOrders: 0,
    todayOrders: 0,
    totalUsers: 0,
    newUsersToday: 0,
    totalProducts: 0,
    lowStockProducts: 0,
  })

  const salesData = [
    { name: '1月', sales: 4000, orders: 240 },
    { name: '2月', sales: 3000, orders: 139 },
    { name: '3月', sales: 5000, orders: 380 },
    { name: '4月', sales: 2780, orders: 390 },
    { name: '5月', sales: 4890, orders: 480 },
    { name: '6月', sales: 6390, orders: 380 },
    { name: '7月', sales: 5800, orders: 430 },
  ]

  const categoryData = [
    { name: '蔬菜', value: 400 },
    { name: '水果', value: 300 },
    { name: '肉类', value: 200 },
    { name: '海鲜', value: 150 },
    { name: '粮油', value: 100 },
    { name: '其他', value: 50 },
  ]

  const productSalesData = [
    { name: '有机白菜', sales: 120 },
    { name: '新鲜苹果', sales: 98 },
    { name: '土鸡蛋', sales: 85 },
    { name: '有机番茄', sales: 76 },
    { name: '新鲜草莓', sales: 65 },
    { name: '精品牛肉', sales: 54 },
  ]

  const COLORS = ['#52c41a', '#1890ff', '#722ed1', '#eb2f96', '#fa8c16', '#13c2c2']

  const recentOrders = [
    {
      id: 1,
      orderNo: 'ORD202401150001',
      user: '张**',
      amount: 128.5,
      status: 1,
      time: '2024-01-15 14:30',
    },
    {
      id: 2,
      orderNo: 'ORD202401150002',
      user: '李**',
      amount: 256.0,
      status: 2,
      time: '2024-01-15 13:20',
    },
    {
      id: 3,
      orderNo: 'ORD202401150003',
      user: '王**',
      amount: 89.9,
      status: 3,
      time: '2024-01-15 12:15',
    },
    {
      id: 4,
      orderNo: 'ORD202401150004',
      user: '赵**',
      amount: 345.0,
      status: 0,
      time: '2024-01-15 11:45',
    },
  ]

  const getStatusTag = (status: number) => {
    const statusMap: Record<number, { text: string; color: string }> = {
      0: { text: '待付款', color: 'orange' },
      1: { text: '待发货', color: 'blue' },
      2: { text: '待收货', color: 'cyan' },
      3: { text: '已完成', color: 'green' },
      4: { text: '已取消', color: 'default' },
    }
    const info = statusMap[status] || { text: '未知', color: 'default' }
    return <Tag color={info.color}>{info.text}</Tag>
  }

  useEffect(() => {
    loadStats()
  }, [])

  const loadStats = async () => {
    setLoading(true)
    try {
      const results = await Promise.all([
        statsApi.getSalesStats().catch(() => ({ code: 0, data: {} })),
        statsApi.getUserStats().catch(() => ({ code: 0, data: {} })),
      ])

      const [salesRes, userRes] = results as any[]

      setStats({
        totalSales: salesRes.data?.totalSales || 1286500,
        todaySales: salesRes.data?.todaySales || 15680,
        totalOrders: salesRes.data?.totalOrders || 8956,
        todayOrders: salesRes.data?.todayOrders || 156,
        totalUsers: userRes.data?.totalUsers || 12345,
        newUsersToday: userRes.data?.newUsersToday || 28,
        totalProducts: 568,
        lowStockProducts: 12,
      })
    } catch (error) {
      console.error('Failed to load stats:', error)
    } finally {
      setLoading(false)
    }
  }

  return (
    <Spin spinning={loading}>
      <div>
        <Row gutter={[16, 16]} className="mb-6">
          <Col xs={12} sm={12} md={6}>
            <Card className="stat-card">
              <Statistic
                title="总销售额"
                value={stats.totalSales}
                precision={2}
                prefix={<DollarOutlined />}
                suffix="元"
                valueStyle={{ color: '#52c41a' }}
              />
              <div className="mt-2 text-sm text-gray-500">
                今日: <span className="text-green-500">¥{stats.todaySales.toLocaleString()}</span>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={12} md={6}>
            <Card className="stat-card">
              <Statistic
                title="总订单数"
                value={stats.totalOrders}
                prefix={<ShoppingOutlined />}
                valueStyle={{ color: '#1890ff' }}
              />
              <div className="mt-2 text-sm text-gray-500">
                今日: <span className="text-blue-500">{stats.todayOrders} 笔</span>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={12} md={6}>
            <Card className="stat-card">
              <Statistic
                title="总用户数"
                value={stats.totalUsers}
                prefix={<UserOutlined />}
                valueStyle={{ color: '#722ed1' }}
              />
              <div className="mt-2 text-sm text-gray-500">
                今日新增: <span className="text-purple-500">{stats.newUsersToday} 人</span>
              </div>
            </Card>
          </Col>
          <Col xs={12} sm={12} md={6}>
            <Card className="stat-card">
              <Statistic
                title="商品总数"
                value={stats.totalProducts}
                prefix={<RiseOutlined />}
                valueStyle={{ color: '#fa8c16' }}
              />
              <div className="mt-2 text-sm text-gray-500">
                库存不足: <span className="text-red-500">{stats.lowStockProducts} 件</span>
              </div>
            </Card>
          </Col>
        </Row>

        <Row gutter={[16, 16]}>
          <Col xs={24} md={14}>
            <Card title="销售趋势">
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={salesData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Line type="monotone" dataKey="sales" stroke="#52c41a" strokeWidth={2} name="销售额" />
                  <Line type="monotone" dataKey="orders" stroke="#1890ff" strokeWidth={2} name="订单数" />
                </LineChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col xs={24} md={10}>
            <Card title="商品分类占比">
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={categoryData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                    outerRadius={100}
                    fill="#8884d8"
                    dataKey="value"
                  >
                    {categoryData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </Card>
          </Col>
        </Row>

        <Row gutter={[16, 16]} className="mt-4">
          <Col xs={24} md={14}>
            <Card title="商品销售排行">
              <ResponsiveContainer width="100%" height={280}>
                <BarChart data={productSalesData} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" />
                  <YAxis dataKey="name" type="category" width={80} />
                  <Tooltip />
                  <Bar dataKey="sales" fill="#52c41a" radius={[0, 4, 4, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </Card>
          </Col>
          <Col xs={24} md={10}>
            <Card title="最新订单">
              <List
                dataSource={recentOrders}
                renderItem={(item) => (
                  <List.Item>
                    <List.Item.Meta
                      avatar={<Avatar icon={<ShoppingOutlined />} />}
                      title={
                        <div className="flex justify-between items-center">
                          <span className="text-sm">{item.orderNo}</span>
                          {getStatusTag(item.status)}
                        </div>
                      }
                      description={
                        <div className="flex justify-between">
                          <span className="text-gray-500">{item.user}</span>
                          <span className="text-red-500 font-bold">¥{item.amount}</span>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>
          </Col>
        </Row>
      </div>
    </Spin>
  )
}

export default Dashboard
