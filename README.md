# Course-Selection-System

字节跳动 ByteCamp 2021 后端训练营项目 - 选排课系统

[![Go Version](https://img.shields.io/badge/Go-1.21-blue)](https://go.dev/)
[![Gin](https://img.shields.io/badge/Gin-1.9-green)](https://gin-gonic.com/)
[![GORM](https://img.shields.io/badge/GORM-2.0-orange)](https://gorm.io/)

## 项目简介

本项目是一个完整的选排课系统，面向字节跳动 ByteCamp 2021 后端训练营大作业。系统实现了登录认证、成员管理、课程管理、教师绑定、排课算法和学生选课等核心功能。

在压测评估中取得 Top 1 的成绩，成功应对 60 万并发请求。

## 核心特性

- **高并发选课**: Redis 原子操作 + RocketMQ 异步处理 + 令牌桶限流
- **排课算法**: 基于二分图最大匹配的自动排课系统
- **优雅架构**: DDD 四层架构设计，清晰的职责分离
- **现代技术栈**: Go 1.21 + Gin + GORM v2 + Redis + RocketMQ
- **完整监控**: Prometheus 指标 + 结构化日志
- **容器化部署**: Docker + docker-compose 一键部署

## 技术栈

| 层次 | 技术 | 用途 |
|------|------|------|
| 语言 | Go 1.21 | 主力开发语言 |
| Web 框架 | Gin 1.9 | HTTP 路由和中间件 |
| ORM | GORM v2 | MySQL 数据库操作 |
| 缓存 | Redis (redigo) | 高并发场景缓存 |
| 消息队列 | RocketMQ | 异步写入数据库 |
| 配置 | Viper | 配置管理 |
| 日志 | Zap | 结构化日志 |
| 监控 | Prometheus | 指标监控 |
| 部署 | Docker | 容器化部署 |
| 验证 | go-playground/validator | 参数校验 |

## DDD 架构设计

本项目采用 **领域驱动设计 (DDD)** 的四层架构，实现了清晰的职责分离和高度可维护性。

### 架构分层

```
┌─────────────────────────────────────────────────────────────┐
│                     Interface Layer (接口层)                  │
│  负责与外部系统交互，处理 HTTP 请求、路由、中间件              │
│  ┌──────────┐  ┌────────────┐  ┌────────┐  ┌────────────┐  │
│  │ Handlers │  │ Middleware │  │ Router │  │  Consumer  │  │
│  └────┬─────┘  └─────┬──────┘  └────┬───┘  └──────┬─────┘  │
└───────┼──────────────┼─────────────┼────────────┼─────────┘
        │              │             │            │
        v              v             v            v
┌─────────────────────────────────────────────────────────────┐
│                  Application Layer (应用层)                   │
│  协调领域对象执行用例，处理事务、限流、日志记录                 │
│  ┌────────────────┐  ┌──────────────────────────────────┐  │
│  │  Member App    │  │  Course App (选课高并发逻辑)     │  │
│  │  Service       │  │  - 令牌桶限流                    │  │
│  └───────┬────────┘  │  - Redis 原子操作                │  │
│          │           │  - RocketMQ 异步写入             │  │
│          v           └──────────────────────────────────┘  │
│  ┌────────────────┐                                         │
│  │  DTOs          │                                         │
│  └───────────────┘                                         │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            v
┌─────────────────────────────────────────────────────────────┐
│                    Domain Layer (领域层)                      │
│  核心业务逻辑所在，包含实体、值对象、领域服务                   │
│  ┌─────────────────────────────────────────────────────┐   │
│  │  Domain Services                                    │   │
│  │  - AuthService (登录/登出/会话)                     │   │
│  │  - MemberService (成员CRUD)                        │   │
│  │  - CourseService (课程管理)                        │   │
│  │  - SchedulingService (二分图排课算法)              │   │
│  └─────────────────────────────────────────────────────┘   │
│  ┌─────────────┐  ┌─────────────────────────────────────┐  │
│  │  Entities   │  │  Repository Interfaces (依赖倒置)   │  │
│  │  Member     │  │  IMemberRepo                        │  │
│  │  Course     │  │  ICourseRepo                        │  │
│  │  Bind       │  │  IChoiceRepo                        │  │
│  │  Choice     │  │  IScheduleRepo                      │  │
│  └─────────────┘  └─────────────────────────────────────┘  │
└───────────────────────────┬─────────────────────────────────┘
                            │ 实现接口
                            v
┌─────────────────────────────────────────────────────────────┐
│               Infrastructure Layer (基础设施层)                │
│  提供技术实现：数据库、缓存、消息队列、外部服务                 │
│  ┌──────────┐  ┌───────────┐  ┌────────┐  ┌──────────────┐  │
│  │ GORM v2  │  │ Redis     │  │ Rocket │  │  Encryption  │  │
│  │ MySQL    │  │ (redigo)  │  │  MQ    │  │  (MD5/Bcrypt)│  │
│  └──────────┘  └───────────┘  └────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### 各层职责

| 层级 | 职责 | 包含内容 |
|------|------|----------|
| **Interface** | 协议转换、请求路由、异常处理 | Handler、Middleware、Router、Consumer |
| **Application** | 用例编排、事务管理、限流、日志 | AppService、DTO、Assembler |
| **Domain** | 业务规则、领域模型、领域服务 | Entity、ValueObject、DomainService、Repository Interface |
| **Infrastructure** | 技术实现、外部集成 | GORM、Redis、RocketMQ、Encryption |

### 依赖关系

```
┌─────────────────────────────────────────────────────────┐
│                    依赖方向                              │
└─────────────────────────────────────────────────────────┘

Interface Layer    ─────依赖─────►    Application Layer
        │                              │
        │                              │
        ▼                              ▼
Infrastructure ◄────────依赖──────── Domain Layer
        │                              │
        │                              │
        └────────实现接口──────────────┘
             (Dependency Inversion)
```

- **上层依赖下层的抽象接口**，而非具体实现
- **基础设施层实现仓储接口**，通过依赖注入注入到领域层
- **依赖倒置原则**：高层模块不依赖低层模块，都依赖于抽象
- 这种设计使得各层可以独立测试和替换实现

### 领域模型

#### 核心实体 (Entities)

| 实体 | 说明 | 主要属性 |
|------|------|----------|
| Member | 成员 | ID、Username、Password、Type、Status、CreateTime |
| Course | 课程 | ID、Name、Code、Credit、Capacity、TeacherID |
| Choice | 选课记录 | ID、StudentID、CourseID、Status、CreateTime |
| Schedule | 排课记录 | ID、TeacherID、CourseID、TimeSlot、WeekDay |
| Bind | 教师课程绑定 | ID、TeacherID、CourseID、CreateTime |

#### 仓储接口 (Repository Interfaces)

```go
// 仓储接口定义在 domain/repository/ 目录
type IMemberRepo interface {
    Create(ctx context.Context, member *Member) error
    GetByID(ctx context.Context, id uint) (*Member, error)
    GetByUsername(ctx context.Context, username string) (*Member, error)
    Update(ctx context.Context, member *Member) error
    Delete(ctx context.Context, id uint) error
}

// 基础设施层实现这些接口
type MemberRepo struct {
    db *gorm.DB
}

func (r *MemberRepo) Create(ctx context.Context, member *Member) error {
    return r.db.WithContext(ctx).Create(member).Error
}
```

## 快速开始

### 环境要求

- Go 1.21+
- MySQL 8.0+
- Redis 7.0+
- RocketMQ 4.9+

### 本地运行

```bash
# 1. 克隆项目
git clone https://github.com/pearFL/Course-Selection-System.git
cd Course-Selection-System

# 2. 配置环境变量或修改 config.yaml
cp config.yaml config.yaml.example
# 编辑 config.yaml 设置数据库、Redis、RocketMQ 连接信息

# 3. 下载依赖
go mod tidy

# 4. 运行服务
go run ./cmd/server/main.go
```

### Docker 部署

```bash
cd deploy

# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f app
```

RocketMQ 控制台访问: http://localhost:8081

### 配置文件说明

`config.yaml` 是项目的核心配置文件，支持环境变量覆盖：

```yaml
# 应用配置
app:
  name: "course-selection-system"
  host: "0.0.0.0"   # 监听地址
  port: 8080         # 监听端口
  env: "development" # 运行环境

# 数据库配置
database:
  host: "localhost"
  port: 3306
  username: "root"
  password: "${DB_PASSWORD}"  # 环境变量支持
  name: "course_select"

# Redis 配置
redis:
  host: "localhost"
  port: 6379
  pool_size: 100

# RocketMQ 配置
rocketmq:
  nameserver: "localhost:9876"
  group_id: "course-select-group"
  topic: "course-booking-topic"

# 限流配置
rate_limit:
  qps: 4000
  burst: 5000
```

## API 接口

### 认证模块

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/auth/login` | 用户登录 | 公开 |
| POST | `/api/v1/auth/logout` | 用户登出 | 需登录 |
| GET | `/api/v1/auth/whoami` | 获取当前用户 | 需登录 |

### 成员管理模块

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/member/create` | 创建成员 | 管理员 |
| GET | `/api/v1/member` | 获取单个成员 | 需登录 |
| GET | `/api/v1/member/list` | 获取成员列表 | 需登录 |
| POST | `/api/v1/member/update` | 更新成员 | 管理员 |
| POST | `/api/v1/member/delete` | 删除成员 | 管理员 |

### 课程管理模块

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/course/create` | 创建课程 | 管理员 |
| GET | `/api/v1/course/get` | 获取课程信息 | 需登录 |
| POST | `/api/v1/teacher/bind_course` | 绑定课程到教师 | 管理员 |
| POST | `/api/v1/teacher/unbind_course` | 解绑课程 | 管理员 |
| GET | `/api/v1/teacher/get_course` | 获取教师课程列表 | 需登录 |
| POST | `/api/v1/course/schedule` | 自动排课 | 管理员 |

### 选课模块

| 方法 | 路径 | 说明 | 权限 |
|------|------|------|------|
| POST | `/api/v1/student/book_course` | 学生选课 | 需登录 |
| GET | `/api/v1/student/course` | 获取学生课表 | 需登录 |

### 监控端点

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/health` | 健康检查 |
| GET | `/metrics` | Prometheus 指标 |

## 请求响应格式

### 成功响应
```json
{
    "code": 0,
    "message": "success",
    "data": {
        // 业务数据
    }
}
```

### 错误响应
```json
{
    "code": 1,
    "message": "参数不合法"
}
```

### 错误码说明

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1 | 参数不合法 |
| 2 | 用户名已存在 |
| 3 | 用户已删除 |
| 4 | 用户不存在 |
| 5 | 密码错误 |
| 6 | 用户未登录 |
| 7 | 课程已满 |
| 8 | 课程已绑定 |
| 9 | 课程未绑定 |
| 10 | 没有操作权限 |
| 11 | 学生不存在 |
| 12 | 课程不存在 |
| 13 | 学生没有课程 |
| 14 | 学生已有课程 |
| 15 | 重复请求 |
| 255 | 未知错误 |

## 用户类型

| 类型值 | 说明 |
|--------|------|
| 1 | 管理员 (Admin) |
| 2 | 学生 (Student) |
| 3 | 教师 (Teacher) |

**系统内置管理员**: `JudgeAdmin` / `JudgePassword2022`

## 项目结构

```
Course-Selection-System/
├── cmd/server/
│   └── main.go                 # 应用入口
├── internal/                   # DDD 四层架构
│   ├── config/                 # 配置层
│   │   └── config.go           # Viper 配置加载
│   ├── domain/                 # 领域层 (DDD 核心)
│   │   ├── model/              # 领域实体 (Member, Course, Choice, Schedule)
│   │   ├── repository/         # 仓储接口 (依赖倒置)
│   │   │   ├── imember.go
│   │   │   ├── icourse.go
│   │   │   ├── ichoice.go
│   │   │   └── ischedule.go
│   │   └── service/            # 领域服务
│   │       ├── auth.go
│   │       ├── member.go
│   │       ├── course.go
│   │       └── scheduling.go
│   ├── application/            # 应用层
│   │   ├── dto/                # 数据传输对象
│   │   │   ├── member_dto.go
│   │   │   ├── course_dto.go
│   │   │   └── booking_dto.go
│   │   └── service/            # 应用服务
│   │       ├── member_app.go
│   │       └── course_app.go
│   ├── infrastructure/         # 基础设施层
│   │   ├── database/           # GORM v2 实现
│   │   │   ├── gorm.go
│   │   │   ├── member_repo.go
│   │   │   ├── course_repo.go
│   │   │   └── choice_repo.go
│   │   ├── redis/              # Redis 客户端
│   │   │   └── redis.go
│   │   ├── mq/                 # RocketMQ
│   │   │   └── rocketmq.go
│   │   └── metrics/            # Prometheus 监控
│   │       └── metrics.go
│   ├── interface/              # 接口层
│   │   └── api/                # HTTP API
│   │       ├── handler/        # HTTP Handlers
│   │       ├── middleware/     # 中间件 (Auth, Session, RateLimit)
│   │       └── router/         # 路由注册
│   └── pkg/                    # 公共包
│       ├── errcode/            # 统一错误码
│       ├── response/           # 统一响应
│       └── validator/          # 参数验证器
├── src/                        # 原版代码 (兼容性保留)
│   ├── controller/
│   ├── database/
│   ├── model/
│   ├── router/
│   └── server/
├── deploy/                     # 部署配置
│   ├── Dockerfile
│   └── docker-compose.yml
├── test/                       # 测试
│   ├── unit/
│   └── load/
├── config.yaml                 # 配置文件
└── README.md
```

## 测试

### 单元测试

```bash
# 运行单元测试
go test ./test/unit/... -v

# 运行所有测试
go test ./... -v
```

### 压测 (Go Load Test)

项目使用 Go 编写了并发压测工具，位于 `test/load/` 目录：

```bash
# 进入压测目录
cd test/load

# 选课接口压测 (20000 并发请求, 1000 线程)
go test -v -run TestBookCourseLoad

# 获取课表接口压测
go test -v -run TestGetStudentCourseLoad

# 登录接口压测
go test -v -run TestLoginLoad

# 并发选课详细测试
go test -v -run TestConcurrentBookCourse

# 运行所有压测
go test -v -run "Load$"
```

**压测结果示例 (Top 1 成绩)**:

```
========== Load Test Results ==========
Total Requests:      20000
Success Requests:    19985
Failed Requests:     15
Success Rate:        99.92%
Total Duration:      5.234s
Requests/sec (RPS):  3821.23

Latency:
  Min:    1.235ms
  Avg:    2.567ms
  Max:    45.231ms

Percentiles:
  P50:    2.123ms
  P90:    3.456ms
  P95:    4.789ms
  P99:    8.123ms
========================================
```

**压测配置参数**:
- 并发 Worker 数: 1000
- 总请求数: 20000
- 请求超时: 5 秒
- 目标 RPS: 4000+

## 核心设计

### 选课高并发处理

1. **令牌桶限流**: 使用 `golang.org/x/time/rate` 实现接口限流
2. **Redis 原子操作**: 使用 `HINCRBY` 原子递减课程库存
3. **异步写入**: 选课请求先写入 Redis，通过 RocketMQ 异步同步到 MySQL
4. **超卖防护**: Redis 库存扣减后检查返回值，负数则回滚

```go
// 核心选课逻辑
func (s *SelectionAppService) BookCourse(ctx context.Context, req *BookCourseRequest) error {
    // 1. 限流检查
    if err := s.limiter.Wait(ctx); err != nil {
        return errcode.UnknownError
    }

    // 2. 检查重复选课
    enrolled, _ := s.redis.SIsMember(ctx, studentKey, courseID)
    if enrolled {
        return errcode.RepeatRequest
    }

    // 3. Redis 原子扣减库存
    remaining, _ := s.redis.HIncrBy(ctx, "course:capacity", courseID, -1)
    if remaining < 0 {
        s.redis.HIncrBy(ctx, "course:capacity", courseID, 1) // 回滚
        return errcode.CourseNotAvailable
    }

    // 4. 异步写入 MQ
    s.redis.LPush(ctx, "booking:queue", message)
    return nil
}
```

### 排课算法 (二分图最大匹配)

使用匈牙利算法求解教师与课程的最优匹配：

```go
// 输入: 教师期望课程偏好
TeacherCourseRelationShip: {
    "a": ["1", "4"],
    "b": ["1", "2"],
    "c": ["2"],
    "d": ["3"]
}

// 输出: 教师分配结果
{
    "a": "4",
    "b": "1",
    "c": "2",
    "d": "3"
}
```

## 许可证

MIT License

感谢所有为这个项目做出贡献的团队成员！
