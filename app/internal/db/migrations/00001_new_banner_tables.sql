-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner (
	id SERIAL PRIMARY KEY,
	feature_id INT NOT NULL,
  	is_active BOOL NOT NULL,
	content jsonb NOT NULL,
  	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS banner_relation (
	banner_id INT NOT NULL REFERENCES banner ON DELETE CASCADE,
  	feature_id INT NOT NULL,
  	tag_id INT NOT NULL,
  	UNIQUE(feature_id, tag_id)
);

CREATE INDEX IF NOT EXISTS banner_id_relation ON banner_relation (banner_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS banner_id_relation;
DROP TABLE IF EXISTS banner;
DROP TABLE IF EXISTS banner_relation;
-- +goose StatementEnd
