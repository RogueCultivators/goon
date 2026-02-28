# Blog 示例项目

这是一个使用 Goon 生成的完整博客系统示例，展示了如何使用 Goon 快速构建一个功能完整的 Web API。

## 功能特性

- 用户管理（注册、登录、个人资料）
- 文章管理（CRUD、分页、搜索）
- 评论系统
- 标签分类
- JWT 认证
- PostgreSQL 数据库
- Docker 支持

## 快速开始

### 1. 生成项目

```bash
# 使用 goon 生成项目
goon init blog --example
cd blog
```

### 2. 启动数据库

```bash
docker-compose up -d
```

### 3. 运行迁移

```bash
make migrate-up
```

### 4. 填充测试数据

```bash
make seed
```

### 5. 启动服务

```bash
go run main.go
```

服务将在 `http://localhost:8080` 启动。

## API 端点

### 用户相关

- `POST /api/v1/auth/register` - 用户注册
- `POST /api/v1/auth/login` - 用户登录
- `GET /api/v1/user/:id` - 获取用户信息
- `PUT /api/v1/user/:id` - 更新用户信息

### 文章相关

- `GET /api/v1/post` - 获取文章列表（支持分页、搜索）
- `GET /api/v1/post/:id` - 获取文章详情
- `POST /api/v1/post` - 创建文章（需要认证）
- `PUT /api/v1/post/:id` - 更新文章（需要认证）
- `DELETE /api/v1/post/:id` - 删除文章（需要认证）

### 评论相关

- `GET /api/v1/post/:id/comments` - 获取文章评论
- `POST /api/v1/post/:id/comments` - 添加评论（需要认证）
- `DELETE /api/v1/comment/:id` - 删除评论（需要认证）

## 测试 API

### 注册用户

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "张三",
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

### 登录

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "zhangsan@example.com",
    "password": "password123"
  }'
```

### 创建文章

```bash
curl -X POST http://localhost:8080/api/v1/post \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "title": "我的第一篇博客",
    "content": "这是文章内容...",
    "tags": ["技术", "Go"]
  }'
```

### 获取文章列表

```bash
# 基本查询
curl http://localhost:8080/api/v1/post

# 分页查询
curl "http://localhost:8080/api/v1/post?page=1&page_size=10"

# 搜索文章
curl "http://localhost:8080/api/v1/post?keyword=Go"

# 按标签过滤
curl "http://localhost:8080/api/v1/post?tag=技术"
```

## 项目结构

```
blog/
├── cmd/
│   └── server/
│       └── server.go          # 服务器启动
├── internal/
│   ├── user/                  # 用户模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── schema.go
│   ├── post/                  # 文章模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── schema.go
│   ├── comment/               # 评论模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   ├── model.go
│   │   ├── repository.go
│   │   └── schema.go
│   ├── auth/                  # 认证模块
│   │   ├── handler.go
│   │   ├── service.go
│   │   └── schema.go
│   ├── config/
│   │   └── config.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   └── router/
│       └── router.go
├── pkg/
│   ├── response/
│   ├── logger/
│   ├── database/
│   └── jwt/
├── migrations/                # 数据库迁移
│   ├── 000001_create_users_table.up.sql
│   ├── 000001_create_users_table.down.sql
│   ├── 000002_create_posts_table.up.sql
│   ├── 000002_create_posts_table.down.sql
│   ├── 000003_create_comments_table.up.sql
│   └── 000003_create_comments_table.down.sql
├── scripts/
│   ├── setup.sh               # 初始化脚本
│   └── seed.sh                # 数据填充
├── docs/
│   └── api.md                 # API 文档
├── docker-compose.yaml
├── Makefile
├── .env.example
└── main.go
```

## 数据库设计

### users 表

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    avatar VARCHAR(500),
    bio TEXT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### posts 表

```sql
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id),
    title VARCHAR(500) NOT NULL,
    content TEXT NOT NULL,
    summary VARCHAR(1000),
    cover_image VARCHAR(500),
    tags TEXT[],
    view_count INT DEFAULT 0,
    status VARCHAR(20) DEFAULT 'published',
    published_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

### comments 表

```sql
CREATE TABLE comments (
    id BIGSERIAL PRIMARY KEY,
    post_id BIGINT NOT NULL REFERENCES posts(id),
    user_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    parent_id BIGINT REFERENCES comments(id),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
```

## 开发命令

```bash
# 构建项目
make build

# 运行项目
make run

# 运行测试
make test

# 代码格式化
make fmt

# 代码检查
make lint

# 启动数据库
make docker-up

# 停止数据库
make docker-down

# 执行迁移
make migrate-up

# 回滚迁移
make migrate-down

# 填充测试数据
make seed
```

## 扩展功能

这个示例可以轻松扩展以下功能：

1. **文件上传** - 添加图片上传功能
2. **点赞收藏** - 添加文章点赞和收藏
3. **关注系统** - 用户之间的关注关系
4. **全文搜索** - 集成 Elasticsearch
5. **缓存** - 使用 Redis 缓存热门文章
6. **消息队列** - 异步处理通知
7. **WebSocket** - 实时评论通知

## 学习要点

通过这个示例，你可以学习到：

1. 如何使用 Goon 快速生成项目结构
2. 模块化的代码组织方式
3. RESTful API 设计最佳实践
4. JWT 认证实现
5. 数据库迁移管理
6. Docker 容器化部署
7. 分页、搜索、过滤的实现
8. 关联查询和数据关系处理

## 许可证

MIT
