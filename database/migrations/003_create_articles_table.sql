-- 文章表
CREATE TABLE IF NOT EXISTS `articles` (
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

-- 插入测试数据
INSERT INTO `articles` (`title`, `content`, `summary`, `author_id`, `author_name`, `status`, `is_top`, `view_count`, `created_at`, `updated_at`)
VALUES 
('欢迎使用 GinForge 框架', '# GinForge 框架介绍\n\nGinForge 是一个强大的 Go Web 开发框架...', '快速了解 GinForge 框架的核心特性', 1, 'Admin', 1, 1, 100, NOW(), NOW()),
('如何使用代码生成器', '# 代码生成器使用指南\n\n代码生成器可以帮助您快速生成 CRUD 代码...', '5分钟学会使用代码生成器', 1, 'Admin', 1, 0, 50, NOW(), NOW()),
('GinForge 最佳实践', '# 最佳实践\n\n在使用 GinForge 开发时，建议遵循以下最佳实践...', '提升开发效率的最佳实践', 1, 'Admin', 0, 0, 20, NOW(), NOW());
