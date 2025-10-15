-- 创建系统配置表
CREATE TABLE IF NOT EXISTS `admin_system_configs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `key` varchar(100) NOT NULL COMMENT '配置键',
  `value` text DEFAULT NULL COMMENT '配置值',
  `type` varchar(20) DEFAULT 'string' COMMENT '配置类型:string,number,boolean,json',
  `description` varchar(255) DEFAULT NULL COMMENT '配置描述',
  `group` varchar(50) DEFAULT 'default' COMMENT '配置分组',
  `sort` int(11) DEFAULT 0 COMMENT '排序',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_key` (`key`),
  KEY `idx_group` (`group`),
  KEY `idx_sort` (`sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='系统配置表';

-- 创建操作日志表
CREATE TABLE IF NOT EXISTS `admin_operation_logs` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(50) DEFAULT NULL COMMENT '操作用户名',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `path` varchar(255) NOT NULL COMMENT '请求路径',
  `ip` varchar(45) DEFAULT NULL COMMENT '请求IP',
  `user_agent` text DEFAULT NULL COMMENT '用户代理',
  `request_data` json DEFAULT NULL COMMENT '请求数据',
  `response_data` json DEFAULT NULL COMMENT '响应数据',
  `status_code` int(11) DEFAULT 200 COMMENT '响应状态码',
  `duration` int(11) DEFAULT 0 COMMENT '请求耗时(毫秒)',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_method` (`method`),
  KEY `idx_path` (`path`),
  KEY `idx_ip` (`ip`),
  KEY `idx_status_code` (`status_code`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='操作日志表';

-- 插入默认系统配置
INSERT INTO `admin_system_configs` (`key`, `value`, `type`, `description`, `group`, `sort`) VALUES
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
('cache.default_expiration', '3600', 'number', '默认过期时间(秒)', 'cache', 50);
