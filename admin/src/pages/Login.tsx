import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Form, Input, Button, Card, message } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useAppDispatch } from '@/store'
import { login } from '@/store/userSlice'
import { authApi } from '@/services/api'

const Login: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const onFinish = async (values: { username: string; password: string }) => {
    setLoading(true)
    try {
      const res: any = await authApi.login({
        username: values.username,
        password: values.password,
      })

      if (res.code === 0) {
        const userData = res.data.user || res.data
        const tokenData = res.data.token || res.data.accessToken

        dispatch(
          login({
            user: userData,
            token: tokenData,
          })
        )
        message.success('登录成功')
        navigate('/dashboard')
      } else {
        message.error(res.message || '登录失败')
      }
    } catch (error: any) {
      message.error(error.message || '登录失败，请检查用户名和密码')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="login-container">
      <Card className="login-card">
        <div className="login-title">
          <h1>🌿 农产品商城</h1>
          <p>管理系统后台</p>
        </div>

        <Form
          form={form}
          name="admin-login"
          onFinish={onFinish}
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              className="w-full h-12 text-lg"
              loading={loading}
            >
              登录
            </Button>
          </Form.Item>
        </Form>

        <div className="text-center text-gray-400 text-sm">
          <p>管理员账号: admin / 123456</p>
          <p>商户账号: merchant / 123456</p>
        </div>
      </Card>
    </div>
  )
}

export default Login
