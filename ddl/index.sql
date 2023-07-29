# 指数
create table `index` (
                       id bigint unsigned not null AUTO_INCREMENT comment '主键',
                       ts_code varchar(64) default '' comment 'TS代码',
                       `name` varchar(128) default '' comment '指数全称',
                       fullname varchar(256) default '' comment '股票全称',
                       market varchar(64) default '' comment '市场类型',
                       list_date varchar(128) default '' comment '发布日期',

                       publisher varchar(128) default  '' comment '发布方',
                       index_type varchar(64) default  '' comment '指数风格',
                       category varchar(64) default '' comment '指数类别',
                       base_date varchar(64) default '' comment '基期',
                       base_point decimal(18,3) default 0 comment '基点',
                       weight_rule varchar(64) default '' comment '加权方式',
                       `desc` varchar(2048) default '' comment '描述',
                       exp_date varchar(64) default '' comment '终止日期',

                       created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
                       updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
                       deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

                       PRIMARY KEY(`id`),
                       UNIQUE KEY `uk_tscode`(`ts_code`),
                       KEY `idx_name`(`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='股票列表';