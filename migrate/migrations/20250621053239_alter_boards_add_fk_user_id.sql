-- +goose Up
-- +goose StatementBegin
ALTER TABLE boards
ADD CONSTRAINT fk_boards_user_id
FOREIGN KEY (user_id) REFERENCES users (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE boards
DROP CONSTRAINT fk_boards_user_id;
-- +goose StatementEnd
