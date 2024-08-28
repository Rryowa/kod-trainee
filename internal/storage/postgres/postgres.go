package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"kod/internal/models"
	"kod/internal/models/config"
	"kod/internal/storage"
	"kod/internal/util"
)

type Database struct {
	Pool      *pgxpool.Pool
	zapLogger *zap.SugaredLogger
}

func (d *Database) AddUser(ctx context.Context, user *models.User) (models.User, error) {
	const op = "storage.AddUser"
	query := `INSERT INTO users (username, password)
				VALUES ($1, $2) returning id, username, password`

	rows, err := d.Pool.Query(ctx, query, user.Username, user.Password)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newUser models.User
	err = pgxscan.ScanOne(&newUser, rows)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newUser, nil
}

func (d *Database) GetUser(ctx context.Context, userId int) (models.User, error) {
	const op = "storage.GetUser"
	query := `SELECT id, username, password FROM users
				WHERE id = $1`

	rows, err := d.Pool.Query(ctx, query, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newUser models.User
	err = pgxscan.ScanOne(&newUser, rows)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op2, err)
	}

	return newUser, nil
}

func (d *Database) AddNote(ctx context.Context, note models.Note) (models.Note, error) {
	const op = "storage.AddNote"
	query := `INSERT INTO notes (user_id, text, created_at)
				VALUES ($1, $2, $3) returning id, user_id, text, created_at`

	rows, err := d.Pool.Query(ctx, query, note.Id, note.UserId, note.Text, note.CreatedAt)
	if err != nil {
		return models.Note{}, fmt.Errorf("%s: %w", op, err)
	}

	const op2 = op + "pgxscan"
	var newNote models.Note
	err = pgxscan.ScanOne(&newNote, rows)
	if err != nil {
		return models.Note{}, fmt.Errorf("%s: %w", op2, err)
	}

	//if err = d.cacheService.Put(note.Id, newNote); err != nil {
	//	return models.Note{}, err
	//}

	return newNote, nil
}

func (d *Database) GetNotes(ctx context.Context, userId, offset, limit int) ([]models.Note, error) {
	const op = "storage.GetNotes"
	query := `SELECT id, user_id, text, created_at
				FROM notes
				WHERE user_id=$1
				ORDER BY created_at DESC
				OFFSET $2
				FETCH NEXT $3 ROWS ONLY`

	rows, err := d.Pool.Query(ctx, query, userId, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	const op2 = op + "pgxscan"
	var notes []models.Note
	if err := pgxscan.ScanAll(&notes, rows); err != nil {
		return nil, fmt.Errorf("%s: %w", op2, err)
	}
	return notes, err
}

//// Exists returns true if note exists in cache or repository
//func (d *Database) Exists(ctx context.Context, id string) (models.Note, bool) {
//
//	//if note, ok := d.cacheService.Get(id); ok {
//	//	return note, true
//	//}
//	if note, err := d.get(ctx, id); err == nil {
//		return note, true
//	}
//	return models.Note{}, false
//}
//
//func (d *Database) get(ctx context.Context, userId string) (models.Note, error) {
//	const op = "storage.Get"
//	var note models.Note
//	query := `SELECT id, user_id, text, created_at
//				FROM notes
//				WHERE user_id=$1`
//
//	if err := pgxscan.Get(ctx, d.Pool, &note, query, userId); err != nil {
//		if pgxscan.NotFound(err) {
//			return models.Note{}, util.ErrNoteNotFound
//		} else {
//			return models.Note{}, fmt.Errorf("%s: %w", op, err)
//		}
//	}
//	return note, nil
//}

func NewPostgresRepository(ctx context.Context, cfg *config.DbConfig, zap *zap.SugaredLogger) storage.Storage {
	connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	var pool *pgxpool.Pool
	var err error

	err = util.DoWithTries(func() error {
		ctxTimeout, cancel := context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()

		pool, err = pgxpool.New(ctxTimeout, connStr)
		if err != nil {
			zap.Fatalln(err, "db connection error")
		}

		return nil
	}, cfg.Attempts, cfg.Timeout)

	if err != nil {
		zap.Fatalln(err, "DoWithTries error")
	}
	zap.Infoln("Connected to db")

	return &Database{
		Pool:      pool,
		zapLogger: zap,
	}
}