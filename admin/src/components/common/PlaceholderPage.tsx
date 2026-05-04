import React from 'react'
import { Result, Button } from 'antd'
import { SmileOutlined } from '@ant-design/icons'

interface PlaceholderPageProps {
  title?: string
}

const PlaceholderPage: React.FC<PlaceholderPageProps> = ({ title = '页面开发中' }) => {
  return (
    <Result
      icon={<SmileOutlined />}
      title={title}
      subTitle="该功能正在开发中，敬请期待..."
      extra={
        <Button type="primary" onClick={() => window.history.back()}>
          返回上一页
        </Button>
      }
    />
  )
}

export default PlaceholderPage
