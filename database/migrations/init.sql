-- GinForge 数据库初始化脚本
-- 表前缀: gf_
-- 执行方式: make db-init

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ============================================
-- 1. 管理后台相关表
-- ============================================

-- 1.1 管理员用户表
CREATE TABLE IF NOT EXISTS `gf_admin_users` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `email` varchar(100) NOT NULL COMMENT '邮箱',
    `phone` varchar(20) DEFAULT NULL COMMENT '手机号',
    `password` varchar(255) NOT NULL COMMENT '密码',
    `name` varchar(50) DEFAULT NULL COMMENT '真实姓名',
    `avatar` varchar(255) DEFAULT NULL COMMENT '头像URL',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
    `last_login_at` timestamp NULL DEFAULT NULL COMMENT '最后登录时间',
    `last_login_ip` varchar(45) DEFAULT NULL COMMENT '最后登录IP',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_username` (`username`),
    UNIQUE KEY `uk_email` (`email`),
    KEY `idx_phone` (`phone`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='管理员用户表';

-- 1.2 角色表
CREATE TABLE IF NOT EXISTS `gf_admin_roles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '角色ID',
    `name` varchar(50) NOT NULL COMMENT '角色名称',
    `code` varchar(50) NOT NULL COMMENT '角色编码',
    `description` text COMMENT '角色描述',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_name` (`name`),
    KEY `idx_status` (`status`),
    KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- 1.3 权限表
CREATE TABLE IF NOT EXISTS `gf_admin_permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `name` varchar(50) NOT NULL COMMENT '权限名称',
    `code` varchar(100) NOT NULL COMMENT '权限编码',
    `type` varchar(20) NOT NULL DEFAULT 'menu' COMMENT '权限类型: menu-菜单, button-按钮, api-接口',
    `description` text COMMENT '权限描述',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_name` (`name`),
    KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 1.4 菜单表
CREATE TABLE IF NOT EXISTS `gf_admin_menus` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '菜单ID',
    `parent_id` bigint(20) unsigned DEFAULT '0' COMMENT '父菜单ID',
    `name` varchar(50) NOT NULL COMMENT '菜单名称',
    `code` varchar(50) NOT NULL COMMENT '菜单编码',
    `type` varchar(20) NOT NULL DEFAULT 'menu' COMMENT '菜单类型: directory-目录, menu-菜单, button-按钮',
    `path` varchar(200) DEFAULT NULL COMMENT '路由路径',
    `component` varchar(200) DEFAULT NULL COMMENT '组件路径',
    `icon` varchar(50) DEFAULT NULL COMMENT '菜单图标',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `visible` tinyint(1) NOT NULL DEFAULT '1' COMMENT '是否显示: 1-显示, 0-隐藏',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用, 0-禁用',
    `description` text COMMENT '菜单描述',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_parent_id` (`parent_id`),
    KEY `idx_name` (`name`),
    KEY `idx_type` (`type`),
    KEY `idx_status` (`status`),
    KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- 1.5 用户角色关联表
CREATE TABLE IF NOT EXISTS `gf_admin_user_roles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role_id` (`role_id`),
    CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `gf_admin_users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `gf_admin_roles` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 1.6 角色权限关联表
CREATE TABLE IF NOT EXISTS `gf_admin_role_permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `permission_id` bigint(20) unsigned NOT NULL COMMENT '权限ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_permission_id` (`permission_id`),
    CONSTRAINT `fk_role_permissions_role` FOREIGN KEY (`role_id`) REFERENCES `gf_admin_roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_permissions_permission` FOREIGN KEY (`permission_id`) REFERENCES `gf_admin_permissions` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 1.7 角色菜单关联表
CREATE TABLE IF NOT EXISTS `gf_admin_role_menus` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `menu_id` bigint(20) unsigned NOT NULL COMMENT '菜单ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_menu_id` (`menu_id`),
    CONSTRAINT `fk_role_menus_role` FOREIGN KEY (`role_id`) REFERENCES `gf_admin_roles` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_role_menus_menu` FOREIGN KEY (`menu_id`) REFERENCES `gf_admin_menus` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ============================================
-- 2. 系统配置相关表
-- ============================================

-- 2.1 系统配置表
CREATE TABLE IF NOT EXISTS `gf_admin_system_configs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '配置ID',
    `key` varchar(100) NOT NULL COMMENT '配置键',
    `value` text COMMENT '配置值',
    `type` varchar(20) NOT NULL DEFAULT 'string' COMMENT '配置类型: string, number, boolean, json',
    `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
    `group` varchar(50) DEFAULT 'default' COMMENT '配置分组',
    `sort` int(11) NOT NULL DEFAULT '0' COMMENT '排序',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_key` (`key`),
    KEY `idx_group` (`group`),
    KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='系统配置表';

-- 2.2 操作日志表
CREATE TABLE IF NOT EXISTS `gf_admin_operation_logs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作用户ID',
    `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
    `method` varchar(10) NOT NULL COMMENT '请求方法',
    `path` varchar(255) NOT NULL COMMENT '请求路径',
    `ip` varchar(45) DEFAULT NULL COMMENT '请求IP',
    `user_agent` text COMMENT '用户代理',
    `request_data` json DEFAULT NULL COMMENT '请求数据',
    `response_data` json DEFAULT NULL COMMENT '响应数据',
    `status_code` int(11) NOT NULL DEFAULT '200' COMMENT '响应状态码',
    `duration` int(11) NOT NULL DEFAULT '0' COMMENT '请求耗时(毫秒)',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_method` (`method`),
    KEY `idx_path` (`path`),
    KEY `idx_ip` (`ip`),
    KEY `idx_status_code` (`status_code`),
    KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='操作日志表';

-- 2.3 登录记录表
CREATE TABLE IF NOT EXISTS `gf_admin_login_logs` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `username` varchar(50) NOT NULL COMMENT '用户名',
    `login_ip` varchar(45) DEFAULT NULL COMMENT '登录IP',
    `user_agent` text COMMENT '用户代理',
    `login_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '登录时间',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '登录状态: 1-成功, 0-失败',
    `failure_reason` varchar(255) DEFAULT NULL COMMENT '失败原因',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_username` (`username`),
    KEY `idx_login_time` (`login_time`),
    KEY `idx_status` (`status`),
    KEY `idx_login_ip` (`login_ip`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='登录记录表';

-- ============================================
-- 3. 文件服务相关表
-- ============================================

-- 3.1 文件上传日志表
CREATE TABLE IF NOT EXISTS `gf_file_upload_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    `file_name` varchar(255) NOT NULL COMMENT '文件名',
    `file_size` bigint DEFAULT NULL COMMENT '文件大小',
    `uploaded_by` bigint unsigned DEFAULT NULL COMMENT '上传用户ID',
    `upload_ip` varchar(50) DEFAULT NULL COMMENT '上传IP',
    `upload_time` datetime DEFAULT NULL COMMENT '上传时间',
    `success` tinyint(1) DEFAULT NULL COMMENT '是否成功: 1-成功, 0-失败',
    `error_message` text COMMENT '错误信息',
    `duration` bigint DEFAULT NULL COMMENT '上传耗时(毫秒)',
    `upload_type` varchar(20) DEFAULT NULL COMMENT '上传类型: simple-普通上传, chunked-分片上传',
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_uploaded_by` (`uploaded_by`),
    KEY `idx_upload_time` (`upload_time`),
    KEY `idx_success` (`success`),
    KEY `idx_upload_type` (`upload_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件上传日志表';

-- 3.2 文件下载日志表
CREATE TABLE IF NOT EXISTS `gf_file_download_logs` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '日志ID',
    `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime(3) DEFAULT NULL COMMENT '删除时间',
    `file_id` bigint unsigned DEFAULT NULL COMMENT '文件ID',
    `file_name` varchar(255) DEFAULT NULL COMMENT '文件名',
    `downloaded_by` bigint unsigned DEFAULT NULL COMMENT '下载用户ID',
    `download_ip` varchar(50) DEFAULT NULL COMMENT '下载IP',
    `download_time` datetime DEFAULT NULL COMMENT '下载时间',
    `success` tinyint(1) DEFAULT NULL COMMENT '是否成功: 1-成功, 0-失败',
    PRIMARY KEY (`id`),
    KEY `idx_deleted_at` (`deleted_at`),
    KEY `idx_file_id` (`file_id`),
    KEY `idx_downloaded_by` (`downloaded_by`),
    KEY `idx_download_time` (`download_time`),
    KEY `idx_success` (`success`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文件下载日志表';

-- ============================================
-- 4. 业务表（示例）
-- ============================================

-- 4.1 文章表
CREATE TABLE IF NOT EXISTS `gf_articles` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '文章ID',
    `title` varchar(255) NOT NULL COMMENT '标题',
    `content` text NOT NULL COMMENT '内容',
    `summary` varchar(500) DEFAULT NULL COMMENT '摘要',
    `author_id` bigint unsigned NOT NULL COMMENT '作者ID',
    `author_name` varchar(100) DEFAULT NULL COMMENT '作者名称',
    `category_id` bigint unsigned DEFAULT NULL COMMENT '分类ID',
    `cover_image` varchar(500) DEFAULT NULL COMMENT '封面图片',
    `tags` varchar(255) DEFAULT NULL COMMENT '标签（逗号分隔）',
    `status` tinyint NOT NULL DEFAULT '0' COMMENT '状态:0草稿,1已发布,2已下线',
    `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶',
    `view_count` int unsigned NOT NULL DEFAULT '0' COMMENT '浏览次数',
    `like_count` int unsigned NOT NULL DEFAULT '0' COMMENT '点赞次数',
    `comment_count` int unsigned NOT NULL DEFAULT '0' COMMENT '评论次数',
    `published_at` datetime DEFAULT NULL COMMENT '发布时间',
    `created_at` datetime DEFAULT NULL COMMENT '创建时间',
    `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
    `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    KEY `idx_author` (`author_id`),
    KEY `idx_category` (`category_id`),
    KEY `idx_status` (`status`),
    KEY `idx_created` (`created_at`),
    KEY `idx_published` (`published_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='文章表';

SET FOREIGN_KEY_CHECKS = 1;

-- ============================================
-- 插入初始数据
-- ============================================

-- 插入默认管理员用户 (密码: admin123)
INSERT INTO `gf_admin_users` (`username`, `email`, `password`, `name`, `status`) VALUES
('admin', 'admin@ginforge.com', '$2a$10$3Idb4DTJjoxAsyo6SD2NX.oBckKvrRheHS8V13yDl7LWNNvnQIkgG', '超级管理员', 1)
ON DUPLICATE KEY UPDATE `username`=`username`;

-- 插入默认角色
INSERT INTO `gf_admin_roles` (`name`, `code`, `description`, `sort`, `status`) VALUES
('超级管理员', 'super_admin', '拥有所有权限的超级管理员角色', 1, 1),
('管理员', 'admin', '普通管理员角色', 2, 1),
('普通用户', 'user', '普通用户角色', 3, 1)
ON DUPLICATE KEY UPDATE `code`=`code`;

-- 插入默认权限
INSERT INTO `gf_admin_permissions` (`name`, `code`, `type`, `description`, `status`) VALUES
-- 用户管理权限
('用户查看', 'user:read', 'api', '查看用户列表', 1),
('用户创建', 'user:create', 'api', '创建新用户', 1),
('用户编辑', 'user:update', 'api', '编辑用户信息', 1),
('用户删除', 'user:delete', 'api', '删除用户', 1),
-- 角色管理权限
('角色查看', 'role:read', 'api', '查看角色列表', 1),
('角色创建', 'role:create', 'api', '创建新角色', 1),
('角色编辑', 'role:update', 'api', '编辑角色信息', 1),
('角色删除', 'role:delete', 'api', '删除角色', 1),
-- 菜单管理权限
('菜单查看', 'menu:read', 'api', '查看菜单列表', 1),
('菜单创建', 'menu:create', 'api', '创建新菜单', 1),
('菜单编辑', 'menu:update', 'api', '编辑菜单信息', 1),
('菜单删除', 'menu:delete', 'api', '删除菜单', 1),
-- 权限管理权限
('权限查看', 'permission:read', 'api', '查看权限列表', 1),
('权限创建', 'permission:create', 'api', '创建新权限', 1),
('权限编辑', 'permission:update', 'api', '编辑权限信息', 1),
('权限删除', 'permission:delete', 'api', '删除权限', 1),
-- 系统管理权限
('系统查看', 'system:read', 'api', '查看系统信息', 1),
('系统配置', 'system:config', 'api', '系统配置管理', 1),
-- 文章管理权限
('文章查看', 'article:read', 'api', '查看文章列表', 1),
('文章创建', 'article:create', 'api', '创建新文章', 1),
('文章编辑', 'article:update', 'api', '编辑文章信息', 1),
('文章删除', 'article:delete', 'api', '删除文章', 1)
ON DUPLICATE KEY UPDATE `code`=`code`;

-- 插入默认菜单（先插入一级菜单）
INSERT INTO `gf_admin_menus` (`parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `sort`, `visible`, `status`, `description`) VALUES
-- 一级菜单
(0, '仪表盘', 'dashboard', 'menu', '/dashboard', 'Dashboard', 'House', 1, 1, 1, '系统仪表盘'),
(0, '系统管理', 'system_management', 'directory', '', '', 'Setting', 2, 1, 1, '系统管理模块'),
(0, '文章管理', 'articles_management', 'menu', '/dashboard/articleses', 'Articles', 'Document', 3, 1, 1, '文章管理')
ON DUPLICATE KEY UPDATE `code`=`code`;

-- 插入系统管理下的二级菜单（需要先获取系统管理菜单ID）
SET @system_menu_id = (SELECT id FROM gf_admin_menus WHERE code = 'system_management' LIMIT 1);

INSERT INTO `gf_admin_menus` (`parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `sort`, `visible`, `status`, `description`) VALUES
-- 系统管理下的二级菜单
(@system_menu_id, '系统配置', 'system_config', 'menu', '/dashboard/system', 'System', 'Setting', 1, 1, 1, '系统配置管理'),
(@system_menu_id, '用户管理', 'user_management', 'directory', '', '', 'User', 2, 1, 1, '用户管理模块'),
(@system_menu_id, '角色管理', 'role_management', 'directory', '', '', 'UserFilled', 3, 1, 1, '角色管理模块'),
(@system_menu_id, '菜单管理', 'menu_management', 'directory', '', '', 'Menu', 4, 1, 1, '菜单管理模块'),
(@system_menu_id, '权限管理', 'permission_management', 'menu', '/dashboard/permissions', 'Permissions', 'Key', 5, 1, 1, '权限管理')
ON DUPLICATE KEY UPDATE `code`=`code`;

-- 插入三级菜单（用户管理、角色管理、菜单管理的子菜单）
SET @user_menu_id = (SELECT id FROM gf_admin_menus WHERE code = 'user_management' LIMIT 1);
SET @role_menu_id = (SELECT id FROM gf_admin_menus WHERE code = 'role_management' LIMIT 1);
SET @menu_menu_id = (SELECT id FROM gf_admin_menus WHERE code = 'menu_management' LIMIT 1);

INSERT INTO `gf_admin_menus` (`parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `sort`, `visible`, `status`, `description`) VALUES
(@user_menu_id, '用户列表', 'user_list', 'menu', '/dashboard/users', 'Users', '', 1, 1, 1, '用户列表页面'),
(@user_menu_id, '创建用户', 'user_create', 'button', '', '', '', 2, 0, 1, '创建用户按钮'),
(@role_menu_id, '角色列表', 'role_list', 'menu', '/dashboard/roles', 'Roles', '', 1, 1, 1, '角色列表页面'),
(@role_menu_id, '创建角色', 'role_create', 'button', '', '', '', 2, 0, 1, '创建角色按钮'),
(@menu_menu_id, '菜单列表', 'menu_list', 'menu', '/dashboard/menus', 'Menus', '', 1, 1, 1, '菜单列表页面'),
(@menu_menu_id, '创建菜单', 'menu_create', 'button', '', '', '', 2, 0, 1, '创建菜单按钮')
ON DUPLICATE KEY UPDATE `code`=`code`;

-- 关联管理员用户和超级管理员角色
INSERT IGNORE INTO `gf_admin_user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 为超级管理员角色分配所有权限
INSERT IGNORE INTO `gf_admin_role_permissions` (`role_id`, `permission_id`) 
SELECT 1, id FROM `gf_admin_permissions`;

-- 为超级管理员角色分配所有菜单
INSERT IGNORE INTO `gf_admin_role_menus` (`role_id`, `menu_id`) 
SELECT 1, id FROM `gf_admin_menus`;

-- 插入默认系统配置
INSERT INTO `gf_admin_system_configs` (`key`, `value`, `type`, `description`, `group`, `sort`) VALUES
-- 基本配置
('system.name', 'GinForge 管理后台', 'string', '系统名称', 'basic', 10),
('system.version', '1.0.0', 'string', '系统版本', 'basic', 20),
('system.description', '基于 Go + Gin 的企业级微服务开发框架', 'string', '系统描述', 'basic', 30),
('system.logo', '/logo.svg', 'string', '系统Logo', 'basic', 40),
('system.default_language', 'zh-CN', 'string', '默认语言', 'basic', 50),
-- 安全配置
('security.min_password_length', '8', 'number', '密码最小长度', 'security', 10),
('security.password_complexity', '["lowercase","numbers"]', 'json', '密码复杂度要求', 'security', 20),
('security.max_login_attempts', '5', 'number', '最大登录失败次数', 'security', 30),
('security.lockout_duration', '15', 'number', '账户锁定时间(分钟)', 'security', 40),
('security.session_timeout', '120', 'number', '会话超时时间(分钟)', 'security', 50),
-- 邮件配置
('email.smtp_host', 'smtp.example.com', 'string', 'SMTP服务器地址', 'email', 10),
('email.smtp_port', '587', 'number', 'SMTP服务器端口', 'email', 20),
('email.from_email', 'noreply@example.com', 'string', '发送邮箱', 'email', 30),
('email.email_password', '', 'string', '邮箱密码', 'email', 40),
('email.enable_ssl', 'true', 'boolean', '启用SSL', 'email', 50),
-- 缓存配置
('cache.type', 'redis', 'string', '缓存类型', 'cache', 10),
('cache.redis_host', 'localhost', 'string', 'Redis地址', 'cache', 20),
('cache.redis_port', '6379', 'number', 'Redis端口', 'cache', 30),
('cache.redis_password', '', 'string', 'Redis密码', 'cache', 40),
('cache.default_expiration', '3600', 'number', '默认过期时间(秒)', 'cache', 50)
ON DUPLICATE KEY UPDATE `key`=`key`;

-- 插入测试文章数据
INSERT INTO `gf_articles` (`title`, `content`, `summary`, `author_id`, `author_name`, `status`, `is_top`, `view_count`, `created_at`, `updated_at`) VALUES
('欢迎使用 GinForge 框架', '# GinForge 框架介绍\n\nGinForge 是一个强大的 Go Web 开发框架...', '快速了解 GinForge 框架的核心特性', 1, 'Admin', 1, 1, 100, NOW(), NOW()),
('如何使用代码生成器', '# 代码生成器使用指南\n\n代码生成器可以帮助您快速生成 CRUD 代码...', '5分钟学会使用代码生成器', 1, 'Admin', 1, 0, 50, NOW(), NOW()),
('GinForge 最佳实践', '# 最佳实践\n\n在使用 GinForge 开发时，建议遵循以下最佳实践...', '提升开发效率的最佳实践', 1, 'Admin', 0, 0, 20, NOW(), NOW())
ON DUPLICATE KEY UPDATE `title`=`title`;

-- 初始化完成提示
SELECT 'Database initialization completed successfully!' AS message;

