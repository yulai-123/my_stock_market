#现金流量表
create table income (
                        id bigint unsigned not null AUTO_INCREMENT comment '主键',

                        ts_code varchar(64) default '' comment 'TS代码',
                        ann_date varchar(64) default '' comment '公告日期',
                        f_ann_date varchar(64) default '' comment '实际公告日期',
                        end_date varchar(64) default '' comment '报告期',
                        report_type varchar(64) default '' comment '报告类型： 参考下表说明',
                        comp_type varchar(64) default '' comment '公司类型：1一般工商业 2银行 3保险 4证券',
                        end_type varchar(64) default '' comment '报告期：1第一季度 2第二季度 3第三季度 4第四季度',

                        n_cashflow_act decimal(20,4) default 0 comment '经营活动产生的现金流量净额',

                        created_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '创建时间',
                        updated_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '更新时间',
                        deleted_at bigint unsigned NOT NULL DEFAULT 0 COMMENT '删除时间',

                        PRIMARY KEY(`id`),
                        UNIQUE KEY `uk_tscode`(`ts_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='现金流量表';