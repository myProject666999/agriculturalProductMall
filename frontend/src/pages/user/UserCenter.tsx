import React, { useState, useEffect } from 'react'
import { useSearchParams } from 'react-router-dom'
import { Layout, Menu, Card, Avatar, Form, Input, Button, message, Tabs, TabsProps, List, Table, Empty, Tag, Space, Popconfirm, Modal, InputNumber, Select } from 'antd'
import { UserOutlined, LockOutlined, EnvironmentOutlined, HeartOutlined, ShoppingCartOutlined, EditOutlined, DeleteOutlined, PlusOutlined } from '@ant-design/icons'
import { useAppSelector, useAppDispatch } from '@/store'
import { updateUser } from '@/store/userSlice'
import { userApi } from '@/services/api'
import type { MenuProps } from 'antd'

const { Sider, Content } = Layout
const { Option } = Select

const UserCenter: React.FC = () => {
  const [searchParams, setSearchParams] = useSearchParams()
  const user = useAppSelector((state) => state.user.user)
  const dispatch = useAppDispatch()
  const [activeKey, setActiveKey] = useState(searchParams.get('tab') || 'profile')
  const [loading, setLoading] = useState(false)
  const [addresses, setAddresses] = useState<any[]>([])
  const [favorites, setFavorites] = useState<any[]>([])
  const [isAddressModalOpen, setIsAddressModalOpen] = useState(false)
  const [editingAddress, setEditingAddress] = useState<any>(null)
  const [form] = Form.useForm()
  const [passwordForm] = Form.useForm()

  useEffect(() => {
    const tab = searchParams.get('tab') || 'profile'
    setActiveKey(tab)
    if (tab === 'addresses') {
      loadAddresses()
    } else if (tab === 'favorites') {
      loadFavorites()
    }
  }, [searchParams])

  const loadAddresses = async () => {
    try {
      const res: any = await userApi.getAddresses()
      if (res.code === 0) {
        setAddresses(res.data || [])
      }
    } catch (error) {
      console.error('Failed to load addresses:', error)
    }
  }

  const loadFavorites = async () => {
    try {
      const res: any = await userApi.getFavorites()
      if (res.code === 0) {
        setFavorites(res.data || [])
      }
    } catch (error) {
      console.error('Failed to load favorites:', error)
    }
  }

  const handleMenuClick: MenuProps['onClick'] = (e) => {
    setActiveKey(e.key)
    setSearchParams({ tab: e.key })
    if (e.key === 'addresses') {
      loadAddresses()
    } else if (e.key === 'favorites') {
      loadFavorites()
    }
  }

  const handleProfileSubmit = async (values: any) => {
    setLoading(true)
    try {
      const res: any = await userApi.updateInfo(values)
      if (res.code === 0) {
        message.success('修改成功')
        if (res.data) {
          dispatch(updateUser(res.data))
        }
      } else {
        message.error(res.message || '修改失败')
      }
    } catch (error) {
      message.error('修改失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  const handlePasswordSubmit = async (values: any) => {
    if (values.newPassword !== values.confirmPassword) {
      message.error('两次输入的新密码不一致')
      return
    }
    setLoading(true)
    try {
      const res: any = await userApi.updatePassword({
        oldPassword: values.oldPassword,
        newPassword: values.newPassword,
      })
      if (res.code === 0) {
        message.success('密码修改成功')
        passwordForm.resetFields()
      } else {
        message.error(res.message || '修改失败')
      }
    } catch (error) {
      message.error('修改失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  const handleAddAddress = () => {
    setEditingAddress(null)
    form.resetFields()
    setIsAddressModalOpen(true)
  }

  const handleEditAddress = (address: any) => {
    setEditingAddress(address)
    form.setFieldsValue(address)
    setIsAddressModalOpen(true)
  }

  const handleDeleteAddress = async (id: number) => {
    try {
      const res: any = await userApi.deleteAddress(id)
      if (res.code === 0) {
        message.success('删除成功')
        loadAddresses()
      } else {
        message.error(res.message || '删除失败')
      }
    } catch (error) {
      message.error('删除失败，请稍后重试')
    }
  }

  const handleAddressSubmit = async (values: any) => {
    setLoading(true)
    try {
      let res: any
      if (editingAddress) {
        res = await userApi.updateAddress(editingAddress.id, values)
      } else {
        res = await userApi.addAddress(values)
      }
      if (res.code === 0) {
        message.success(editingAddress ? '修改成功' : '添加成功')
        setIsAddressModalOpen(false)
        loadAddresses()
      } else {
        message.error(res.message || '操作失败')
      }
    } catch (error) {
      message.error('操作失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  const handleRemoveFavorite = async (productId: number) => {
    try {
      const res: any = await userApi.removeFavorite(productId)
      if (res.code === 0) {
        message.success('已取消收藏')
        loadFavorites()
      } else {
        message.error(res.message || '操作失败')
      }
    } catch (error) {
      message.error('操作失败，请稍后重试')
    }
  }

  const menuItems: MenuProps['items'] = [
    { key: 'profile', icon: <UserOutlined />, label: '个人信息' },
    { key: 'password', icon: <LockOutlined />, label: '修改密码' },
    { key: 'addresses', icon: <EnvironmentOutlined />, label: '收货地址' },
    { key: 'favorites', icon: <HeartOutlined />, label: '我的收藏' },
  ]

  const addressColumns = [
    {
      title: '收货人',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '联系电话',
      dataIndex: 'phone',
      key: 'phone',
    },
    {
      title: '收货地址',
      key: 'address',
      render: (_: any, record: any) => (
        <span>
          {record.province} {record.city} {record.district} {record.detail}
        </span>
      ),
    },
    {
      title: '默认地址',
      dataIndex: 'isDefault',
      key: 'isDefault',
      render: (isDefault: boolean) => (
        isDefault ? <Tag color="green">是</Tag> : <Tag>否</Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button type="link" size="small" icon={<EditOutlined />} onClick={() => handleEditAddress(record)}>
            编辑
          </Button>
          <Popconfirm
            title="确定要删除这个地址吗？"
            onConfirm={() => handleDeleteAddress(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="text" danger size="small" icon={<DeleteOutlined />}>
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <Layout style={{ background: 'transparent' }}>
        <Sider width={200} style={{ background: '#fff', borderRadius: 8 }} className="sidebar-menu">
          <Card className="text-center mb-4">
            <Avatar size={64} icon={<UserOutlined />} src={user?.avatar} />
            <p className="mt-2 font-medium">{user?.username || '用户'}</p>
            <Tag color="green">{user?.role === 'merchant' ? '商户' : user?.role === 'admin' ? '管理员' : '普通用户'}</Tag>
          </Card>
          <Menu
            mode="inline"
            selectedKeys={[activeKey]}
            onClick={handleMenuClick}
            items={menuItems}
            style={{ borderRight: 'none' }}
          />
        </Sider>
        <Content className="content-area">
          {activeKey === 'profile' && (
            <Card title="个人信息">
              <Form
                layout="vertical"
                initialValues={{
                  username: user?.username,
                  email: user?.email,
                  phone: user?.phone,
                }}
                onFinish={handleProfileSubmit}
              >
                <Form.Item label="用户名" name="username">
                  <Input disabled />
                </Form.Item>
                <Form.Item label="邮箱" name="email">
                  <Input disabled />
                </Form.Item>
                <Form.Item label="手机号" name="phone">
                  <Input placeholder="请输入手机号" />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}>
                    保存修改
                  </Button>
                </Form.Item>
              </Form>
            </Card>
          )}

          {activeKey === 'password' && (
            <Card title="修改密码">
              <Form
                form={passwordForm}
                layout="vertical"
                onFinish={handlePasswordSubmit}
              >
                <Form.Item
                  label="原密码"
                  name="oldPassword"
                  rules={[{ required: true, message: '请输入原密码' }]}
                >
                  <Input.Password placeholder="请输入原密码" />
                </Form.Item>
                <Form.Item
                  label="新密码"
                  name="newPassword"
                  rules={[
                    { required: true, message: '请输入新密码' },
                    { min: 6, message: '密码至少6个字符' },
                  ]}
                >
                  <Input.Password placeholder="请输入新密码（至少6个字符）" />
                </Form.Item>
                <Form.Item
                  label="确认新密码"
                  name="confirmPassword"
                  dependencies={['newPassword']}
                  rules={[
                    { required: true, message: '请确认新密码' },
                    ({ getFieldValue }) => ({
                      validator(_, value) {
                        if (!value || getFieldValue('newPassword') === value) {
                          return Promise.resolve()
                        }
                        return Promise.reject(new Error('两次输入的密码不一致'))
                      },
                    }),
                  ]}
                >
                  <Input.Password placeholder="请再次输入新密码" />
                </Form.Item>
                <Form.Item>
                  <Button type="primary" htmlType="submit" loading={loading} style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}>
                    确认修改
                  </Button>
                </Form.Item>
              </Form>
            </Card>
          )}

          {activeKey === 'addresses' && (
            <Card
              title="收货地址"
              extra={
                <Button type="primary" icon={<PlusOutlined />} onClick={handleAddAddress} style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}>
                  添加地址
                </Button>
              }
            >
              {addresses.length > 0 ? (
                <Table
                  columns={addressColumns}
                  dataSource={addresses}
                  rowKey="id"
                  pagination={false}
                />
              ) : (
                <Empty description="暂无收货地址" />
              )}
            </Card>
          )}

          {activeKey === 'favorites' && (
            <Card title="我的收藏">
              {favorites.length > 0 ? (
                <List
                  grid={{ gutter: 16, xs: 1, sm: 2, md: 3, lg: 4, xl: 5 }}
                  dataSource={favorites}
                  renderItem={(item: any) => (
                    <List.Item>
                      <Card
                        hoverable
                        cover={
                          <img
                            alt={item.productName || item.name}
                            src={item.image || 'https://picsum.photos/200/200'}
                            style={{ height: 150, objectFit: 'cover' }}
                          />
                        }
                        actions={[
                          <Popconfirm
                            title="确定要取消收藏吗？"
                            onConfirm={() => handleRemoveFavorite(item.productId || item.id)}
                            okText="确定"
                            cancelText="取消"
                          >
                            <Button type="text" danger icon={<HeartOutlined />}>
                              取消收藏
                            </Button>
                          </Popconfirm>,
                        ]}
                      >
                        <Card.Meta
                          title={item.productName || item.name}
                          description={<span className="text-red-500 font-bold">¥{item.price}</span>}
                        />
                      </Card>
                    </List.Item>
                  )}
                />
              ) : (
                <Empty description="暂无收藏商品" />
              )}
            </Card>
          )}
        </Content>
      </Layout>

      <Modal
        title={editingAddress ? '编辑收货地址' : '添加收货地址'}
        open={isAddressModalOpen}
        onCancel={() => setIsAddressModalOpen(false)}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleAddressSubmit}
        >
          <Form.Item
            label="收货人"
            name="name"
            rules={[{ required: true, message: '请输入收货人' }]}
          >
            <Input placeholder="请输入收货人姓名" />
          </Form.Item>
          <Form.Item
            label="联系电话"
            name="phone"
            rules={[{ required: true, message: '请输入联系电话' }]}
          >
            <Input placeholder="请输入联系电话" />
          </Form.Item>
          <Form.Item
            label="省份"
            name="province"
            rules={[{ required: true, message: '请选择省份' }]}
          >
            <Select placeholder="请选择省份">
              <Option value="北京市">北京市</Option>
              <Option value="上海市">上海市</Option>
              <Option value="广东省">广东省</Option>
              <Option value="浙江省">浙江省</Option>
              <Option value="江苏省">江苏省</Option>
            </Select>
          </Form.Item>
          <Form.Item
            label="城市"
            name="city"
            rules={[{ required: true, message: '请输入城市' }]}
          >
            <Input placeholder="请输入城市" />
          </Form.Item>
          <Form.Item
            label="区县"
            name="district"
            rules={[{ required: true, message: '请输入区县' }]}
          >
            <Input placeholder="请输入区县" />
          </Form.Item>
          <Form.Item
            label="详细地址"
            name="detail"
            rules={[{ required: true, message: '请输入详细地址' }]}
          >
            <Input.TextArea placeholder="请输入详细地址" rows={2} />
          </Form.Item>
          <Form.Item
            label="是否默认地址"
            name="isDefault"
            valuePropName="checked"
          >
            <Select>
              <Option value={true}>是</Option>
              <Option value={false}>否</Option>
            </Select>
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading} block style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}>
              保存
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default UserCenter
