CREATE TABLE skuder_user (
    login TEXT PRIMARY KEY,
    password_hash BYTEA NOT NULL,
    role TEXT NOT NULL
);

INSERT INTO skuder_user (login, password_hash, role) VALUES
    -- password: Cuzee2motof6aiJe
    ('admin', '$2a$10$6oSCRf0RN3l2qs1zaX0Qie11eMcMaljgGml4VsBlz3pucxzTwniTO', 'admin')

CREATE TABLE person (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    position TEXT NOT NULL,
    unit TEXT NOT NULL
);

CREATE TABLE person_face (
    id BIGSERIAL PRIMARY KEY,
    person_id BIGINT NOT NULL REFERENCES person (id) ON DELETE CASCADE,
    descriptor REAL[128] NOT NULL,
    photo_id TEXT NOT NULL
);

CREATE TABLE event (
    id BIGSERIAL PRIMARY KEY,
    time TIMESTAMP WITH TIME ZONE NOT NULL,
    type TEXT NOT NULL,
    data JSONB
);
