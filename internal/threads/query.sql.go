// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package threads

import (
	"context"
	"database/sql"
)

const createThread = `-- name: CreateThread :one
INSERT INTO threads (content, thread_id, user_id, created_at, updated_at)
VALUES (?, ?, ?, ?, ?)
RETURNING id, content, thread_id, user_id, created_at, updated_at, "foreign"
`

type CreateThreadParams struct {
	Content   string
	ThreadID  sql.NullInt64
	UserID    int64
	CreatedAt string
	UpdatedAt string
}

func (q *Queries) CreateThread(ctx context.Context, arg CreateThreadParams) (Thread, error) {
	row := q.db.QueryRowContext(ctx, createThread,
		arg.Content,
		arg.ThreadID,
		arg.UserID,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i Thread
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.ThreadID,
		&i.UserID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Foreign,
	)
	return i, err
}

const getAllThreads = `-- name: GetAllThreads :many
SELECT threads.id, content, username, avatar
FROM threads
JOIN profiles ON profiles.user_id = threads.user_id
ORDER BY threads.created_at desc
`

type GetAllThreadsRow struct {
	ID       int64
	Content  string
	Username string
	Avatar   sql.NullString
}

func (q *Queries) GetAllThreads(ctx context.Context) ([]GetAllThreadsRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllThreads)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllThreadsRow
	for rows.Next() {
		var i GetAllThreadsRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.Username,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThreadByID = `-- name: GetThreadByID :one
SELECT threads.id, content, username, avatar
FROM threads
JOIN profiles ON profiles.user_id = threads.user_id
WHERE threads.id = ?
`

type GetThreadByIDRow struct {
	ID       int64
	Content  string
	Username string
	Avatar   sql.NullString
}

func (q *Queries) GetThreadByID(ctx context.Context, id int64) (GetThreadByIDRow, error) {
	row := q.db.QueryRowContext(ctx, getThreadByID, id)
	var i GetThreadByIDRow
	err := row.Scan(
		&i.ID,
		&i.Content,
		&i.Username,
		&i.Avatar,
	)
	return i, err
}

const getUserThreads = `-- name: GetUserThreads :many
SELECT threads.id, content, username, avatar
FROM threads
JOIN profiles ON profiles.user_id = threads.user_id
WHERE threads.user_id = ?
ORDER BY threads.created_at desc
`

type GetUserThreadsRow struct {
	ID       int64
	Content  string
	Username string
	Avatar   sql.NullString
}

func (q *Queries) GetUserThreads(ctx context.Context, userID int64) ([]GetUserThreadsRow, error) {
	rows, err := q.db.QueryContext(ctx, getUserThreads, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUserThreadsRow
	for rows.Next() {
		var i GetUserThreadsRow
		if err := rows.Scan(
			&i.ID,
			&i.Content,
			&i.Username,
			&i.Avatar,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
