CREATE TABLE branches (
    id BIGINT PRIMARY KEY,
    name VARCHAR(255) NULL,
    status VARCHAR(50) NULL,
    type_id INT NULL,
    type_name VARCHAR(50) NULL,
    type_description TEXT NULL,
    token VARCHAR(255) NULL,
    location_id UUID NULL,
    location_name VARCHAR(255) NULL,
    partner_id UUID NULL,
    partner_name VARCHAR(255) NULL,
    partner_logo TEXT NULL
);