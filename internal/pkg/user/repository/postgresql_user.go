package repository

import (
	"context"
	"dripapp/internal/pkg/models"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgreUserRepo struct {
	conn *sqlx.DB
}

func NewPostgresStaffRepository(conn *sqlx.DB) PostgreUserRepo {
	return PostgreUserRepo{conn}
}

func (p *PostgreUserRepo) AddLogin(ctx context.Context, loginUser models.LoginUser) (models.LoginUser, error) {
	query := `INSERT into drip_profile(
                  email, 
                  password) 
                  VALUES ($1,$2) 
                  RETURNING email, password`

	var RespLoginUser models.LoginUser
	err := p.conn.GetContext(ctx, &RespLoginUser, query, loginUser.Email, loginUser.Password)
	return RespLoginUser, err
}
