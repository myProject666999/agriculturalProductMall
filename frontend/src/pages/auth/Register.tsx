import React, { useState } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { Form, Input, Button, Card, message, InputNumber, Row, Col } from 'antd'
import { UserOutlined, LockOutlined, MailOutlined } from '@ant-design/icons'
import { authApi } from '@/services/api'

const Register: React.FC = () => {
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [codeLoading, setCodeLoading] = useState(false)
  const [countdown, setCountdown] = useState(0)
  const [form] = Form.useForm()

  React.useEffect(() => {
    let timer: NodeJS.Timeout
    if (countdown > 0) {
      timer = setTimeout(() => setCountdown(countdown - 1), 1000)
    }
    return () => clearTimeout(timer)
  }, [countdown])

  const handleSendCode = async () => {
    const email = form.getFieldValue('email')
    if (!email) {
      message.warning('请先输入邮箱')
      return
    }

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
      message.warning('请输入有效的邮箱地址')
      return
    }

    setCodeLoading(true)
    try {
      const res: any = await authApi.sendCode(email)
      if (res.code === 0) {
        message.success('验证码已发送到您的邮箱')
        setCountdown(60)
      } else {
        message.error(res.message || '发送失败')
      }
    } catch (error: any) {
      message.error(error.message || '发送失败，请稍后重试')
    } finally {
      setCodeLoading(false)
    }
  }

  const onFinish = async (values: {
    username: string
    email: string
    password: string
    confirmPassword: string
    code: string
  }) => {
    if (values.password !== values.confirmPassword) {
      message.error('两次输入的密码不一致')
      return
    }

    setLoading(true)
    try {
      const res: any = await authApi.register({
        username: values.username,
        email: values.email,
        password: values.password,
        confirmPassword: values.confirmPassword,
        code: values.code,
      })

      if (res.code === 0) {
        message.success('注册成功，请登录')
        navigate('/login')
      } else {
        message.error(res.message || '注册失败')
      }
    } catch (error: any) {
      message.error(error.message || '注册失败，请稍后重试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-blue-50 py-12 px-4">
      <Card className="w-full max-w-md shadow-xl">
        <div className="text-center mb-8">
          <h1 className="text-3xl font-bold text-green-600 mb-2">🌿 农产品商城</h1>
          <p className="text-gray-500">创建您的账户，开始购物之旅</p>
        </div>

        <Form
          form={form}
          name="register"
          onFinish={onFinish}
          size="large"
          scrollToFirstError
        >
          <Form.Item
            name="username"
            rules={[
              { required: true, message: '请输入用户名' },
              { min: 3, message: '用户名至少3个字符' },
              { max: 20, message: '用户名最多20个字符' },
            ]}
          >
            <Input prefix={<UserOutlined />} placeholder="用户名 (3-20个字符)" />
          </Form.Item>

          <Form.Item
            name="email"
            rules={[
              { required: true, message: '请输入邮箱' },
              { type: 'email', message: '请输入有效的邮箱地址' },
            ]}
          >
            <Input prefix={<MailOutlined />} placeholder="邮箱地址" />
          </Form.Item>

          <Form.Item
            name="code"
            rules={[{ required: true, message: '请输入验证码' }]}
          >
            <Row gutter={8}>
              <Col span={16}>
                <Input placeholder="邮箱验证码" />
              </Col>
              <Col span={8}>
                <Button
                  type="default"
                  className="w-full h-full"
                  loading={codeLoading}
                  disabled={countdown > 0}
                  onClick={handleSendCode}
                >
                  {countdown > 0 ? `${countdown}s` : '获取验证码'}
                </Button>
              </Col>
            </Row>
          </Form.Item>

          <Form.Item
            name="password"
            rules={[
              { required: true, message: '请输入密码' },
              { min: 6, message: '密码至少6个字符' },
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="密码 (至少6个字符)" />
          </Form.Item>

          <Form.Item
            name="confirmPassword"
            dependencies={['password']}
            rules={[
              { required: true, message: '请确认密码' },
              ({ getFieldValue }) => ({
                validator(_, value) {
                  if (!value || getFieldValue('password') === value) {
                    return Promise.resolve()
                  }
                  return Promise.reject(new Error('两次输入的密码不一致'))
                },
              }),
            ]}
          >
            <Input.Password prefix={<LockOutlined />} placeholder="确认密码" />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              className="w-full h-12 text-lg"
              loading={loading}
              style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
            >
              注册
            </Button>
          </Form.Item>

          <div className="text-center">
            <span className="text-gray-500">已有账户？</span>
            <Link to="/login" className="text-green-600 ml-1 font-medium">
              立即登录
            </Link>
          </div>
        </Form>

        <div className="mt-6 pt-6 border-t border-gray-200">
          <p className="text-center text-gray-400 text-sm">
            注册即表示同意<a href="#" className="text-green-600">用户协议</a>和
            <a href="#" className="text-green-600">隐私政策</a>
          </p>
        </div>
      </Card>
    </div>
  )
}

export default Register
