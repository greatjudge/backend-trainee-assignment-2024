-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS banner (
	id SERIAL PRIMARY KEY,
  	is_active BOOL,
  	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS banner_relation (
	banner_id INT REFERENCES banner,
  	feature_id INT NOT NULL,
  	tag_id INT NOT NULL,
  	UNIQUE(feature_id, tag_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS banner;
DROP TABLE IF EXISTS banner_relation;
-- +goose StatementEnd
