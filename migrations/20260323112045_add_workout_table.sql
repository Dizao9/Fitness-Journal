-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts ADD COLUMN status TEXT DEFAULT 'finished';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts DROP COLUMN status;
-- +goose StatementEnd
