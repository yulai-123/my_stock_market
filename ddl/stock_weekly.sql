# 周线行情
create table stock_weekly (
                              id bigint unsigned not null AUTO_INCREMENT comment '主键',
                              ts_code varchar(64) default '' comment 'TS代码',
                              trade_date varchar(64) default '' comment '交易日期',
                              `open` decimal(12,3) default 0 comment '开盘价',
                              high decimal(12,3) default 0 comment '最高价',
                              low decimal(12,3) default 0 comment '最低价',
                              `close` decimal(12,3) default 0 comment '收盘价',
                              pre_close decimal(12,3) default 0 comment '昨收价',
                              `change` decimal(12,3) default 0 comment '涨跌额',
                              pct_chg decimal(12,3) default 0 comment '涨跌幅',
                              vol decimal(18,3) default 0 comment '成交量',
                              amount decimal(18,3) default 0 comment '成交量',

                              created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
                              updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
                              deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

                              PRIMARY KEY(`id`),
                              UNIQUE KEY `uk_tscode_tradedate`(`ts_code`, trade_date),
                              KEY `idx_tradedate`(`trade_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='周线行情';
