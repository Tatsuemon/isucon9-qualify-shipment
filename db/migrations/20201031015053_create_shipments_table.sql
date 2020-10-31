
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE `shipments` (
    `id` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '予約ID',
    `to_address` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '配達先の住所',
    `from_address` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '配達元の住所',
    `to_address` BIGINT COLLATE utf8mb4_bin NOT NULL COMMENT '配達先の住所',
    `to_name` BIGINT COLLATE utf8mb4_bin NOT NULL COMMENT '配達先の名前',
    `from_address` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '配達元の住所',
    `from_name` varchar(255) COLLATE utf8mb4_bin NOT NULL COMMENT '配達元の名前',
    `status` varchar(127) COLLATE utf8mb4_bin NOT NULL COMMENT '配達ステータス',
    `reserve_time`DATETIME COLLATE utf8mb4_bin NOT NULL COMMENT '配達予約時間',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE `shipments`;
