import React, { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Form, Input, Button, Card, message, Checkbox } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons'
import { useAppDispatch } from '@/store'
import { login } from '@/store/userSlice'
import { authApi } from '@/services/api'

const Login: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [loading, setLoading] = useState(false)
  const [form] = Form.useForm()

  const onFinish = async (values: { username: string; password: string; remember: boolean }) => {
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
        navigate('/')
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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-blue-50 py-12 px-4">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-green-600 mb-2">🌿 农产品商城</h1>
          <p className="text-gray-500">欢迎回来，请登录您的账户</p>
        </div>

        <Form
          form={form}
          name="login"
          initialValues={{ remember: true }}
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
            <div className="flex justify-between items-center">
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>记住我</Checkbox>
              </Form.Item>
              <Link to="/forgot-password" className="text-green-600">
                忘记密码?
              </Link>
            </div>
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              className="w-full h-12 text-lg"
              loading={loading}
              style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
            >
              登录
            </Button>
          </Form.Item>

          <div className="text-center">
            <span className="text-gray-500">还没有账户？</span>
            <Link to="/register" className="text-green-600 ml-1 font-medium">
              立即注册
            </Link>
          </div>
        </Form>

        <div className="mt-6 pt-6 border-t border-gray-200">
          <p className="text-center text-gray-400 text-sm">
            测试账号: admin / 123456
          </p>
        </div>
      </Card>
    </div>
  )
}

export default Login
