import React from 'react'
import { Layout, Row, Col, Space } from 'antd'
import { Link } from 'react-router-dom'

const { Footer: AntFooter } = Layout

const Footer: React.FC = () => {
  return (
    <AntFooter className="footer-container">
      <div className="max-w-7xl mx-auto">
        <Row gutter={[32, 32]}>
          <Col xs={24} md={6}>
            <h3 className="text-white text-lg font-bold mb-4">关于我们</h3>
            <p className="text-gray-400 text-sm leading-relaxed">
              农产品商城致力于为您提供最优质、最新鲜的农产品。我们严格把控质量，让您吃得放心、吃得健康。
            </p>
          </Col>
          <Col xs={24} md={6}>
            <h3 className="text-white text-lg font-bold mb-4">快速链接</h3>
            <Space direction="vertical" size="small">
              <Link to="/" className="text-gray-400 hover:text-white transition-colors">
                首页
              </Link>
              <Link to="/products" className="text-gray-400 hover:text-white transition-colors">
                商品中心
              </Link>
              <Link to="/news" className="text-gray-400 hover:text-white transition-colors">
                资讯中心
              </Link>
              <Link to="/cart" className="text-gray-400 hover:text-white transition-colors">
                购物车
              </Link>
            </Space>
          </Col>
          <Col xs={24} md={6}>
            <h3 className="text-white text-lg font-bold mb-4">帮助中心</h3>
            <Space direction="vertical" size="small">
              <a href="#" className="text-gray-400 hover:text-white transition-colors">
                新手上路
              </a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">
                支付方式
              </a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">
                配送说明
              </a>
              <a href="#" className="text-gray-400 hover:text-white transition-colors">
                售后服务
              </a>
            </Space>
          </Col>
          <Col xs={24} md={6}>
            <h3 className="text-white text-lg font-bold mb-4">联系我们</h3>
            <Space direction="vertical" size="small" className="text-gray-400">
              <p>客服热线：400-123-4567</p>
              <p>工作时间：9:00 - 21:00</p>
              <p>邮箱：service@agrimall.com</p>
              <p>地址：北京市朝阳区科技园</p>
            </Space>
          </Col>
        </Row>
        <div className="border-t border-gray-700 mt-8 pt-6 text-center text-gray-500 text-sm">
          <p>Copyright © 2024 农产品商城. All Rights Reserved.</p>
          <p className="mt-2">ICP备案号：京ICP备12345678号-1</p>
        </div>
      </div>
    </AntFooter>
  )
}

export default Footer
