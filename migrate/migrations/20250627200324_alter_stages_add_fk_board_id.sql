-- +goose Up
-- +goose StatementBegin
ALTER TABLE stages
ADD CONSTRAINT fk_stages_board_id
FOREIGN KEY (board_id) REFERENCES boards (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE boards
DROP CONSTRAINT fk_stages_board_id;
-- +goose StatementEnd
