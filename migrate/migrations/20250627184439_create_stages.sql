-- +goose Up
-- +goose StatementBegin
CREATE TABLE stages (
  id SERIAL PRIMARY KEY,
  title VARCHAR NOT NULL,
  position INT NOT NULL,
  wip_limit INT,
  stage_type INT NOT NULL,
  board_id INT NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE stages;
-- +goose StatementEnd
