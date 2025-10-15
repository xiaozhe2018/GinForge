-- 创建管理后台相关表结构
-- 数据库: gin_forge

-- 1. 用户表
CREATE TABLE IF NOT EXISTS `admin_users` (
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

-- 2. 角色表
CREATE TABLE IF NOT EXISTS `admin_roles` (
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

-- 3. 权限表
CREATE TABLE IF NOT EXISTS `admin_permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '权限ID',
    `name` varchar(50) NOT NULL COMMENT '权限名称',
    `code` varchar(100) NOT NULL COMMENT '权限编码',
    `type` varchar(20) NOT NULL DEFAULT 'menu' COMMENT '权限类型: menu-菜单, button-按钮, api-接口',
    `description` text COMMENT '权限描述',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_code` (`code`),
    KEY `idx_name` (`name`),
    KEY `idx_type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='权限表';

-- 4. 菜单表
CREATE TABLE IF NOT EXISTS `admin_menus` (
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

-- 5. 用户角色关联表
CREATE TABLE IF NOT EXISTS `admin_user_roles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `user_id` bigint(20) unsigned NOT NULL COMMENT '用户ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_user_role` (`user_id`, `role_id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_role_id` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户角色关联表';

-- 6. 角色权限关联表
CREATE TABLE IF NOT EXISTS `admin_role_permissions` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `permission_id` bigint(20) unsigned NOT NULL COMMENT '权限ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_permission_id` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色权限关联表';

-- 7. 角色菜单关联表
CREATE TABLE IF NOT EXISTS `admin_role_menus` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `role_id` bigint(20) unsigned NOT NULL COMMENT '角色ID',
    `menu_id` bigint(20) unsigned NOT NULL COMMENT '菜单ID',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_role_menu` (`role_id`, `menu_id`),
    KEY `idx_role_id` (`role_id`),
    KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- 8. 系统配置表
CREATE TABLE IF NOT EXISTS `admin_system_configs` (
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

-- 9. 操作日志表
CREATE TABLE IF NOT EXISTS `admin_operation_logs` (
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

-- 插入初始数据

-- 插入默认管理员用户 (密码: admin123)
INSERT INTO `admin_users` (`username`, `email`, `password`, `name`, `status`) VALUES
('admin', 'admin@ginforge.com', '$2a$10$3Idb4DTJjoxAsyo6SD2NX.oBckKvrRheHS8V13yDl7LWNNvnQIkgG', '超级管理员', 1);

-- 插入默认角色
INSERT INTO `admin_roles` (`name`, `code`, `description`, `sort`, `status`) VALUES
('超级管理员', 'super_admin', '拥有所有权限的超级管理员角色', 1, 1),
('管理员', 'admin', '普通管理员角色', 2, 1),
('普通用户', 'user', '普通用户角色', 3, 1);

-- 插入默认权限
INSERT INTO `admin_permissions` (`name`, `code`, `type`, `description`) VALUES
-- 用户管理权限
('用户查看', 'user:read', 'api', '查看用户列表'),
('用户创建', 'user:create', 'api', '创建新用户'),
('用户编辑', 'user:update', 'api', '编辑用户信息'),
('用户删除', 'user:delete', 'api', '删除用户'),
-- 角色管理权限
('角色查看', 'role:read', 'api', '查看角色列表'),
('角色创建', 'role:create', 'api', '创建新角色'),
('角色编辑', 'role:update', 'api', '编辑角色信息'),
('角色删除', 'role:delete', 'api', '删除角色'),
-- 菜单管理权限
('菜单查看', 'menu:read', 'api', '查看菜单列表'),
('菜单创建', 'menu:create', 'api', '创建新菜单'),
('菜单编辑', 'menu:update', 'api', '编辑菜单信息'),
('菜单删除', 'menu:delete', 'api', '删除菜单'),
-- 权限管理权限
('权限查看', 'permission:read', 'api', '查看权限列表'),
('权限创建', 'permission:create', 'api', '创建新权限'),
('权限编辑', 'permission:update', 'api', '编辑权限信息'),
('权限删除', 'permission:delete', 'api', '删除权限'),
-- 系统管理权限
('系统查看', 'system:read', 'api', '查看系统信息'),
('系统配置', 'system:config', 'api', '系统配置管理');

-- 插入默认菜单
INSERT INTO `admin_menus` (`parent_id`, `name`, `code`, `type`, `path`, `component`, `icon`, `sort`, `visible`, `status`, `description`) VALUES
-- 一级菜单
(0, '仪表盘', 'dashboard', 'menu', '/dashboard', 'Dashboard', 'House', 1, 1, 1, '系统仪表盘'),
(0, '用户管理', 'user_management', 'directory', '', '', 'User', 2, 1, 1, '用户管理模块'),
(0, '角色管理', 'role_management', 'directory', '', '', 'UserFilled', 3, 1, 1, '角色管理模块'),
(0, '菜单管理', 'menu_management', 'directory', '', '', 'Menu', 4, 1, 1, '菜单管理模块'),
(0, '权限管理', 'permission_management', 'menu', '/dashboard/permissions', 'Permissions', 'Key', 5, 1, 1, '权限管理'),
(0, '系统管理', 'system_management', 'menu', '/dashboard/system', 'System', 'Setting', 6, 1, 1, '系统管理'),
-- 二级菜单
(2, '用户列表', 'user_list', 'menu', '/dashboard/users', 'Users', '', 1, 1, 1, '用户列表页面'),
(2, '创建用户', 'user_create', 'button', '', '', '', 2, 0, 1, '创建用户按钮'),
(3, '角色列表', 'role_list', 'menu', '/dashboard/roles', 'Roles', '', 1, 1, 1, '角色列表页面'),
(3, '创建角色', 'role_create', 'button', '', '', '', 2, 0, 1, '创建角色按钮'),
(4, '菜单列表', 'menu_list', 'menu', '/dashboard/menus', 'Menus', '', 1, 1, 1, '菜单列表页面'),
(4, '创建菜单', 'menu_create', 'button', '', '', '', 2, 0, 1, '创建菜单按钮');

-- 关联管理员用户和超级管理员角色
INSERT INTO `admin_user_roles` (`user_id`, `role_id`) VALUES (1, 1);

-- 为超级管理员角色分配所有权限
INSERT INTO `admin_role_permissions` (`role_id`, `permission_id`) 
SELECT 1, id FROM `admin_permissions`;

-- 为超级管理员角色分配所有菜单
INSERT INTO `admin_role_menus` (`role_id`, `menu_id`) 
SELECT 1, id FROM `admin_menus`;

-- 插入默认系统配置
INSERT INTO `admin_system_configs` (`key`, `value`, `type`, `description`, `group`, `sort`) VALUES
('site_name', 'GinForge 管理后台', 'string', '网站名称', 'basic', 1),
('site_logo', '/logo.svg', 'string', '网站Logo', 'basic', 2),
('site_favicon', '/favicon.ico', 'string', '网站图标', 'basic', 3),
('login_captcha', 'true', 'boolean', '登录验证码', 'security', 1),
('password_min_length', '6', 'number', '密码最小长度', 'security', 2),
('session_timeout', '7200', 'number', '会话超时时间(秒)', 'security', 3);
