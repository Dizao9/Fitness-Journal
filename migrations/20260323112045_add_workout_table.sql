-- +goose Up
-- +goose StatementBegin
ALTER TABLE workouts ADD COLUMN status TEXT DEFAULT 'completed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE workouts DROP COLUMN workouts;
-- +goose StatementEnd
