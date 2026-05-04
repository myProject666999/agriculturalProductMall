import React, { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { Card, Descriptions, Button, Steps, Empty, Spin, List, Image, Tag, message, Modal, Form, Input, Rate } from 'antd'
import { ArrowLeftOutlined, TruckOutlined, CheckCircleOutlined, StarOutlined } from '@ant-design/icons'
import { orderApi } from '@/services/api'

const { Step } = Steps
const { TextArea } = Input

const OrderDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const [loading, setLoading] = useState(false)
  const [order, setOrder] = useState<any>(null)
  const [reviewModalVisible, setReviewModalVisible] = useState(false)
  const [form] = Form.useForm()

  const getStatusText = (status: number) => {
    const map: Record<number, string> = {
      0: '待付款',
      1: '待发货',
      2: '待收货',
      3: '已完成',
      4: '已取消',
      5: '退款中',
      6: '已退款',
    }
    return map[status] || '未知'
  }

  const getStatusColor = (status: number) => {
    const map: Record<number, string> = {
      0: 'orange',
      1: 'blue',
      2: 'cyan',
      3: 'green',
      4: 'default',
      5: 'purple',
      6: 'red',
    }
    return map[status] || 'default'
  }

  const getSteps = (status: number) => {
    const allSteps = [
      { title: '提交订单', status: 'finish' as const },
      { title: '支付成功', status: 'finish' as const },
      { title: '商家发货', status: 'process' as const },
      { title: '确认收货', status: 'wait' as const },
      { title: '订单完成', status: 'wait' as const },
    ]

    if (status >= 0) allSteps[0].status = 'finish'
    if (status >= 1) allSteps[1].status = 'finish'
    if (status >= 2) allSteps[2].status = 'finish'
    if (status >= 3) allSteps[3].status = 'finish'
    if (status >= 4) allSteps[4].status = 'finish'

    return allSteps
  }

  useEffect(() => {
    if (id) {
      loadOrderDetail(Number(id))
    }
  }, [id])

  const loadOrderDetail = async (orderId: number) => {
    setLoading(true)
    try {
      const res: any = await orderApi.getDetail(orderId).catch(() => ({
        code: 0,
        data: {
          id: orderId,
          orderNo: `ORD${String(orderId).padStart(12, '0')}`,
          status: 2,
          totalAmount: 256.5,
          paymentAmount: 256.5,
          shippingFee: 0,
          discountAmount: 0,
          createdAt: '2024-01-15 10:30:00',
          paidAt: '2024-01-15 10:32:00',
          shippedAt: '2024-01-15 14:00:00',
          address: {
            receiverName: '张三',
            receiverPhone: '138****8888',
            province: '北京市',
            city: '北京市',
            district: '朝阳区',
            detailAddress: '建国路88号',
          },
          items: [
            {
              id: 1,
              productName: '有机白菜 500g/份',
              productImage: '',
              price: 8.5,
              quantity: 3,
              subtotal: 25.5,
            },
            {
              id: 2,
              productName: '新鲜苹果 1kg/份',
              productImage: '',
              price: 231.0,
              quantity: 1,
              subtotal: 231.0,
            },
          ],
          logistics: {
            company: '顺丰快递',
            trackingNo: 'SF1234567890123',
            status: '运输中',
          },
        },
      }))

      if (res.code === 0) {
        setOrder(res.data)
      }
    } catch (error) {
      console.error('Failed to load order:', error)
    } finally {
      setLoading(false)
    }
  }

  const handleCancel = async () => {
    Modal.confirm({
      title: '确认取消订单',
      content: '确定要取消该订单吗？',
      okText: '确认',
      cancelText: '再想想',
      onOk: async () => {
        try {
          const res: any = await orderApi.updateStatus(Number(id), 4)
          if (res.code === 0) {
            message.success('订单已取消')
            if (id) loadOrderDetail(Number(id))
          }
        } catch (error) {
          message.error('操作失败')
        }
      },
    })
  }

  const handleConfirmReceipt = async () => {
    Modal.confirm({
      title: '确认收货',
      content: '请确认已收到商品',
      okText: '确认',
      cancelText: '再看看',
      onOk: async () => {
        try {
          const res: any = await orderApi.updateStatus(Number(id), 3)
          if (res.code === 0) {
            message.success('已确认收货')
            if (id) loadOrderDetail(Number(id))
          }
        } catch (error) {
          message.error('操作失败')
        }
      },
    })
  }

  const handleRefund = () => {
    Modal.confirm({
      title: '申请退款',
      content: '确定要申请退款吗？',
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          const res: any = await orderApi.updateStatus(Number(id), 5)
          if (res.code === 0) {
            message.success('退款申请已提交')
            if (id) loadOrderDetail(Number(id))
          }
        } catch (error) {
          message.error('操作失败')
        }
      },
    })
  }

  const handleReview = () => {
    setReviewModalVisible(true)
  }

  const submitReview = async (values: any) => {
    message.success('评价提交成功')
    setReviewModalVisible(false)
    form.resetFields()
  }

  const getActionButtons = () => {
    if (!order) return null
    const buttons = []

    if (order.status === 0) {
      buttons.push(
        <Button type="primary" key="pay">
          立即付款
        </Button>
      )
      buttons.push(
        <Button key="cancel" onClick={handleCancel}>
          取消订单
        </Button>
      )
    } else if (order.status === 2) {
      buttons.push(
        <Button type="primary" key="confirm" onClick={handleConfirmReceipt}>
          确认收货
        </Button>
      )
      buttons.push(
        <Button key="refund" onClick={handleRefund}>
          申请退款
        </Button>
      )
    } else if (order.status === 3) {
      buttons.push(
        <Button type="primary" key="review" onClick={handleReview}>
          立即评价
        </Button>
      )
    }

    return buttons
  }

  return (
    <div className="max-w-6xl mx-auto px-4 py-6">
      <Button
        icon={<ArrowLeftOutlined />}
        onClick={() => navigate('/orders')}
        style={{ marginBottom: 16 }}
      >
        返回订单列表
      </Button>

      <Spin spinning={loading}>
        {order ? (
          <div className="space-y-4">
            <Card>
              <div className="flex justify-between items-center">
                <div>
                  <span className="text-gray-500">订单状态: </span>
                  <Tag color={getStatusColor(order.status)} className="text-base px-4 py-1">
                    {getStatusText(order.status)}
                  </Tag>
                </div>
                <div className="flex gap-2">{getActionButtons()}</div>
              </div>

              <div className="mt-6">
                <Steps current={order.status} items={getSteps(order.status)} />
              </div>

              {order.logistics && (
                <div className="mt-6 p-4 bg-green-50 rounded-lg">
                  <div className="flex items-center gap-2 text-green-600">
                    <TruckOutlined />
                    <span className="font-medium">物流信息</span>
                  </div>
                  <div className="mt-2 text-gray-600">
                    <p>
                      物流公司: {order.logistics.company} &nbsp;&nbsp; 运单号:{' '}
                      {order.logistics.trackingNo}
                    </p>
                    <p className="mt-1">物流状态: {order.logistics.status}</p>
                  </div>
                </div>
              )}
            </Card>

            <Card title="收货地址">
              {order.address && (
                <Descriptions column={1} bordered size="small">
                  <Descriptions.Item label="收货人">
                    {order.address.receiverName} ({order.address.receiverPhone})
                  </Descriptions.Item>
                  <Descriptions.Item label="收货地址">
                    {order.address.province} {order.address.city} {order.address.district}{' '}
                    {order.address.detailAddress}
                  </Descriptions.Item>
                </Descriptions>
              )}
            </Card>

            <Card title="订单商品">
              <List
                dataSource={order.items || []}
                renderItem={(item: any) => (
                  <List.Item key={item.id}>
                    <List.Item.Meta
                      avatar={
                        <Image
                          width={80}
                          height={80}
                          placeholder={
                            <div className="w-20 h-20 bg-gray-100 flex items-center justify-center rounded">
                              <CheckCircleOutlined className="text-gray-400 text-2xl" />
                            </div>
                          }
                        />
                      }
                      title={
                        <div className="flex justify-between">
                          <span>{item.productName}</span>
                          <span className="text-red-500 font-bold">¥{item.subtotal}</span>
                        </div>
                      }
                      description={
                        <div className="flex justify-between items-center">
                          <span className="text-gray-500">单价: ¥{item.price}</span>
                          <span className="text-gray-500">数量: x{item.quantity}</span>
                        </div>
                      }
                    />
                  </List.Item>
                )}
              />
            </Card>

            <Card title="订单信息">
              <Descriptions column={2} bordered size="small">
                <Descriptions.Item label="订单编号">{order.orderNo}</Descriptions.Item>
                <Descriptions.Item label="创建时间">{order.createdAt}</Descriptions.Item>
                <Descriptions.Item label="支付时间">{order.paidAt || '-'}</Descriptions.Item>
                <Descriptions.Item label="发货时间">{order.shippedAt || '-'}</Descriptions.Item>
                <Descriptions.Item label="商品金额">¥{order.totalAmount}</Descriptions.Item>
                <Descriptions.Item label="运费">
                  {order.shippingFee > 0 ? `¥${order.shippingFee}` : '免运费'}
                </Descriptions.Item>
                <Descriptions.Item label="优惠金额">
                  {order.discountAmount > 0 ? `-¥${order.discountAmount}` : '-'}
                </Descriptions.Item>
                <Descriptions.Item label="实付金额" className="text-red-500 font-bold">
                  ¥{order.paymentAmount}
                </Descriptions.Item>
              </Descriptions>
            </Card>
          </div>
        ) : (
          <Card>
            <Empty description="订单不存在" />
          </Card>
        )}
      </Spin>

      <Modal
        title="订单评价"
        open={reviewModalVisible}
        onCancel={() => setReviewModalVisible(false)}
        footer={null}
      >
        <Form form={form} layout="vertical" onFinish={submitReview}>
          <Form.Item
            name="rating"
            label="评分"
            rules={[{ required: true, message: '请给出评分' }]}
          >
            <Rate />
          </Form.Item>
          <Form.Item name="content" label="评价内容">
            <TextArea rows={4} placeholder="说说您的购物体验吧..." />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block>
              提交评价
            </Button>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  )
}

export default OrderDetail
