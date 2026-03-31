-- 与 game/model/model.go 对齐的基准表结构（新库可整文件执行）。
-- 表结构变更请改 model 后由维护者提供增量 ALTER，勿依赖程序内 AutoMigrate。

SET NAMES utf8mb4;

CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `phone` varchar(20) NOT NULL COMMENT '手机号',
  `password` varchar(255) NOT NULL COMMENT '密码哈希',
  `username` varchar(50) NOT NULL COMMENT '昵称',
  `avatar` varchar(500) NOT NULL DEFAULT '' COMMENT '头像URL或emoji',
  `source` varchar(50) NOT NULL DEFAULT '' COMMENT '注册来源游戏标识',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

CREATE TABLE IF NOT EXISTS `scores` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `game_key` varchar(50) NOT NULL DEFAULT 'match3' COMMENT '游戏标识',
  `completion_time_ms` int NOT NULL COMMENT '完成用时毫秒',
  `user_agent` text COMMENT '上报时User-Agent',
  `ip` varchar(45) COMMENT '上报时客户端IP',
  `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_scores_user_id` (`user_id`),
  KEY `idx_scores_game_key` (`game_key`),
  KEY `idx_scores_completion_time_ms` (`completion_time_ms`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户游戏成绩记录表';

CREATE TABLE IF NOT EXISTS `games` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `game_key` varchar(50) NOT NULL COMMENT '游戏唯一标识',
  `name` varchar(100) NOT NULL COMMENT '展示名称',
  `icon` varchar(50) DEFAULT NULL COMMENT '封面图标或emoji',
  `desc` varchar(500) DEFAULT NULL COMMENT '简介',
  `url` varchar(200) DEFAULT NULL COMMENT '入口相对路径',
  `status` varchar(20) NOT NULL DEFAULT 'online' COMMENT '状态 online/coming/offline',
  `sort_order` int NOT NULL DEFAULT 0 COMMENT '排序权重升序',
  `created_at` datetime(3) DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime(3) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_games_game_key` (`game_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='游戏大厅配置表';

CREATE TABLE IF NOT EXISTS `game_runs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `run_id` char(36) NOT NULL COMMENT '对局凭证ID',
  `user_id` bigint unsigned NOT NULL COMMENT '用户ID',
  `game_key` varchar(50) NOT NULL COMMENT '游戏标识',
  `started_at` datetime(3) NOT NULL COMMENT '对局开始时间',
  `expires_at` datetime(3) NOT NULL COMMENT '凭证过期时间',
  `used_at` datetime(3) DEFAULT NULL COMMENT '成绩提交时间',
  `created_at` datetime(3) DEFAULT NULL COMMENT '记录创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_game_runs_run_id` (`run_id`),
  KEY `idx_game_run_user_game` (`user_id`,`game_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='对局凭证表';

CREATE TABLE IF NOT EXISTS `api_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键',
  `method` varchar(10) NOT NULL COMMENT 'HTTP方法',
  `path` varchar(500) NOT NULL COMMENT '请求路径',
  `query` text COMMENT '查询串',
  `req_body` text COMMENT '请求体截断',
  `status_code` int NOT NULL COMMENT '响应状态码',
  `resp_body` text COMMENT '响应体截断',
  `user_id` bigint unsigned NOT NULL DEFAULT 0 COMMENT '已登录用户ID，未登录为0',
  `ip` varchar(45) COMMENT '客户端IP',
  `user_agent` text COMMENT 'User-Agent',
  `latency_ms` bigint NOT NULL COMMENT '处理耗时毫秒',
  `created_at` datetime(3) DEFAULT NULL COMMENT '记录时间',
  PRIMARY KEY (`id`),
  KEY `idx_api_logs_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='API访问日志表';
