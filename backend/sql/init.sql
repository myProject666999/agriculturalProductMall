-- 农产品销售系统数据库初始化脚本
-- 创建数据库
CREATE DATABASE IF NOT EXISTS agricultural_mall DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE agricultural_mall;

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE COMMENT '用户名',
    password VARCHAR(255) NOT NULL COMMENT '密码',
    email VARCHAR(100) NOT NULL UNIQUE COMMENT '邮箱',
    phone VARCHAR(20) DEFAULT NULL COMMENT '手机号',
    avatar VARCHAR(255) DEFAULT NULL COMMENT '头像',
    nickname VARCHAR(50) DEFAULT NULL COMMENT '昵称',
    gender TINYINT DEFAULT 0 COMMENT '性别 0未知 1男 2女',
    role VARCHAR(20) DEFAULT 'user' COMMENT '角色 user用户 merchant商户 admin管理员',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    last_login_at DATETIME DEFAULT NULL COMMENT '最后登录时间',
    last_login_ip VARCHAR(50) DEFAULT NULL COMMENT '最后登录IP',
    email_verified_at DATETIME DEFAULT NULL COMMENT '邮箱验证时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_username (username),
    INDEX idx_email (email),
    INDEX idx_role (role)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 地址表
CREATE TABLE IF NOT EXISTS addresses (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    name VARCHAR(50) NOT NULL COMMENT '收货人姓名',
    phone VARCHAR(20) NOT NULL COMMENT '收货人电话',
    province VARCHAR(50) DEFAULT NULL COMMENT '省份',
    city VARCHAR(50) DEFAULT NULL COMMENT '城市',
    district VARCHAR(50) DEFAULT NULL COMMENT '区县',
    detail VARCHAR(255) NOT NULL COMMENT '详细地址',
    is_default TINYINT DEFAULT 0 COMMENT '是否默认 0否 1是',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收货地址表';

-- 商户表
CREATE TABLE IF NOT EXISTS merchants (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '用户ID',
    store_name VARCHAR(100) NOT NULL COMMENT '店铺名称',
    store_logo VARCHAR(255) DEFAULT NULL COMMENT '店铺Logo',
    store_desc TEXT COMMENT '店铺描述',
    contact_name VARCHAR(50) DEFAULT NULL COMMENT '联系人',
    contact_phone VARCHAR(20) DEFAULT NULL COMMENT '联系电话',
    contact_email VARCHAR(100) DEFAULT NULL COMMENT '联系邮箱',
    address VARCHAR(255) DEFAULT NULL COMMENT '店铺地址',
    business_license VARCHAR(255) DEFAULT NULL COMMENT '营业执照',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    rating DECIMAL(3,2) DEFAULT 5.00 COMMENT '评分',
    total_sales BIGINT DEFAULT 0 COMMENT '总销量',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商户表';

-- 分类表
CREATE TABLE IF NOT EXISTS categories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '分类名称',
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父分类ID',
    icon VARCHAR(255) DEFAULT NULL COMMENT '图标',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品分类表';

-- 商品表
CREATE TABLE IF NOT EXISTS products (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(200) NOT NULL COMMENT '商品名称',
    category_id BIGINT UNSIGNED NOT NULL COMMENT '分类ID',
    merchant_id BIGINT UNSIGNED NOT NULL COMMENT '商户ID',
    description TEXT COMMENT '商品描述',
    main_image VARCHAR(255) DEFAULT NULL COMMENT '主图',
    sub_images TEXT COMMENT '副图 JSON数组',
    price DECIMAL(10,2) NOT NULL COMMENT '销售价格',
    original_price DECIMAL(10,2) DEFAULT NULL COMMENT '原价',
    stock INT DEFAULT 0 COMMENT '库存',
    sales INT DEFAULT 0 COMMENT '销量',
    unit VARCHAR(20) DEFAULT '斤' COMMENT '单位',
    specifications TEXT COMMENT '规格 JSON',
    is_hot TINYINT DEFAULT 0 COMMENT '是否热销 0否 1是',
    is_new TINYINT DEFAULT 0 COMMENT '是否新品 0否 1是',
    is_recommend TINYINT DEFAULT 0 COMMENT '是否推荐 0否 1是',
    status TINYINT DEFAULT 1 COMMENT '状态 0下架 1上架',
    rating DECIMAL(3,2) DEFAULT 5.00 COMMENT '评分',
    review_count INT DEFAULT 0 COMMENT '评价数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_category_id (category_id),
    INDEX idx_merchant_id (merchant_id),
    INDEX idx_status (status),
    FULLTEXT INDEX idx_name (name),
    FOREIGN KEY (category_id) REFERENCES categories(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='商品表';

-- 库存表
CREATE TABLE IF NOT EXISTS inventories (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '商品ID',
    quantity INT DEFAULT 0 COMMENT '库存数量',
    min_stock INT DEFAULT 10 COMMENT '最小库存预警',
    max_stock INT DEFAULT 1000 COMMENT '最大库存',
    last_updated_at DATETIME DEFAULT CURRENT_TIMESTAMP COMMENT '最后更新时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存表';

-- 库存记录表
CREATE TABLE IF NOT EXISTS inventory_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    type VARCHAR(20) NOT NULL COMMENT '类型 in入库 out出库',
    quantity INT NOT NULL COMMENT '数量',
    before_qty INT NOT NULL COMMENT '变更前数量',
    after_qty INT NOT NULL COMMENT '变更后数量',
    operator_id BIGINT UNSIGNED NOT NULL COMMENT '操作人ID',
    remark VARCHAR(255) DEFAULT NULL COMMENT '备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_product_id (product_id),
    INDEX idx_operator_id (operator_id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (operator_id) REFERENCES users(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='库存记录表';

-- 轮播图表
CREATE TABLE IF NOT EXISTS banners (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '标题',
    image VARCHAR(255) NOT NULL COMMENT '图片URL',
    link_type VARCHAR(20) DEFAULT 'none' COMMENT '链接类型 none product category',
    link_value VARCHAR(255) DEFAULT NULL COMMENT '链接值',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='轮播图表';

-- 购物车表
CREATE TABLE IF NOT EXISTS carts (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    quantity INT DEFAULT 1 COMMENT '数量',
    selected TINYINT DEFAULT 1 COMMENT '是否选中 0否 1是',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE KEY idx_user_product (user_id, product_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='购物车表';

-- 订单表
CREATE TABLE IF NOT EXISTS orders (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_no VARCHAR(50) NOT NULL UNIQUE COMMENT '订单号',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    merchant_id BIGINT UNSIGNED NOT NULL COMMENT '商户ID',
    status VARCHAR(20) DEFAULT 'pending' COMMENT '状态 pending待付款 paid待发货 shipped待收货 completed已完成 cancelled已取消',
    total_amount DECIMAL(10,2) NOT NULL COMMENT '商品总金额',
    discount_amount DECIMAL(10,2) DEFAULT 0 COMMENT '优惠金额',
    pay_amount DECIMAL(10,2) NOT NULL COMMENT '实付金额',
    shipping_fee DECIMAL(10,2) DEFAULT 0 COMMENT '运费',
    receiver_name VARCHAR(50) NOT NULL COMMENT '收货人',
    receiver_phone VARCHAR(20) NOT NULL COMMENT '收货电话',
    receiver_address VARCHAR(255) NOT NULL COMMENT '收货地址',
    pay_method VARCHAR(20) DEFAULT NULL COMMENT '支付方式',
    pay_time DATETIME DEFAULT NULL COMMENT '支付时间',
    shipping_time DATETIME DEFAULT NULL COMMENT '发货时间',
    complete_time DATETIME DEFAULT NULL COMMENT '完成时间',
    refund_status VARCHAR(20) DEFAULT 'none' COMMENT '退款状态 none无 applying申请中 rejected已拒绝 completed已完成',
    refund_reason TEXT COMMENT '退款原因',
    refund_amount DECIMAL(10,2) DEFAULT 0 COMMENT '退款金额',
    remark VARCHAR(255) DEFAULT NULL COMMENT '订单备注',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_merchant_id (merchant_id),
    INDEX idx_order_no (order_no),
    INDEX idx_status (status),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单表';

-- 订单商品表
CREATE TABLE IF NOT EXISTS order_items (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    product_name VARCHAR(200) NOT NULL COMMENT '商品名称',
    product_image VARCHAR(255) DEFAULT NULL COMMENT '商品图片',
    price DECIMAL(10,2) NOT NULL COMMENT '商品单价',
    quantity INT NOT NULL COMMENT '购买数量',
    sub_total DECIMAL(10,2) NOT NULL COMMENT '小计',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_order_id (order_id),
    INDEX idx_product_id (product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='订单商品表';

-- 评价表
CREATE TABLE IF NOT EXISTS reviews (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    order_id BIGINT UNSIGNED NOT NULL COMMENT '订单ID',
    merchant_id BIGINT UNSIGNED NOT NULL COMMENT '商户ID',
    rating INT DEFAULT 5 COMMENT '评分 1-5',
    content TEXT NOT NULL COMMENT '评价内容',
    images TEXT COMMENT '图片 JSON数组',
    is_anonymous TINYINT DEFAULT 0 COMMENT '是否匿名 0否 1是',
    status TINYINT DEFAULT 1 COMMENT '状态 0隐藏 1显示',
    reply TEXT COMMENT '商家回复',
    reply_time DATETIME DEFAULT NULL COMMENT '回复时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    INDEX idx_order_id (order_id),
    INDEX idx_merchant_id (merchant_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (merchant_id) REFERENCES merchants(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='评价表';

-- 物流表
CREATE TABLE IF NOT EXISTS logistics (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    order_id BIGINT UNSIGNED NOT NULL UNIQUE COMMENT '订单ID',
    logistics_no VARCHAR(50) NOT NULL COMMENT '物流单号',
    company VARCHAR(50) NOT NULL COMMENT '物流公司',
    current_status VARCHAR(100) DEFAULT NULL COMMENT '当前状态',
    trajectory TEXT COMMENT '物流轨迹 JSON',
    ship_time DATETIME DEFAULT NULL COMMENT '发货时间',
    received_time DATETIME DEFAULT NULL COMMENT '收货时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_order_id (order_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='物流表';

-- 收藏表
CREATE TABLE IF NOT EXISTS favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    UNIQUE KEY idx_user_product (user_id, product_id),
    INDEX idx_user_id (user_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='收藏表';

-- 资讯表
CREATE TABLE IF NOT EXISTS news (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '标题',
    author VARCHAR(50) DEFAULT NULL COMMENT '作者',
    category VARCHAR(50) DEFAULT NULL COMMENT '分类',
    content LONGTEXT NOT NULL COMMENT '内容',
    cover VARCHAR(255) DEFAULT NULL COMMENT '封面图',
    description VARCHAR(500) DEFAULT NULL COMMENT '描述',
    views INT DEFAULT 0 COMMENT '浏览量',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    sort INT DEFAULT 0 COMMENT '排序',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='资讯表';

-- 公告表
CREATE TABLE IF NOT EXISTS announcements (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(200) NOT NULL COMMENT '标题',
    content TEXT NOT NULL COMMENT '内容',
    type VARCHAR(20) DEFAULT 'general' COMMENT '类型 general普通 notice通知 urgent紧急',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    start_at DATETIME DEFAULT NULL COMMENT '开始时间',
    end_at DATETIME DEFAULT NULL COMMENT '结束时间',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='公告表';

-- 菜单表
CREATE TABLE IF NOT EXISTS menus (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    parent_id BIGINT UNSIGNED DEFAULT 0 COMMENT '父菜单ID',
    name VARCHAR(50) NOT NULL COMMENT '菜单名称',
    path VARCHAR(255) DEFAULT NULL COMMENT '路由路径',
    icon VARCHAR(100) DEFAULT NULL COMMENT '图标',
    component VARCHAR(255) DEFAULT NULL COMMENT '组件路径',
    permission VARCHAR(100) DEFAULT NULL COMMENT '权限标识',
    type TINYINT DEFAULT 1 COMMENT '类型 1目录 2菜单 3按钮',
    sort INT DEFAULT 0 COMMENT '排序',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_parent_id (parent_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='菜单表';

-- 权限表
CREATE TABLE IF NOT EXISTS permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    menu_id BIGINT UNSIGNED DEFAULT NULL COMMENT '菜单ID',
    name VARCHAR(50) NOT NULL COMMENT '权限名称',
    code VARCHAR(100) NOT NULL COMMENT '权限代码',
    type VARCHAR(20) DEFAULT 'button' COMMENT '类型',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL,
    INDEX idx_menu_id (menu_id),
    FOREIGN KEY (menu_id) REFERENCES menus(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='权限表';

-- 角色表
CREATE TABLE IF NOT EXISTS roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL COMMENT '角色名称',
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '角色代码',
    description VARCHAR(255) DEFAULT NULL COMMENT '描述',
    status TINYINT DEFAULT 1 COMMENT '状态 0禁用 1启用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色表';

-- 角色权限关联表
CREATE TABLE IF NOT EXISTS role_permissions (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    permission_id BIGINT UNSIGNED NOT NULL COMMENT '权限ID',
    UNIQUE KEY idx_role_permission (role_id, permission_id),
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='角色权限关联表';

-- 用户角色关联表
CREATE TABLE IF NOT EXISTS user_roles (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role_id BIGINT UNSIGNED NOT NULL COMMENT '角色ID',
    UNIQUE KEY idx_user_role (user_id, role_id),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户角色关联表';

-- 验证码表
CREATE TABLE IF NOT EXISTS verification_codes (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(100) NOT NULL COMMENT '邮箱',
    code VARCHAR(10) NOT NULL COMMENT '验证码',
    type VARCHAR(20) NOT NULL COMMENT '类型 register reset_password',
    expires_at DATETIME NOT NULL COMMENT '过期时间',
    used TINYINT DEFAULT 0 COMMENT '是否已使用',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_email (email),
    INDEX idx_expires_at (expires_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='验证码表';

-- 用户行为记录表
CREATE TABLE IF NOT EXISTS user_behaviors (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    behavior VARCHAR(20) NOT NULL COMMENT '行为类型 view浏览 collect收藏 purchase购买 rating评分',
    score DECIMAL(5,2) NOT NULL COMMENT '行为权重分数',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户行为记录表';

-- 推荐记录表
CREATE TABLE IF NOT EXISTS recommend_records (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    product_id BIGINT UNSIGNED NOT NULL COMMENT '商品ID',
    score DECIMAL(5,2) NOT NULL COMMENT '推荐分数',
    recommend_type VARCHAR(20) NOT NULL COMMENT '推荐类型 user基于用户 item基于物品 hybrid混合',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id),
    INDEX idx_product_id (product_id),
    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (product_id) REFERENCES products(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='推荐记录表';

-- 插入初始数据
-- 插入角色数据
INSERT INTO roles (name, code, description, status) VALUES
('超级管理员', 'admin', '系统超级管理员', 1),
('商户', 'merchant', '入驻商户', 1),
('普通用户', 'user', '普通注册用户', 1);

-- 插入管理员账户 (密码: admin123)
INSERT INTO users (username, password, email, phone, nickname, role, status) VALUES
('admin', '$2a$10$rEq1GJQxJQxJQxJQxJQxJO7xJQxJQxJQxJQxJQxJQxJQxJQxJ', 'admin@example.com', '13800138000', '管理员', 'admin', 1);

-- 插入分类数据
INSERT INTO categories (name, parent_id, icon, sort, status) VALUES
('新鲜蔬菜', 0, 'vegetable', 1, 1),
('时令水果', 0, 'fruit', 2, 1),
('肉禽蛋品', 0, 'meat', 3, 1),
('粮油米面', 0, 'grain', 4, 1),
('水产海鲜', 0, 'seafood', 5, 1),
('干货特产', 0, 'dry', 6, 1),
('有机蔬菜', 1, 'organic', 1, 1),
('绿叶菜', 1, 'leaf', 2, 1),
('根茎菜', 1, 'root', 3, 1),
('柑橘类', 2, 'citrus', 1, 1),
('苹果类', 2, 'apple', 2, 1),
('热带水果', 2, 'tropical', 3, 1);

-- 插入轮播图数据
INSERT INTO banners (title, image, link_type, link_value, sort, status) VALUES
('新鲜直达，品质保证', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20vegetables%20and%20fruits%20agriculture%20market%20banner&image_size=landscape_16_9', 'none', '', 1, 1),
('限时特惠，每日低价', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=special%20discount%20agriculture%20products%20promotion%20banner&image_size=landscape_16_9', 'none', '', 2, 1),
('新用户专享福利', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=new%20user%20welcome%20gift%20agriculture%20shopping%20banner&image_size=landscape_16_9', 'none', '', 3, 1);

-- 插入资讯数据
INSERT INTO news (title, author, category, content, cover, description, views, status, sort) VALUES
('春季蔬菜种植技术要点', '农业专家', '种植技术', '春季是蔬菜种植的黄金时期，以下是几个关键技术要点...', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=spring%20vegetable%20planting%20guide%20agriculture%20article&image_size=landscape_4_3', '春季蔬菜种植的关键技术指南，助您获得丰收', 1256, 1, 1),
('如何辨别新鲜农产品', '品质检测', '消费指南', '购买新鲜农产品时，需要注意以下几个方面...', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=fresh%20agricultural%20products%20quality%20inspection&image_size=landscape_4_3', '教您如何辨别新鲜优质的农产品', 892, 1, 2),
('农业数字化转型新趋势', '行业分析师', '行业动态', '随着科技的发展，农业数字化转型已成必然趋势...', 'https://trae-api-cn.mchost.guru/api/ide/v1/text_to_image?prompt=digital%20agriculture%20technology%20transformation%20trend&image_size=landscape_4_3', '探讨农业数字化转型的最新发展趋势', 654, 1, 3);

-- 插入公告数据
INSERT INTO announcements (title, content, type, status) VALUES
('系统升级通知', '为了提供更好的服务，系统将于本周日凌晨2:00-4:00进行例行维护升级。', 'notice', 1),
('新用户注册送优惠券', '欢迎新用户注册！注册即送100元优惠券大礼包，赶快行动吧！', 'general', 1);
