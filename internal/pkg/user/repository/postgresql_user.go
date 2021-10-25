package repository

import (
	"context"
	"dripapp/internal/pkg/models"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgreUserRepo struct {
	conn *sqlx.DB
}

func NewPostgresUserRepository(conn *sqlx.DB) PostgreUserRepo {
	return PostgreUserRepo{conn}
}

func (p *PostgreUserRepo) CreateUser(ctx context.Context, loginUser models.LoginUser) (models.User, error) {
	query := `INSERT into profile(
                  email, 
                  password) 
                  VALUES ($1,$2) 
                  RETURNING email, password;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, loginUser.Email, loginUser.Password)
	return RespUser, err
}

func (p *PostgreUserRepo) CreateUserAndProfile(ctx context.Context, user models.User) (models.User, error) {
	query := `insert into profile(name, email, password, date, description, imgs)
		values($1,$2,$3,$4,$5,$6)
		RETURNING id, name, email, password, email, password, date, description;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, user.Name, user.Email, user.Password, user.Date,
		user.Description, pq.Array(user.Imgs))
	if err != nil {
		return models.User{}, err
	}

	err = p.insertTags(ctx, RespUser.ID, user.Tags)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Imgs, err = p.getImgs(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Age, err = models.GetAgeFromDate(RespUser.Date)
	if err != nil {
		log.Fatal(err)
	}

	RespUser.Tags, err = p.getTagsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, err
}

func (p *PostgreUserRepo) UpdateUser(ctx context.Context, newUserData models.User) (models.User, error) {
	query := `update profile
		set name=$1, date=$3, description=$4, imgs=&5
		where email=$2
		RETURNING id, email, password, name, email, password, date, description;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, newUserData.Name, newUserData.Email, newUserData.Date,
		newUserData.Description, pq.Array(newUserData.Imgs))

	RespUser.Imgs, err = p.getImgs(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}
	RespUser.Tags, err = p.getTagsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, err
}

func (p *PostgreUserRepo) DeleteUser(ctx context.Context, user models.User) error {
	query := `delete from profile where id=$1`

	if err := p.conn.QueryRow(query, user.ID).Scan(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) GetTags(ctx context.Context) (map[uint64]string, error) {
	query := `select id, tag_name from tag;`

	var tags []models.Tag
	err := p.conn.Select(&tags, query)
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[uint64]string)

	for _, val := range tags {
		tagsMap[val.Id] = val.Tag_Name
	}

	return tagsMap, nil
}

func (p *PostgreUserRepo) DropSwipes(ctx context.Context) error {
	query := `delete from reactions`

	if err := p.conn.QueryRow(query).Scan(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) DropUsers(ctx context.Context) error {
	query := `
	delete from profile_tag;
	delete from matches;
	delete from reactions;
	delete from profile;`

	if err := p.conn.QueryRow(query).Scan(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) GetUserByID(ctx context.Context, userID uint64) (models.User, error) {
	query := `select id, name, email, password, date, description
		from profile
		where id = $1;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, userID)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.getTagsByID(ctx, 1)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Imgs, err = p.getImgs(ctx, 1)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, nil
}

func (p *PostgreUserRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	query := `select id, name, email, password, date, description
		from profile
		where email = $1;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, email)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.getTagsByID(ctx, 1)
	if err != nil {
		return models.User{}, err
	}
	RespUser.Imgs, err = p.getImgs(ctx, 1)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, nil
}

func (p *PostgreUserRepo) getTagsByID(ctx context.Context, id uint64) ([]string, error) {
	sel := `select
				tag_name
			from
				profile p
				join profile_tag pt on(pt.profile_id = p.id)
				join tag t on(pt.tag_id = t.id)
			where
				p.id = $1;`

	var tags []string
	err := p.conn.Select(&tags, sel, id)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (p *PostgreUserRepo) CreateTag(ctx context.Context, tag_name string) error {
	sel := "insert into tag(tag_name) values($1);"

	if err := p.conn.QueryRow(sel, tag_name).Scan(); err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) getImgs(ctx context.Context, id uint64) ([]string, error) {
	sel := "SELECT imgs FROM profile WHERE id=$1;"

	var imgs []string
	if err := p.conn.QueryRow(sel, id).Scan(pq.Array(&imgs)); err != nil {
		return nil, err
	}

	return imgs, nil
}

func (p *PostgreUserRepo) insertTags(ctx context.Context, id uint64, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	query := "insert into profile_tag(profile_id, tag_id)\nvalues\n"

	vals := []interface{}{}
	vals = append(vals, id)
	for _, val := range tags {
		vals = append(vals, val)
	}

	var sb strings.Builder
	sb.WriteString(query)
	var inserts []string
	for idx := range tags {
		str := fmt.Sprintf("($1, (select id from tag where tag_name=$%d))", idx+2)
		inserts = append(inserts, str)
	}
	sb.WriteString(strings.Join(inserts, ",\n"))
	sb.WriteString(";")
	query = sb.String()

	fmt.Println(query)
	fmt.Println(vals)

	stmt, _ := p.conn.Prepare(query)
	_, err := stmt.Exec(vals...)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) UpdateImgs(ctx context.Context, id uint64, imgs []string) error {
	query := `update profile set imgs=$2 where id=$1;`

	err := p.conn.QueryRow(query, id, pq.Array(imgs))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (p *PostgreUserRepo) AddSwipedUsers(ctx context.Context, currentUserId uint64, swipedUserId uint64, type_name string) error {
	query := "insert into reactions(id1, id2, type) values ($1,$2,$3);"

	stmt, _ := p.conn.Prepare(query)
	_, err := stmt.Exec(currentUserId, swipedUserId, type_name)
	if err != nil {
		return err
	}

	return nil
}

func (p *PostgreUserRepo) IsSwiped(ctx context.Context, userID, swipedUserID uint64) (bool, error) {
	query := `select exists(select id1, id2 from reactions where id1=$1 and id2=$2)`

	var resp bool
	err := p.conn.GetContext(ctx, &resp, query, userID, swipedUserID)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (p *PostgreUserRepo) GetNextUserForSwipe(ctx context.Context, currentUserId uint64) (models.User, error) {
	return models.User{}, nil
}
