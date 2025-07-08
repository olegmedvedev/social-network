package db

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int
	Name         string
	Email        string
	PasswordHash string
}

type FriendRequest struct {
	ID         int
	FromUserID int
	ToUserID   int
	CreatedAt  string
}

type SocialRepository struct {
	DB *sql.DB
}

func NewSocialRepository(db *sql.DB) *SocialRepository {
	return &SocialRepository{DB: db}
}

// User methods
func (r *SocialRepository) CreateUser(ctx context.Context, name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	var id int
	err = r.DB.QueryRowContext(ctx, `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`, name, email, string(hash)).Scan(&id)
	if err != nil {
		return nil, err
	}
	return &User{ID: id, Name: name, Email: email, PasswordHash: string(hash)}, nil
}

func (r *SocialRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u := &User{}
	err := r.DB.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *SocialRepository) GetUserByID(ctx context.Context, id int) (*User, error) {
	u := &User{}
	err := r.DB.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *SocialRepository) CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

// FriendRequest methods
func (r *SocialRepository) SendFriendRequest(ctx context.Context, fromUserID, toUserID int) (*FriendRequest, error) {
	var id int
	var createdAt string
	err := r.DB.QueryRowContext(ctx, `
        INSERT INTO friend_requests (from_user_id, to_user_id)
        VALUES ($1, $2)
        RETURNING id, created_at
    `, fromUserID, toUserID).Scan(&id, &createdAt)
	if err != nil {
		return nil, err
	}
	return &FriendRequest{ID: id, FromUserID: fromUserID, ToUserID: toUserID, CreatedAt: createdAt}, nil
}

func (r *SocialRepository) AcceptFriendRequest(ctx context.Context, requestID int) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var fromUserID, toUserID int
	err = tx.QueryRowContext(ctx, `SELECT from_user_id, to_user_id FROM friend_requests WHERE id=$1`, requestID).Scan(&fromUserID, &toUserID)
	if err != nil {
		return err
	}
	// Add friendship both ways
	_, err = tx.ExecContext(ctx, `
        INSERT INTO friends (user_id, friend_id) VALUES ($1, $2), ($2, $1)
    `, fromUserID, toUserID)
	if err != nil {
		return err
	}
	// Delete the friend request
	_, err = tx.ExecContext(ctx, `DELETE FROM friend_requests WHERE id=$1`, requestID)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (r *SocialRepository) GetFriends(ctx context.Context, userID int) ([]int, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT friend_id FROM friends WHERE user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *SocialRepository) GetIncomingFriendRequests(ctx context.Context, userID int) ([]*FriendRequest, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, from_user_id, to_user_id, created_at FROM friend_requests WHERE to_user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reqs []*FriendRequest
	for rows.Next() {
		fr := &FriendRequest{}
		if err := rows.Scan(&fr.ID, &fr.FromUserID, &fr.ToUserID, &fr.CreatedAt); err != nil {
			return nil, err
		}
		reqs = append(reqs, fr)
	}
	return reqs, nil
}

func (r *SocialRepository) GetOutgoingFriendRequests(ctx context.Context, userID int) ([]*FriendRequest, error) {
	rows, err := r.DB.QueryContext(ctx, `SELECT id, from_user_id, to_user_id, created_at FROM friend_requests WHERE from_user_id=$1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reqs []*FriendRequest
	for rows.Next() {
		fr := &FriendRequest{}
		if err := rows.Scan(&fr.ID, &fr.FromUserID, &fr.ToUserID, &fr.CreatedAt); err != nil {
			return nil, err
		}
		reqs = append(reqs, fr)
	}
	return reqs, nil
}
