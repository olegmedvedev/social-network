package db

import (
    "database/sql"
    _ "github.com/lib/pq"
    "log"
    "fmt"
    "social-network/internal/config"
    "golang.org/x/crypto/bcrypt"
    "context"
)

var DB *sql.DB

func Init(cfg *config.Config) {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
    )
    var err error
    DB, err = sql.Open("postgres", dsn)
    if err != nil {
        log.Fatal("Failed to connect to DB:", err)
    }
    if err = DB.Ping(); err != nil {
        log.Fatal("Failed to ping DB:", err)
    }
}

type User struct {
    ID    int
    Name  string
    Email string
    PasswordHash string
}

func CreateUser(ctx context.Context, name, email, password string) (*User, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return nil, err
    }
    var id int
    err = DB.QueryRowContext(ctx, `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3) RETURNING id`, name, email, string(hash)).Scan(&id)
    if err != nil {
        return nil, err
    }
    return &User{ID: id, Name: name, Email: email, PasswordHash: string(hash)}, nil
}

func GetUserByEmail(ctx context.Context, email string) (*User, error) {
    u := &User{}
    err := DB.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM users WHERE email=$1`, email).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func GetUserByID(ctx context.Context, id int) (*User, error) {
    u := &User{}
    err := DB.QueryRowContext(ctx, `SELECT id, name, email, password_hash FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash)
    if err != nil {
        return nil, err
    }
    return u, nil
}

func CheckPassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
} 