#股票列表
create table stock (
                            id bigint unsigned not null AUTO_INCREMENT comment '主键',
                            ts_code varchar(64) default '' comment 'TS代码',
                            symbol varchar(64) default '' comment '股票代码',
                            `name` varchar(128) default '' comment '股票名称',
                            area varchar(128) default '' comment '地域',
                            industry varchar(256) default '' comment '所属行业',
                            fullname varchar(256) default '' comment '股票全称',
                            enname varchar(256) default '' comment '英文全称',
                            cnspell varchar(128) default '' comment '拼音缩写',
                            market varchar(64) default '' comment '市场类型',
                            exchange varchar(64) default '' comment '交易所代码',
                            curr_type varchar(128) default '' comment '交易货币',
                            list_status varchar(64) default '' comment '上市状态',
                            list_date varchar(128) default '' comment '上市日期',
                            delist_date varchar(128) default '' comment '退市日期',
                            is_hs varchar(64) default '' comment '是否沪深港通标的',

                            created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
                            updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
                            deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

                            PRIMARY KEY(`id`),
                            UNIQUE KEY `uk_tscode`(`ts_code`),
                            KEY `idx_name`(`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='股票列表';
