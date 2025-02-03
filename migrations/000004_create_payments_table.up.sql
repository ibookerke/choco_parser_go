CREATE TABLE payments (
    id BIGINT PRIMARY KEY,
    created_by BIGINT,
    type VARCHAR(50),
    amount BIGINT,
    discount_amount BIGINT,
    created_at TIMESTAMP,

    location_title TEXT,
    location_partner_id UUID
);
