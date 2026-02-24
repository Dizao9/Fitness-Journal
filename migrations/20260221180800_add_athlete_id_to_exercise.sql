-- +goose Up
-- +goose StatementBegin
ALTER TABLE exercises
ADD COLUMN athlete_id UUID
REFERENCES athletes(id)
ON DELETE CASCADE
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE exercises
DROP COLUMN athlete_id;
-- +goose StatementEnd
