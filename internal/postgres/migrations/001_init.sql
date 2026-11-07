CREATE TABLE IF NOT EXISTS tasks(id BIGSERIAL PRIMARY KEY,
                            title TEXT NOT NULL,
                            description TEXT NOT NULL,
                            done BOOLEAN DEFAULT FALSE,
                            created_at  TIMESTAMPTZ NOT NULL,
                            updated_at  TIMESTAMPTZ NOT NULL);

