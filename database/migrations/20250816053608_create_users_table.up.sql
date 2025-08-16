CREATE TABLE users (
    id CHAR(26) PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    phone VARCHAR(20) NULL,
    province VARCHAR(100) NULL,
    district VARCHAR(100) NULL,
    address VARCHAR(255) NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'customer',
    social_provider VARCHAR(50) NULL,
    social_provider_id VARCHAR(100) NULL,
    last_sign_in_at TIMESTAMP WITH TIME ZONE NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE NULL
);