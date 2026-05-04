import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Row, Col, Image, Button, InputNumber, message, Tag, Spin, Empty, Tabs, TabsProps, Descriptions, Rate, List, Avatar } from 'antd'
import { ShoppingCartOutlined, HeartOutlined, StarOutlined, StarFilled } from '@ant-design/icons'
import { productApi } from '@/services/api'
import { useAppDispatch } from '@/store'
import { addItem } from '@/store/cartSlice'

const ProductDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [product, setProduct] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [quantity, setQuantity] = useState(1)
  const [activeTab, setActiveTab] = useState('detail')

  useEffect(() => {
    if (id) {
      loadProduct(Number(id))
    }
  }, [id])

  const loadProduct = async (productId: number) => {
    setLoading(true)
    try {
      const res: any = await productApi.getDetail(productId)
      if (res.code === 0) {
        setProduct(res.data)
      }
    } catch (error) {
      console.error('Failed to load product:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleAddToCart = () => {
    if (!product) return
    const cartItem = {
      id: product.id,
      productId: product.id,
      productName: product.name,
      price: product.price,
      quantity: quantity,
      image: product.image || 'https://picsum.photos/200/200',
      stock: product.stock,
    }
    dispatch(addItem(cartItem))
    message.success(`已加入购物车，共 ${quantity} 件`)
  }

  const handleBuyNow = () => {
    handleAddToCart()
    navigate('/cart')
  }

  const tabItems: TabsProps['items'] = [
    {
      key: 'detail',
      label: '商品详情',
      children: (
        <div className="p-4">
          <Descriptions column={2}>
            <Descriptions.Item label="商品名称">{product?.name}</Descriptions.Item>
            <Descriptions.Item label="商品分类">{product?.categoryName || '未分类'}</Descriptions.Item>
            <Descriptions.Item label="商品品牌">{product?.brand || '无'}</Descriptions.Item>
            <Descriptions.Item label="产地">{product?.origin || '中国大陆'}</Descriptions.Item>
            <Descriptions.Item label="保质期">{product?.shelfLife || '12个月'}</Descriptions.Item>
            <Descriptions.Item label="储存方式">{product?.storageMethod || '阴凉干燥处'}</Descriptions.Item>
          </Descriptions>
          <div className="mt-6">
            <h3 className="text-lg font-bold mb-2">商品描述</h3>
            <div className="text-gray-600 leading-relaxed">
              {product?.description || '暂无详细描述。这款农产品精选自优质产地，严格把控品质，为您带来新鲜美味的健康食品。'}
            </div>
          </div>
        </div>
      ),
    },
    {
      key: 'specs',
      label: '规格参数',
      children: (
        <div className="p-4">
          <Descriptions column={1} bordered>
            <Descriptions.Item label="商品名称">{product?.name}</Descriptions.Item>
            <Descriptions.Item label="净含量">{product?.weight || '500g'}</Descriptions.Item>
            <Descriptions.Item label="产地">{product?.origin || '中国大陆'}</Descriptions.Item>
            <Descriptions.Item label="保质期">{product?.shelfLife || '12个月'}</Descriptions.Item>
            <Descriptions.Item label="储存方式">{product?.storageMethod || '阴凉干燥处'}</Descriptions.Item>
            <Descriptions.Item label="食用方法">{product?.usage || '洗净后即可食用或烹饪'}</Descriptions.Item>
          </Descriptions>
        </div>
      ),
    },
    {
      key: 'reviews',
      label: `用户评价 (${product?.reviewCount || 0})`,
      children: (
        <div className="p-4">
          <List
            itemLayout="horizontal"
            dataSource={[
              {
                id: 1,
                user: '张**',
                avatar: 'https://picsum.photos/48/48?random=1',
                rating: 5,
                content: '非常新鲜，质量很好，物流也很快，下次还会回购！',
                images: ['https://picsum.photos/100/100?random=1', 'https://picsum.photos/100/100?random=2'],
                time: '2024-01-15',
              },
              {
                id: 2,
                user: '李**',
                avatar: 'https://picsum.photos/48/48?random=2',
                rating: 4,
                content: '口感不错，很新鲜，就是价格有点小贵。',
                images: [],
                time: '2024-01-10',
              },
            ]}
            renderItem={(item) => (
              <List.Item>
                <List.Item.Meta
                  avatar={<Avatar src={item.avatar} />}
                  title={
                    <div className="flex items-center gap-2">
                      <span>{item.user}</span>
                      <Rate disabled value={item.rating} />
                    </div>
                  }
                  description={
                    <div>
                      <p className="text-gray-700">{item.content}</p>
                      {item.images.length > 0 && (
                        <div className="flex gap-2 mt-2">
                          {item.images.map((img, idx) => (
                            <Image key={idx} width={80} height={80} src={img} />
                          ))}
                        </div>
                      )}
                      <p className="text-gray-400 text-sm mt-2">{item.time}</p>
                    </div>
                  }
                />
              </List.Item>
            )}
          />
        </div>
      ),
    },
  ]

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <Spin spinning={loading}>
        {product ? (
          <>
            <Card className="mb-6">
              <Row gutter={32}>
                <Col xs={24} md={12} lg={10}>
                  <div className="bg-gray-50 rounded-lg p-4">
                    <Image
                      width={400}
                      height={400}
                      src={product.image || 'https://picsum.photos/400/400'}
                      alt={product.name}
                      className="object-contain"
                    />
                  </div>
                </Col>
                <Col xs={24} md={12} lg={14}>
                  <div className="h-full flex flex-col">
                    <div>
                      <h1 className="text-2xl font-bold mb-2">{product.name}</h1>
                      <p className="text-gray-500 mb-4">{product.summary || '精选优质农产品'}</p>
                      <div className="flex items-center gap-2 mb-4">
                        {product.isNew && <Tag color="green">新品</Tag>}
                        {product.isHot && <Tag color="red">热卖</Tag>}
                        {product.isRecommend && <Tag color="blue">推荐</Tag>}
                      </div>
                    </div>

                    <div className="bg-orange-50 rounded-lg p-4 mb-6">
                      <div className="flex items-baseline gap-2">
                        <span className="text-gray-500">售价</span>
                        <span className="text-3xl font-bold text-red-500">¥{product.price}</span>
                        {product.originalPrice && product.originalPrice > product.price && (
                          <span className="text-gray-400 line-through">¥{product.originalPrice}</span>
                        )}
                      </div>
                      <div className="flex gap-6 mt-2 text-sm text-gray-500">
                        <span>销量: {product.sales || 0}</span>
                        <span>评价: {product.reviewCount || 0}</span>
                        <span>库存: {product.stock || 0}</span>
                      </div>
                    </div>

                    <div className="space-y-4">
                      <div className="flex items-center gap-4">
                        <span className="text-gray-500 w-16">数量</span>
                        <InputNumber
                          min={1}
                          max={product.stock || 999}
                          value={quantity}
                          onChange={(val) => setQuantity(val || 1)}
                          size="large"
                        />
                        <span className="text-gray-400">库存 {product.stock || 0} 件</span>
                      </div>

                      <div className="flex gap-4 mt-6">
                        <Button
                          type="primary"
                          size="large"
                          icon={<ShoppingCartOutlined />}
                          onClick={handleAddToCart}
                          style={{ backgroundColor: '#52c41a', borderColor: '#52c41a' }}
                        >
                          加入购物车
                        </Button>
                        <Button
                          type="primary"
                          danger
                          size="large"
                          onClick={handleBuyNow}
                        >
                          立即购买
                        </Button>
                        <Button
                          size="large"
                          icon={<HeartOutlined />}
                        >
                          收藏
                        </Button>
                      </div>
                    </div>
                  </div>
                </Col>
              </Row>
            </Card>

            <Card>
              <Tabs
                activeKey={activeTab}
                onChange={setActiveTab}
                items={tabItems}
              />
            </Card>
          </>
        ) : (
          <div className="bg-white rounded-lg shadow-sm p-12">
            <Empty description="商品不存在" />
          </div>
        )}
      </Spin>
    </div>
  )
}

export default ProductDetail
