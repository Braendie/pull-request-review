CREATE TABLE team (
    team_id BIGSERIAL PRIMARY KEY,
    team_name VARCHAR(50)
);

CREATE TABLE users (
    user_id VARCHAR(50) PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    team_id BIGSERIAL NOT NULL REFERENCES team(team_id) ON DELETE CASCADE,
    is_active BOOLEAN NOT NULL
);

CREATE TABLE pull_requests (
    pull_request_id VARCHAR(50) PRIMARY KEY,
    pull_request_name VARCHAR(50) NOT NULL,
    author_id VARCHAR(50) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    status VARCHAR(6) NOT NULL,
    created_at TIMESTAMPTZ, 
    merged_at TIMESTAMPTZ
);

CREATE TABLE pull_requests_users (
    user_id VARCHAR(50) NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    pull_request_id VARCHAR(50) NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, pull_request_id)
);