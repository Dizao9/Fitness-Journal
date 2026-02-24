-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercises DROP CONSTRAINT IF EXISTS exercises_name_key;
CREATE UNIQUE INDEX idx_user_exercise ON exercises (athlete_id, name) WHERE athlete_id IS NOT NULL;

CREATE UNIQUE INDEX idx_system_exercise ON exercises (name) WHERE athlete_id IS NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_user_exercise;
DROP INDEX IF EXISTS idx_system_exercise;
ALTER TABLE exercises ADD CONSTRAINT exercises_name_key UNIQUE (name);
-- +goose StatementEnd
