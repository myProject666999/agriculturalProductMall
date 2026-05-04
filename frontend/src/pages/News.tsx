import React, { useState, useEffect } from 'react'
import { Card, List, Empty, Spin, Tag, Divider, Typography, Pagination } from 'antd'
import { CalendarOutlined, EyeOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { newsApi } from '@/services/api'

const { Title, Paragraph } = Typography

const mockNews = [
  {
    id: 1,
    title: '有机蔬菜种植技术分享：如何种出健康又美味的蔬菜',
    summary: '本文将为您详细介绍有机蔬菜的种植技术，包括土壤选择、肥料使用、病虫害防治等关键环节...',
    content: '完整的新闻内容...',
    cover: '',
    author: '农技专家',
    viewCount: 1256,
    createdAt: '2024-01-15 10:30:00',
  },
  {
    id: 2,
    title: '新鲜水果直送：从果园到餐桌的最短距离',
    summary: '我们的水果直接来自合作果园，确保新鲜度和品质。了解我们的配送流程和质量保证措施...',
    content: '完整的新闻内容...',
    cover: '',
    author: '品质保障部',
    viewCount: 892,
    createdAt: '2024-01-14 15:20:00',
  },
  {
    id: 3,
    title: '春节特惠活动预告：精选农产品低至5折起',
    summary: '为回馈新老用户，春节期间我们将推出一系列特惠活动。精选农产品低至5折，更有满减优惠...',
    content: '完整的新闻内容...',
    cover: '',
    author: '市场运营部',
    viewCount: 2134,
    createdAt: '2024-01-13 09:00:00',
  },
  {
    id: 4,
    title: '农产品溯源系统上线：让每一份食材都有迹可循',
    summary: '我们的农产品溯源系统正式上线，用户可以通过扫描二维码查看产品的详细信息，包括产地、种植过程等...',
    content: '完整的新闻内容...',
    cover: '',
    author: '技术研发部',
    viewCount: 678,
    createdAt: '2024-01-12 14:45:00',
  },
]

const News: React.FC = () => {
  const [loading, setLoading] = useState(false)
  const [news, setNews] = useState(mockNews)
  const [current, setCurrent] = useState(1)
  const [pageSize] = useState(10)

  useEffect(() => {
    loadNews()
  }, [current])

  const loadNews = async () => {
    setLoading(true)
    try {
      const res: any = await newsApi.getList({
        page: current,
        pageSize,
      }).catch(() => ({ code: 0, data: { list: mockNews, total: mockNews.length } }))

      if (res.code === 0 && res.data?.list) {
        setNews(res.data.list)
      }
    } catch (error) {
      console.error('Failed to load news:', error)
    } finally {
      setLoading(false)
    }
  }

  const handlePageChange = (page: number) => {
    setCurrent(page)
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-6">
      <Card>
        <Title level={3} style={{ margin: 0, marginBottom: 16 }}>
          资讯中心
        </Title>
        <Divider style={{ margin: '16px 0' }} />

        <Spin spinning={loading}>
          {news && news.length > 0 ? (
            <List
              dataSource={news}
              renderItem={(item: any) => (
                <List.Item key={item.id} className="cursor-pointer hover:bg-gray-50 -mx-4 px-4 py-4 rounded">
                  <List.Item.Meta
                    title={
                      <div className="flex items-center gap-2">
                        <span className="text-lg font-medium text-gray-800 hover:text-green-600">
                          {item.title}
                        </span>
                        <Tag color="green">最新</Tag>
                      </div>
                    }
                    description={
                      <div>
                        <Paragraph ellipsis={{ rows: 2 }} className="text-gray-600 mt-2">
                          {item.summary}
                        </Paragraph>
                        <div className="flex items-center gap-6 text-gray-400 text-sm">
                          <span>
                            <CalendarOutlined className="mr-1" />
                            {dayjs(item.createdAt).format('YYYY-MM-DD HH:mm')}
                          </span>
                          <span>
                            <EyeOutlined className="mr-1" />
                            {item.viewCount} 次浏览
                          </span>
                          <span>作者: {item.author}</span>
                        </div>
                      </div>
                    }
                  />
                </List.Item>
              )}
            />
          ) : (
            <Empty description="暂无资讯" />
          )}

          {news && news.length > 0 && (
            <div className="flex justify-center mt-8">
              <Pagination
                current={current}
                pageSize={pageSize}
                total={mockNews.length}
                onChange={handlePageChange}
                showSizeChanger={false}
              />
            </div>
          )}
        </Spin>
      </Card>
    </div>
  )
}

export default News
