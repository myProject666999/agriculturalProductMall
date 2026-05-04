import React, { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import { Card, Row, Col, Input, Select, Pagination, Spin, Empty, Button, message, Tag } from 'antd'
import { SearchOutlined, ShoppingCartOutlined, HeartOutlined } from '@ant-design/icons'
import { productApi } from '@/services/api'
import { useAppDispatch } from '@/store'
import { addItem } from '@/store/cartSlice'

const { Meta } = Card
const { Search } = Input
const { Option } = Select

const ProductList: React.FC = () => {
  const navigate = useNavigate()
  const dispatch = useAppDispatch()
  const [searchParams, setSearchParams] = useSearchParams()
  const [products, setProducts] = useState<any[]>([])
  const [categories, setCategories] = useState<any[]>([])
  const [loading, setLoading] = useState(false)
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 12,
    total: 0,
  })
  const [filters, setFilters] = useState({
    keyword: searchParams.get('keyword') || '',
    categoryId: searchParams.get('categoryId') ? Number(searchParams.get('categoryId')) : undefined,
    sort: 'default',
  })

  useEffect(() => {
    loadCategories()
  }, [])

  useEffect(() => {
    loadProducts()
  }, [pagination.current, pagination.pageSize, filters])

  const loadCategories = async () => {
    try {
      const res: any = await productApi.getCategories()
      if (res.code === 0) {
        setCategories(res.data || [])
      }
    } catch (error) {
      console.error('Failed to load categories:', error)
    }
  }

  const loadProducts = async () => {
    setLoading(true)
    try {
      const params: any = {
        page: pagination.current,
        pageSize: pagination.pageSize,
      }
      if (filters.keyword) {
        params.keyword = filters.keyword
      }
      if (filters.categoryId) {
        params.categoryId = filters.categoryId
      }
      if (filters.sort) {
        params.sort = filters.sort
      }

      const res: any = await productApi.getList(params)
      if (res.code === 0) {
        const data = res.data
        setProducts(data.list || data || [])
        setPagination((prev) => ({
          ...prev,
          total: data.total || (data.list ? data.list.length : 0),
        }))
      }
    } catch (error) {
      console.error('Failed to load products:', error)
      setProducts([])
    } finally {
      setLoading(false)
    }
  }

  const handleSearch = (value: string) => {
    setFilters((prev) => ({ ...prev, keyword: value }))
    setPagination((prev) => ({ ...prev, current: 1 }))
    if (value) {
      setSearchParams({ keyword: value })
    } else {
      setSearchParams({})
    }
  }

  const handleCategoryChange = (value: number | undefined) => {
    setFilters((prev) => ({ ...prev, categoryId: value }))
    setPagination((prev) => ({ ...prev, current: 1 }))
  }

  const handleSortChange = (value: string) => {
    setFilters((prev) => ({ ...prev, sort: value }))
    setPagination((prev) => ({ ...prev, current: 1 }))
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

  const handlePageChange = (page: number, pageSize: number) => {
    setPagination((prev) => ({
      ...prev,
      current: page,
      pageSize: pageSize || prev.pageSize,
    }))
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="bg-white rounded-lg shadow-sm p-6 mb-6">
        <div className="flex flex-wrap items-center gap-4">
          <Search
            placeholder="搜索商品..."
            allowClear
            enterButton="搜索"
            size="large"
            defaultValue={filters.keyword}
            onSearch={handleSearch}
            style={{ width: 300 }}
            prefix={<SearchOutlined />}
          />
          <Select
            placeholder="全部分类"
            allowClear
            size="large"
            style={{ width: 180 }}
            value={filters.categoryId}
            onChange={handleCategoryChange}
          >
            {categories.map((cat) => (
              <Option key={cat.id} value={cat.id}>
                {cat.name}
              </Option>
            ))}
          </Select>
          <Select
            placeholder="排序方式"
            size="large"
            style={{ width: 150 }}
            value={filters.sort}
            onChange={handleSortChange}
          >
            <Option value="default">默认排序</Option>
            <Option value="price-asc">价格从低到高</Option>
            <Option value="price-desc">价格从高到低</Option>
            <Option value="sales">销量优先</Option>
            <Option value="new">最新上架</Option>
          </Select>
        </div>
      </div>

      <Spin spinning={loading}>
        {products.length > 0 ? (
          <>
            <Row gutter={[16, 16]}>
              {products.map((product) => (
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
                          {product.originalPrice && (
                            <span className="text-gray-400 text-sm line-through ml-2">
                              ¥{product.originalPrice}
                            </span>
                          )}
                        </div>
                      }
                    />
                    <div className="mt-2">
                      {product.isNew && <Tag color="green">新品</Tag>}
                      {product.isHot && <Tag color="red">热卖</Tag>}
                      <span className="text-gray-400 text-sm ml-2">
                        销量: {product.sales || 0}
                      </span>
                    </div>
                  </Card>
                </Col>
              ))}
            </Row>
            <div className="text-center mt-8">
              <Pagination
                current={pagination.current}
                pageSize={pagination.pageSize}
                total={pagination.total || products.length}
                onChange={handlePageChange}
                showSizeChanger
                showQuickJumper
                showTotal={(total) => `共 ${total} 件商品`}
              />
            </div>
          </>
        ) : (
          <div className="bg-white rounded-lg shadow-sm p-12">
            <Empty description="暂无商品" />
          </div>
        )}
      </Spin>
    </div>
  )
}

export default ProductList
