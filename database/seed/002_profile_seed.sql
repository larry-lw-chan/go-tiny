-- +goose Up
-- +goose StatementBegin
INSERT INTO profiles (name, bio, link, avatar, private, user_id, created_at, updated_at) 
VALUES('Testy One', 'I am a test user', 'www.test.com', null, 0, 1, UNIXEPOCH(), UNIXEPOCH());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM profiles WHERE id = 1;
-- +goose StatementEnd
