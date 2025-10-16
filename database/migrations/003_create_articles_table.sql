-- 创建文章表
CREATE TABLE IF NOT EXISTS `articles` (
    `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '文章ID',
    `title` varchar(200) NOT NULL COMMENT '文章标题',
    `slug` varchar(200) DEFAULT NULL COMMENT 'URL别名',
    `author_id` bigint(20) unsigned NOT NULL COMMENT '作者ID',
    `author_name` varchar(50) DEFAULT NULL COMMENT '作者名称',
    `category_id` bigint(20) unsigned DEFAULT NULL COMMENT '分类ID',
    `summary` varchar(500) DEFAULT NULL COMMENT '文章摘要',
    `content` longtext NOT NULL COMMENT '文章内容',
    `cover_image` varchar(255) DEFAULT NULL COMMENT '封面图片',
    `view_count` int(11) NOT NULL DEFAULT '0' COMMENT '浏览次数',
    `like_count` int(11) NOT NULL DEFAULT '0' COMMENT '点赞次数',
    `comment_count` int(11) NOT NULL DEFAULT '0' COMMENT '评论次数',
    `is_published` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否发布: 1-已发布, 0-草稿',
    `is_top` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否置顶: 1-是, 0-否',
    `is_featured` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否推荐: 1-是, 0-否',
    `published_at` timestamp NULL DEFAULT NULL COMMENT '发布时间',
    `tags` varchar(500) DEFAULT NULL COMMENT '标签(逗号分隔)',
    `seo_title` varchar(200) DEFAULT NULL COMMENT 'SEO标题',
    `seo_keywords` varchar(500) DEFAULT NULL COMMENT 'SEO关键词',
    `seo_description` varchar(500) DEFAULT NULL COMMENT 'SEO描述',
    `status` tinyint(1) NOT NULL DEFAULT '1' COMMENT '状态: 1-正常, 0-禁用',
    `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_slug` (`slug`),
    KEY `idx_author_id` (`author_id`),
    KEY `idx_category_id` (`category_id`),
    KEY `idx_is_published` (`is_published`),
    KEY `idx_is_top` (`is_top`),
    KEY `idx_status` (`status`),
    KEY `idx_created_at` (`created_at`),
    KEY `idx_published_at` (`published_at`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='文章表';

-- 插入测试数据
INSERT INTO `articles` (`title`, `slug`, `author_id`, `author_name`, `summary`, `content`, `cover_image`, `is_published`, `status`) VALUES
('欢迎使用 GinForge 框架', 'welcome-to-ginforge', 1, '超级管理员', '这是一篇测试文章，介绍 GinForge 微服务框架的主要特性。', '# 欢迎使用 GinForge\n\nGinForge 是一个功能完整的企业级微服务开发框架。\n\n## 主要特性\n\n- 微服务架构\n- RBAC 权限管理\n- 代码生成器\n- 完整的后台管理\n\n开始您的开发之旅吧！', '/images/welcome.jpg', 1, 1),
('Go 语言最佳实践', 'golang-best-practices', 1, '超级管理员', '分享 Go 语言开发中的最佳实践和常见陷阱。', '# Go 语言最佳实践\n\n本文总结了 Go 开发中的一些最佳实践...\n\n## 并发编程\n\n- 使用 channel 进行通信\n- 避免共享内存\n- 正确使用 context', '/images/golang.jpg', 1, 1),
('微服务架构设计', 'microservices-architecture', 1, '超级管理员', '深入探讨微服务架构的设计原则和实践经验。', '# 微服务架构设计\n\n微服务是一种架构风格...\n\n## 核心原则\n\n1. 单一职责\n2. 服务自治\n3. 去中心化', '/images/microservices.jpg', 0, 1);

