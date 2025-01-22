package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"music/internal/models"
	"reflect"
	"strings"
)

type Music interface {
	GetMusics(ctx context.Context, group, song, releaseDate, text string, limit, offset int) ([]models.MusicInfo, error)
	GetSongLyricsByVerses(ctx context.Context, ID, couplet, size int) ([]string, error)
	AddMusic(ctx context.Context, music models.MusicInfo) error
	UpdateMusic(ctx context.Context, ID int, updates models.MusicUpdate) error
	DeleteMusic(ctx context.Context, ID int) error
}

type musicPostgres struct {
	db *sql.DB
}

func newMusicPostgres(db *sql.DB) Music {
	return &musicPostgres{db: db}
}

func (r *musicPostgres) GetMusics(
	ctx context.Context, group, song, releaseDate, text string, limit, offset int,
) ([]models.MusicInfo, error) {
	query := `
		SELECT id, music_group, song, release_date, text, link 
		FROM musics 
		WHERE (COALESCE($1, '') = '' OR music_group ILIKE '%' || $1 || '%') 
		  AND (COALESCE($2, '') = '' OR song ILIKE '%' || $2 || '%') 
		  AND (COALESCE($3, '') = '' OR release_date::TEXT ILIKE '%' || $3 || '%') 
		  AND (COALESCE($4, '') = '' OR text ILIKE '%' || $4 || '%') 
		LIMIT $5 
		OFFSET $6;
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}

	rows, err := stmt.QueryContext(ctx, group, song, releaseDate, text, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var result []models.MusicInfo

	for rows.Next() {
		var music models.MusicInfo

		if err := rows.Scan(
			&music.ID,
			&music.Group,
			&music.Song,
			&music.RelaseDate,
			&music.Text,
			&music.Link,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		result = append(result, music)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func (r *musicPostgres) GetSongLyricsByVerses(ctx context.Context, ID, couplet, size int) ([]string, error) {
	query := "SELECT text FROM musics WHERE id = $1;"

	var text string

	err := r.db.QueryRowContext(ctx, query, ID).Scan(&text)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("song not found")
		}

		return nil, fmt.Errorf("failed to fetch song lyrics: %w", err)
	}

	verses := strings.Split(text, "\\n\\n")

	start := (couplet - 1)

	end := start + size
	if end > len(verses) {
		end = len(verses)
	}

	return verses[start:end], nil
}

func (r *musicPostgres) AddMusic(ctx context.Context, music models.MusicInfo) error {
	query := `
		INSERT INTO musics (music_group, song, release_date, text, link) 
		VALUES ($1, $2, $3, $4, $5);
	`

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(
		ctx,
		query,
		music.Group,
		music.Song,
		music.RelaseDate,
		music.Text,
		music.Link,
	); err != nil {
		return err
	}

	return nil
}

func (r *musicPostgres) UpdateMusic(ctx context.Context, ID int, updates models.MusicUpdate) error {
	if reflect.ValueOf(updates).IsZero() {
		return fmt.Errorf("No updates provided")
	}

	query := "UPDATE musics SET "
	var args []any
	argCounter := 1

	val := reflect.ValueOf(updates)
	typ := reflect.TypeOf(updates)

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		value := val.Field(i)

		if value.IsZero() {
			continue
		}

		columnName := field.Tag.Get("db")
		if columnName == "" {
			columnName = field.Tag.Get("json")
		}

		query += fmt.Sprintf("%s = $%d, ", columnName, argCounter)
		args = append(args, value.Interface())
		argCounter++
	}

	query = strings.TrimSuffix(query, ", ")
	query += fmt.Sprintf(" WHERE id = $%d", argCounter)
	args = append(args, ID)

	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	if _, err := stmt.ExecContext(ctx, args...); err != nil {
		return fmt.Errorf("Failed to update music: %w", err)
	}

	return nil

}

func (r *musicPostgres) DeleteMusic(ctx context.Context, ID int) error {
	query := "DELETE FROM music WHERE id = $1;"

	if _, err := r.db.ExecContext(ctx, query, ID); err != nil {
		return err
	}

	return nil
}
