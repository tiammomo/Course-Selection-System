-- Course Selection System Database Schema
-- Version: 1.03
-- Database: course_select
-- Charset: utf8mb4
-- Collation: utf8mb4_unicode_ci

-- ============================================================
-- 1. 创建数据库
-- ============================================================
DROP DATABASE IF EXISTS course_select;
CREATE DATABASE course_select DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
USE course_select;

-- ============================================================
-- 2. 成员表 (member)
--    存储所有用户信息: 管理员、教师、学生
-- ============================================================
DROP TABLE IF EXISTS `member`;
CREATE TABLE `member` (
    `user_id` INT NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` VARCHAR(50) NOT NULL COMMENT '用户名(登录账号)',
    `password` VARCHAR(255) NOT NULL COMMENT '密码(bcrypt加密)',
    `nickname` VARCHAR(100) NOT NULL COMMENT '昵称',
    `user_type` TINYINT NOT NULL DEFAULT 3 COMMENT '用户类型: 1=管理员, 2=教师, 3=学生',
    `is_deleted` TINYINT NOT NULL DEFAULT 0 COMMENT '是否删除: 0=否, 1=是',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`user_id`),
    UNIQUE KEY `uk_username` (`username`),
    KEY `idx_user_type` (`user_type`),
    KEY `idx_is_deleted` (`is_deleted`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户成员表';

-- 插入测试数据 (密码均为: Pass123)
-- 密码说明: bcrypt加密后的值, $2a$10$开头
-- 原始密码 Pass123 对应的 bcrypt 哈希
-- 管理员: JudgeAdmin / Pass123
-- 教师: teacher1 / Pass123
-- 学生: student1 / Pass123
--       student2 / Pass123
LOCK TABLES `member` WRITE;
INSERT INTO `member` VALUES
    (1, 'JudgeAdmin', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '系统管理员', 1, 0, NOW(), NOW()),
    (2, 'teacher1', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '张老师', 2, 0, NOW(), NOW()),
    (3, 'teacher2', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '李老师', 2, 0, NOW(), NOW()),
    (4, 'student1', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '张三', 3, 0, NOW(), NOW()),
    (5, 'student2', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '李四', 3, 0, NOW(), NOW()),
    (6, 'student3', '$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iAt6Z5EHsM8lE9lBOsl7iKTVKIUi', '王五', 3, 0, NOW(), NOW());
UNLOCK TABLES;

-- ============================================================
-- 3. 课程表 (course)
--    存储课程基本信息
-- ============================================================
DROP TABLE IF EXISTS `course`;
CREATE TABLE `course` (
    `course_id` INT NOT NULL AUTO_INCREMENT COMMENT '课程ID',
    `name` VARCHAR(100) NOT NULL COMMENT '课程名称',
    `capacity` INT NOT NULL COMMENT '课程容量(最大选课人数)',
    `teacher_id` INT NULL COMMENT '授课教师ID',
    `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`course_id`),
    KEY `idx_teacher_id` (`teacher_id`),
    CONSTRAINT `fk_course_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `member` (`user_id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 插入测试课程数据
LOCK TABLES `course` WRITE;
INSERT INTO `course` (`course_id`, `name`, `capacity`, `teacher_id`) VALUES
    (1, '高等数学', 100, 2),
    (2, '线性代数', 50, 2),
    (3, '信号与系统', 60, 3),
    (4, '数字信号处理', 40, 3),
    (5, '复变函数', 80, NULL);
UNLOCK TABLES;

-- ============================================================
-- 4. 选课关系表 (choice)
--    学生与课程的选课关系
-- ============================================================
DROP TABLE IF EXISTS `choice`;
CREATE TABLE `choice` (
    `student_id` INT NOT NULL COMMENT '学生ID',
    `course_id` INT NOT NULL COMMENT '课程ID',
    `selected_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '选课时间',
    PRIMARY KEY (`student_id`, `course_id`),
    KEY `idx_choice_course` (`course_id`),
    CONSTRAINT `fk_choice_student` FOREIGN KEY (`student_id`) REFERENCES `member` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_choice_course` FOREIGN KEY (`course_id`) REFERENCES `course` (`course_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生选课关系表';

-- 插入测试选课数据
LOCK TABLES `choice` WRITE;
INSERT INTO `choice` VALUES
    (4, 1, NOW()),
    (4, 2, NOW()),
    (5, 1, NOW()),
    (5, 3, NOW()),
    (6, 2, NOW()),
    (6, 4, NOW());
UNLOCK TABLES;

-- ============================================================
-- 5. 教师-课程绑定表 (bind)
--    教师授课关系 (用于排课)
-- ============================================================
DROP TABLE IF EXISTS `bind`;
CREATE TABLE `bind` (
    `teacher_id` INT NOT NULL COMMENT '教师ID',
    `course_id` INT NOT NULL COMMENT '课程ID',
    `assigned_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '绑定时间',
    PRIMARY KEY (`teacher_id`, `course_id`),
    KEY `idx_bind_course` (`course_id`),
    CONSTRAINT `fk_bind_teacher` FOREIGN KEY (`teacher_id`) REFERENCES `member` (`user_id`) ON DELETE CASCADE,
    CONSTRAINT `fk_bind_course` FOREIGN KEY (`course_id`) REFERENCES `course` (`course_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师课程绑定表(排课用)';

-- 插入测试绑定数据
LOCK TABLES `bind` WRITE;
INSERT INTO `bind` VALUES
    (2, 1, NOW()),
    (2, 2, NOW()),
    (3, 3, NOW()),
    (3, 4, NOW());
UNLOCK TABLES;

-- ============================================================
-- 6. 查看已选课人数的视图
-- ============================================================
DROP VIEW IF EXISTS v_course_selection;
CREATE VIEW v_course_selection AS
SELECT
    c.course_id,
    c.name,
    c.capacity,
    COUNT(ch.student_id) AS selected_count
FROM course c
LEFT JOIN choice ch ON c.course_id = ch.course_id
GROUP BY c.course_id, c.name, c.capacity;

-- ============================================================
-- 7. 查看用户详情的视图
-- ============================================================
DROP VIEW IF EXISTS v_member_detail;
CREATE VIEW v_member_detail AS
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

-- ============================================================
-- 8. 初始化 Redis 缓存数据 (可选)
--    以下为课程剩余容量的初始值，实际运行时由应用维护
-- ============================================================
-- Redis Key: course:capacity
-- Redis Hash Structure: { course_id: remaining_count }
-- 示例: HSET course:capacity 1 100 2 50 3 60 4 40 5 80
--
-- Redis Key: student:{student_id}:courses
-- Redis Set Structure: { course_id1, course_id2, ... }
-- 示例: SADD student:4:courses 1 2
--       SADD student:5:courses 1 3
