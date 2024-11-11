-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS settings (
  id UUID DEFAULT gen_random_uuid() NOT NULL,
  key VARCHAR NOT NULL,
  value TEXT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
  deleted_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS settings;

-- +goose StatementEnd
