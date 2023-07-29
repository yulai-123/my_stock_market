package ddl

create table `fund_basic` (
    id bigint unsigned not null AUTO_INCREMENT comment '主键',
    ts_code varchar(64) default '' comment 'TS代码',
    `name` varchar(128) default '' comment '全称',
    `management` varchar(128) default '' comment '管理人',
    `custodian` varchar(128) default '' comment '托管人',
    `fund_type` varchar(128) default '' comment '投资类型',
    `found_date` varchar(128) default '' comment '投资类型',
    `due_date` varchar(128) default '' comment '到期日期',
    list_date varchar(128) default '' comment '发布日期',
    `issue_date` varchar(128) default '' comment '发行日期',
    `delist_date` varchar(128) default '' comment '退市日期',
    `issue_amount` decimal(18,3) default 0 comment '发行份额',
    `m_fee` decimal(18,3) default 0 comment '管理费',
    `c_fee` decimal(18,3) default 0 comment '托管费',
    `duration_year` decimal(18,3) default 0 comment '存续期',
    `p_value` decimal(18,3) default 0 comment '面值',
    `min_amount` decimal(18,3) default 0 comment '起点金额',
    `exp_return` decimal(18,3) default 0 comment '预期收益率',
    `benchmark` varchar(128) default '' comment '业绩比较基准',
    `status` varchar(64) default '' comment '存续状态D摘牌 I发行 L已上市',
    `invest_type` varchar(64) default '' comment '投资风格',
    `type` varchar(64) default '' comment '基金类型',
    `trustee` varchar(64) default '' comment '委托人',
    `purc_startdate` varchar(64) default '' comment '日常申购起始日',
    `redm_startdate` varchar(64) default '' comment '日常赎回起始日',
    `market` varchar(64) default '' comment 'E场内O场外',

    created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
    updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
    deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

    PRIMARY KEY(`id`),
    UNIQUE KEY `uk_tscode`(`ts_code`),
    KEY `idx_name`(`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='基金列表';