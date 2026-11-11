package repository

import (
	"context"
	"grades/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

// ================= CREATE =================

func (r *UserRepository) Create(ctx context.Context, email, password string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1,$2)`,
		email, password)

	return err
}

// ================= LOGIN =================

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (id int, hash string, err error) {
	err = r.db.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE email=$1`,
		email,
	).Scan(&id, &hash)

	return
}

// ================= GET BY ID =================

func (r *UserRepository) GetByID(ctx context.Context, id int) (*model.User, error) {
	var user model.User

	err := r.db.QueryRow(ctx, `
		SELECT id, email,name, age, grade
		FROM users WHERE id=$1`, id).
		Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Age,
			&user.Grade,
		)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// ================= UPDATE =================

func (r *UserRepository) Update(ctx context.Context, id int, name string, age int, grade string) error {

	_, err := r.db.Exec(ctx,
		`UPDATE users SET
			name  = COALESCE(NULLIF($1, ''), name),
			age   = CASE WHEN $2 = 0 THEN age ELSE $2 END,
			grade = COALESCE(NULLIF($3, ''), grade)
		WHERE id = $4;`,
		name, age, grade, id)
	return err
}

// ================= DELETE =================

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM users WHERE id=$1`, id)
	return err
}

// Create
func (r *UserRepository) CreateSubj(ctx context.Context, subject model.Subjects) error {

	query := `
		INSERT INTO subjects (
			user_id,
			alg,
			bio,
			whis,
			geo,
			lang,
			it,
			khis,
			kaz,
			rights,
			lit,
			rus,
			phys,
			chem
		)
		VALUES (
			$1, $2, $3, $4, $5, $6, $7,
			$8, $9, $10, $11, $12, $13, $14
		)
		RETURNING
			id,
			user_id,
			alg,
			bio,
			whis,
			geo,
			lang,
			it,
			khis,
			kaz,
			rights,
			lit,
			rus,
			phys,
			chem;
	`

	var result model.Subjects

	err := r.db.QueryRow(
		ctx,
		query,
		subject.UserId,
		subject.Alg,
		subject.Bio,
		subject.Whis,
		subject.Geo,
		subject.Lang,
		subject.It,
		subject.Khis,
		subject.Kaz,
		subject.Rights,
		subject.Lit,
		subject.Rus,
		subject.Phys,
		subject.Chem,
	).Scan(
		&result.ID,
		&result.UserId,
		&result.Alg,
		&result.Bio,
		&result.Whis,
		&result.Geo,
		&result.Lang,
		&result.It,
		&result.Khis,
		&result.Kaz,
		&result.Rights,
		&result.Lit,
		&result.Rus,
		&result.Phys,
		&result.Chem,
	)

	if err != nil {
		return err
	}

	return nil
}

// GET
func (r *UserRepository) GetSubj(
	ctx context.Context,
	userID int,
) (*model.Subjects, error) {

	query := `
		SELECT
			id,
			user_id,
			alg,
			bio,
			whis,
			geo,
			lang,
			it,
			khis,
			kaz,
			rights,
			lit,
			rus,
			phys,
			chem
		FROM subjects
		WHERE user_id = $1
	`

	var subject model.Subjects

	err := r.db.QueryRow(
		ctx,
		query,
		userID,
	).Scan(
		&subject.ID,
		&subject.UserId,
		&subject.Alg,
		&subject.Bio,
		&subject.Whis,
		&subject.Geo,
		&subject.Lang,
		&subject.It,
		&subject.Khis,
		&subject.Kaz,
		&subject.Rights,
		&subject.Lit,
		&subject.Rus,
		&subject.Phys,
		&subject.Chem,
	)

	if err != nil {
		return nil, err
	}

	return &subject, nil
}

// Update
func (r *UserRepository) UpdateSubj(
	ctx context.Context,
	subject model.Subjects,
) error {

	query := `
		UPDATE subjects
		SET
			alg = COALESCE($1, alg),
			bio = COALESCE($2, bio),
			whis = COALESCE($3, whis),
			geo = COALESCE($4, geo),
			lang = COALESCE($5, lang),
			it = COALESCE($6, it),
			khis = COALESCE($7, khis),
			kaz = COALESCE($8, kaz),
			rights = COALESCE($9, rights),
			lit = COALESCE($10, lit),
			rus = COALESCE($11, rus),
			phys = COALESCE($12, phys),
			chem = COALESCE($13, chem)
		WHERE user_id = $14

	`

	_, err := r.db.Exec(
		ctx,
		query,
		subject.Alg,
		subject.Bio,
		subject.Whis,
		subject.Geo,
		subject.Lang,
		subject.It,
		subject.Khis,
		subject.Kaz,
		subject.Rights,
		subject.Lit,
		subject.Rus,
		subject.Phys,
		subject.Chem,
		subject.UserId,
	)

	return err
}

// DELETE
func (r *UserRepository) DeleteSubj(
	ctx context.Context,
	userID int,
) error {

	query := `
		DELETE FROM subjects
		WHERE user_id = $1
	`

	_, err := r.db.Exec(
		ctx,
		query,
		userID,
	)

	return err
}
