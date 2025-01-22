-- +goose Up
-- +goose StatementBegin
ALTER TABLE property_expenses
ALTER COLUMN type SET DATA TYPE varchar(500);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE property_expenses
ALTER COLUMN type SET DATA TYPE varchar(100);
-- +goose StatementEnd
