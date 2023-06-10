#每日指标
create table daily_basic (
                       id bigint unsigned not null AUTO_INCREMENT comment '主键',
                       ts_code varchar(64) default '' comment 'TS代码',
                       trade_date varchar(64) default '' comment '交易日期',
                       `close` decimal(18,3) default 0 comment '当日收盘价',
                       turnover_rate decimal(18,5) default 0 comment '换手率 %',
                       turnover_rate_f decimal(18,5) default 0 comment '换手率 自由流通股',
                       volume_ratio decimal(12,3) default 0 comment '量比',
                       pe decimal(18,5) default 0 comment '市盈率（总市值/净利润， 亏损的PE为空）',
                       pe_ttm decimal(18,5) default 0 comment '市盈率（TTM，亏损的PE为空）',
                       pb decimal(18,5) default 0 comment '市净率（总市值/净资产）',
                       ps decimal(18,5) default 0 comment '市销率',
                       ps_ttm decimal(18,5) default 0 comment '市销率 TTM',
                       dv_ratio decimal(18,5) default  comment '股息率 （%）',
                       dv_ttm decimal(18,5) default 0 comment '股息率（TTM）（%）',
                       total_share decimal(18,5) default 0 comment '总股本 万股',
                       float_share decimal(18,5) default 0 comment '流通股本 万股',
                       free_share decimal(18,5) default 0 comment '自由流通股本 万',
                       `total_mv` decimal(18,5) default 0 comment '总市值 （万元）',
                       circ_mv decimal(18,5) default 0 comment '流通市值（万元）'


                       created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
                       updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
                       deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

                       PRIMARY KEY(`id`),
                       UNIQUE KEY `uk_tscode_tradedate`(`ts_code`, `trade_date`),
                       KEY `idx_tradedate`(`trade_date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='每日指标';