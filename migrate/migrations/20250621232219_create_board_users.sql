-- +goose Up
-- +goose StatementBegin
CREATE TABLE board_users (
  user_id INT NOT NULL,
  board_id INT NOT NULL,
  PRIMARY KEY (user_id, board_id),
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (board_id) REFERENCES boards(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE board_users;
-- +goose StatementEnd
