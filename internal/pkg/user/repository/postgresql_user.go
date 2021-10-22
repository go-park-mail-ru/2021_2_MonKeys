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

func (p *PostgreUserRepo) CreateUser(ctx context.Context, loginUser models.LoginUser) (models.User, error) {
	query := `INSERT into drip_profile(
                  email, 
                  password) 
                  VALUES ($1,$2) 
                  RETURNING email, password`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, loginUser.Email, loginUser.Password)
	return RespUser, err
}

func (p *PostgreUserRepo) GetUserByID(ctx context.Context, userID uint64) (*models.User, error) {
	query := `select id, name, email, password, date, description, img 
		from drip_profile
		where id = $1`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, userID)
	return &RespUser, err
}

func (p *PostgreUserRepo) UpdateUser(ctx context.Context, newUserData *models.User) error {
	query := `update drip_profile
		set name=$1, date=$3, description=$4, img=$5
		where email=$2,;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, newUserData.Name, newUserData.Email, newUserData.Date,
		newUserData.Description, newUserData.ImgSrc)
	return err
}
