import React, { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { Carousel, Card, Row, Col, Button, message, Spin, Empty, Tag, Tabs, TabsProps } from 'antd'
import { RightOutlined, ShoppingCartOutlined, HeartOutlined } from '@ant-design/icons'
import { bannerApi, productApi, newsApi } from '@/services/api'
import { useAppDispatch } from '@/store'
import { addItem } from '@/store/cartSlice'

const { Meta } = Card

const Home: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [banners, setBanners] = useState<any[]>([])
  const [recommendProducts, setRecommendProducts] = useState<any[]>([])
  const [hotProducts, setHotProducts] = useState<any[]>([])
  const [categories, setCategories] = useState<any[]>([])
  const [news, setNews] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    loadData()
  }, [])

  const loadData = async () => {
    setLoading(true)
    try {
      const [bannerRes, recommendRes, categoryRes, newsRes] = await Promise.all([
        bannerApi.getList(),
        productApi.getRecommend(),
        productApi.getCategories(),
        newsApi.getList({ page: 1, pageSize: 5 }),
      ])

      if ((bannerRes as any).code === 0) {
        setBanners((bannerRes as any).data || [])
      }
      if ((recommendRes as any).code === 0) {
        setRecommendProducts((recommendRes as any).data || [])
      }
      if ((categoryRes as any).code === 0) {
        setCategories((categoryRes as any).data || [])
      }
      if ((newsRes as any).code === 0) {
        setNews((newsRes as any).data?.list || [])
      }
    } catch (error) {
      console.error('Failed to load home data:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = (product: any) => {
    const cartItem = {
      id: product.id,
      productId: product.id,
      productName: product.name,
      price: product.price,
      quantity: 1,
      image: product.image || 'https://picsum.photos/200/200',
      stock: product.stock,
    }
    dispatch(addItem(cartItem))
    message.success('已加入购物车')
  }

  const tabItems: TabsProps['items'] = [
    {
      key: 'recommend',
      label: '个性化推荐',
      children: (
        <Row gutter={[16, 16]}>
          {recommendProducts.length > 0 ? (
            recommendProducts.slice(0, 8).map((product) => (
              <Col xs={12} sm={8} md={6} lg={4} key={product.id}>
                <Card
                  hoverable
                  className="product-card"
                  cover={
                    <img
                      alt={product.name}
                      src={product.image || 'https://picsum.photos/200/200'}
                      className="product-image"
                      onClick={() => navigate(`/products/${product.id}`)}
                    />
                  }
                  actions={[
                    <ShoppingCartOutlined
                      key="cart"
                      onClick={() => handleAddToCart(product)}
                    />,
                    <HeartOutlined key="heart" />,
                  ]}
                >
                  <Meta
                    title={
                      <div
                        className="truncate"
                        onClick={() => navigate(`/products/${product.id}`)}
                      >
                        {product.name}
                      </div>
                    }
                    description={
                      <div className="price-tag">
                        {product.price}
                        <span className="text-gray-400 text-sm line-through ml-2">
                          ¥{product.originalPrice || product.price * 1.2}
                        </span>
                      </div>
                    }
                  />
                  {product.isNew && <Tag color="green">新品</Tag>}
                  {product.isHot && <Tag color="red">热卖</Tag>}
                </Card>
              </Col>
            ))
          ) : (
            <Col span={24}>
              <Empty description="暂无推荐商品" />
            </Col>
          )}
        </Row>
      ),
    },
    {
      key: 'hot',
      label: '热销商品',
      children: (
        <Row gutter={[16, 16]}>
          {recommendProducts.length > 0 ? (
            recommendProducts.slice(0, 8).map((product) => (
              <Col xs={12} sm={8} md={6} lg={4} key={product.id}>
                <Card
                  hoverable
                  className="product-card"
                  cover={
                    <img
                      alt={product.name}
                      src={product.image || 'https://picsum.photos/200/200'}
                      className="product-image"
                      onClick={() => navigate(`/products/${product.id}`)}
                    />
                  }
                  actions={[
                    <ShoppingCartOutlined
                      key="cart"
                      onClick={() => handleAddToCart(product)}
                    />,
                    <HeartOutlined key="heart" />,
                  ]}
                >
                  <Meta
                    title={
                      <div
                        className="truncate"
                        onClick={() => navigate(`/products/${product.id}`)}
                      >
                        {product.name}
                      </div>
                    }
                    description={
                      <div className="price-tag">
                        {product.price}
                        <span className="text-gray-400 text-sm line-through ml-2">
                          ¥{product.originalPrice || product.price * 1.2}
                        </span>
                      </div>
                    }
                  />
                </Card>
              </Col>
            ))
          ) : (
            <Col span={24}>
              <Empty description="暂无热销商品" />
            </Col>
          )}
        </Row>
      ),
    },
  ]

  const carouselItems = [
    {
      id: 1,
      image: 'https://picsum.photos/1200/400?random=1',
      title: '新鲜直达，有机蔬菜限时特惠',
      description: '全场蔬菜满99减30',
    },
    {
      id: 2,
      image: 'https://picsum.photos/1200/400?random=2',
      title: '水果节狂欢季',
      description: '进口水果第二件半价',
    },
    {
      id: 3,
      image: 'https://picsum.photos/1200/400?random=3',
      title: '新用户专享优惠',
      description: '注册即送100元优惠券',
    },
  ]

  return (
    <div className="pb-8">
      <Spin spinning={loading}>
        <Carousel autoplay effect="fade" className="h-96">
          {carouselItems.map((item) => (
            <div
              key={item.id}
              className="carousel-slide"
              style={{
                backgroundImage: `url(${item.image})`,
              }}
            >
              <div className="text-center text-white bg-black bg-opacity-40 p-8 rounded-lg">
                <h2 className="text-4xl font-bold mb-4">{item.title}</h2>
                <p className="text-xl mb-6">{item.description}</p>
                <Button type="primary" size="large" onClick={() => navigate('/products')}>
                  立即抢购
                </Button>
              </div>
            </div>
          ))}
        </Carousel>

        <div className="max-w-7xl mx-auto px-4">
          <div className="recommend-section">
            <div className="flex justify-between items-center mb-6">
              <h2 className="section-title">商品分类</h2>
            </div>
            <Row gutter={[16, 16]}>
              {categories.slice(0, 8).map((category) => (
                <Col xs={6} sm={4} md={3} key={category.id}>
                  <div
                    className="category-item bg-gray-50"
                    onClick={() => navigate(`/products?categoryId=${category.id}`)}
                  >
                    <img
                      src={category.icon || `https://picsum.photos/64/64?random=${category.id}`}
                      alt={category.name}
                      className="w-16 h-16 mx-auto mb-2 rounded-full"
                    />
                    <p className="font-medium">{category.name}</p>
                  </div>
                </Col>
              ))}
            </Row>
          </div>

          <div className="recommend-section">
            <div className="flex justify-between items-center mb-6">
              <h2 className="section-title">精选商品</h2>
              <Button type="link" onClick={() => navigate('/products')}>
                查看更多 <RightOutlined />
              </Button>
            </div>
            <Tabs defaultActiveKey="recommend" items={tabItems} />
          </div>

          <div className="news-section">
            <div className="flex justify-between items-center mb-6">
              <h2 className="section-title">最新资讯</h2>
              <Button type="link" onClick={() => navigate('/news')}>
                查看更多 <RightOutlined />
              </Button>
            </div>
            <Row gutter={[16, 16]}>
              {news.length > 0 ? (
                news.slice(0, 4).map((item) => (
                  <Col xs={24} sm={12} md={6} key={item.id}>
                    <Card
                      hoverable
                      cover={
                        <img
                          alt={item.title}
                          src={item.cover || `https://picsum.photos/300/200?random=${item.id}`}
                          style={{ height: 150, objectFit: 'cover' }}
                        />
                      }
                    >
                      <Meta
                        title={item.title}
                        description={
                          <div className="text-gray-500 text-sm">
                            {item.summary?.slice(0, 50)}...
                          </div>
                        }
                      />
                    </Card>
                  </Col>
                ))
              ) : (
                <Col span={24}>
                  <Empty description="暂无资讯" />
                </Col>
              )}
            </Row>
          </div>
        </div>
      </Spin>
    </div>
  )
}

export default Home
