# ER 图说明

## 数据库架构

![ER Diagram](https://s3.bmp.ovh/imgs/2022/01/5b77d2f959361d2e.png)

### 实体集

| 实体集 | 说明 |
|--------|------|
| **Member** | 用户成员实体集，包含管理员、教师、学生三种角色 |
| **Course** | 课程实体集，存储课程基本信息 |

### 联系集

| 联系集 | 类型 | 说明 |
|--------|------|------|
| **Bind** | 1:1 | 教师授课绑定，一个教师绑定一门课程，一门课程由一个教师教授 |
| **Choice** | M:N | 学生选课关系，多个学生可选同一门课，一个学生可选多门课 |

---

## 数据表结构

### 1. Member 表 (用户成员表)

```
Member(UserID, Username, Password, Nickname, UserType, IsDeleted, CreatedAt, UpdatedAt)
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| UserID | INT | PRIMARY KEY, AUTO_INCREMENT | 用户ID |
| Username | VARCHAR(50) | UNIQUE, NOT NULL | 用户名（登录账号） |
| Password | VARCHAR(255) | NOT NULL | 密码（bcrypt加密） |
| Nickname | VARCHAR(100) | NOT NULL | 用户昵称 |
| UserType | TINYINT | NOT NULL | 用户类型 |
| IsDeleted | TINYINT | NOT NULL DEFAULT 0 | 是否删除（软删除） |
| CreatedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| UpdatedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE | 更新时间 |

**UserType 枚举值**：

| 值 | 类型 | 说明 |
|----|------|------|
| 1 | 管理员 | 系统管理员，拥有最高权限 |
| 2 | 教师 | 授课教师，可管理所授课程 |
| 3 | 学生 | 选课学生，可进行选课操作 |

---

### 2. Course 表 (课程表)

```
Course(CourseID, Name, Capacity, TeacherID, CreatedAt, UpdatedAt)
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| CourseID | INT | PRIMARY KEY, AUTO_INCREMENT | 课程ID |
| Name | VARCHAR(100) | NOT NULL | 课程名称 |
| Capacity | INT | NOT NULL | 课程容量（最大选课人数） |
| TeacherID | INT | FOREIGN KEY → Member(UserID) | 授课教师ID |
| CreatedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 创建时间 |
| UpdatedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP ON UPDATE | 更新时间 |

---

### 3. Choice 表 (学生选课关系表)

```
Choice(StudentID, CourseID, SelectedAt)
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| StudentID | INT | FOREIGN KEY → Member(UserID), PRIMARY KEY | 学生ID |
| CourseID | INT | FOREIGN KEY → Course(CourseID), PRIMARY KEY | 课程ID |
| SelectedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 选课时间 |

**级联删除策略**：删除学生或课程时，自动删除相关选课记录

---

### 4. Bind 表 (教师课程绑定表)

```
Bind(TeacherID, CourseID, AssignedAt)
```

| 字段 | 类型 | 约束 | 说明 |
|------|------|------|------|
| TeacherID | INT | FOREIGN KEY → Member(UserID), PRIMARY KEY | 教师ID |
| CourseID | INT | FOREIGN KEY → Course(CourseID), PRIMARY KEY | 课程ID |
| AssignedAt | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP | 绑定时间 |

**级联删除策略**：删除教师或课程时，自动删除相关绑定记录

---

## 视图

### v_course_selection (课程选课统计视图)

```sql
SELECT
    c.course_id,
    c.name,
    c.capacity,
    COUNT(ch.student_id) AS selected_count
FROM course c
LEFT JOIN choice ch ON c.course_id = ch.course_id
GROUP BY c.course_id, c.name, c.capacity;
```

### v_member_detail (用户详情视图)

```sql
SELECT
    user_id,
    username,
    nickname,
    CASE user_type
        WHEN 1 THEN '管理员'
        WHEN 2 THEN '教师'
        WHEN 3 THEN '学生'
        ELSE '未知'
    END AS user_type_name,
    is_deleted,
    created_at,
    updated_at
FROM member;
```

---

## Redis 缓存数据结构

### 课程剩余容量 (Hash)

```
Key: course:capacity
Field: course_id
Value: remaining_count (剩余可选数量)
```

**示例**：
```
HSET course:capacity 1 100 2 50 3 60 4 40 5 80
```

### 学生已选课程 (Set)

```
Key: student:{student_id}:courses
Value: course_id (课程ID集合)
```

**示例**：
```
SADD student:4:courses 1 2
SADD student:5:courses 1 3
```

### Session 存储

```
Key: session:{sessionId}
Value: JSON {
    "user_id": "1",
    "username": "admin",
    "nickname": "管理员",
    "user_type": 1
}
TTL: 3600 seconds
```

---

## ER 图关系说明

```
┌─────────────┐       ┌─────────────┐       ┌─────────────┐
│   Member    │       │   Course    │       │   Choice    │
├─────────────┤       ├─────────────┤       ├─────────────┤
│ UserID (PK) │───┐   │ CourseID(PK)│   ┌───│ StudentID   │
│ Username    │   │   │ Name        │   │   │ CourseID    │
│ Password    │   │   │ Capacity    │   │   │ SelectedAt  │
│ Nickname    │   │   │ TeacherID ──┼───┘   └─────────────┘
│ UserType    │   │   └─────────────┘
│ IsDeleted   │   │
└─────────────┘   │
      │           │
      │ Bind      │
      ▼           │
┌─────────────┐   │
│   Bind      │◄──┘
├─────────────┤
│ TeacherID   │
│ CourseID    │
│ AssignedAt  │
└─────────────┘
```

### 关系说明

| 关系 | 说明 |
|------|------|
| Member → Course (via TeacherID) | 一位教师可教授多门课程，一门课程有一位授课教师 |
| Member ↔ Choice (via StudentID) | 一位学生可选多门课程，一门课程可被多位学生选择 |
| Member ↔ Bind (via TeacherID) | 一位教师可绑定多门课程，一门课程只能被一位教师绑定 |

---

## 索引优化

| 表 | 索引类型 | 索引字段 | 说明 |
|----|----------|----------|------|
| Member | UNIQUE | username | 防止重复用户名 |
| Member | INDEX | user_type | 按用户类型快速查询 |
| Member | INDEX | is_deleted | 按删除状态快速筛选 |
| Course | INDEX | teacher_id | 按教师快速查询课程 |
| Choice | INDEX | course_id | 按课程快速查询选课记录 |
| Bind | INDEX | course_id | 按课程快速查询绑定关系 |
