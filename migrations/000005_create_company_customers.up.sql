CREATE TABLE company_customers (
    id SERIAL PRIMARY KEY,
    company TEXT,
    user_id BIGINT,
    full_name TEXT,
    phone VARCHAR(50),
    turnover FLOAT,
    last_visit_date TIMESTAMP,
    visits_count BIGINT,
    average_bill FLOAT
);
