package repository

import (
	"context"
	"dripapp/configs"
	"dripapp/internal/pkg/models"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const success = "Connection success (postgre) on: "

type PostgreUserRepo struct {
	conn sqlx.DB
}

func NewPostgresUserRepository(config configs.PostgresConfig) (*PostgreUserRepo, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		configs.Postgres.User,
		configs.Postgres.Password,
		configs.Postgres.DBName)

	conn, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// query, err := ioutil.ReadFile("docker/postgres_scripts/dump.sql")
	// if err != nil {
	// 	return nil, err
	// }
	// strQuery := string(query)
	// if _, err := conn.Exec(strQuery); err != nil {
	// 	return nil, err
	// }

	log.Printf("%s%s", success, connStr)
	return &PostgreUserRepo{*conn}, nil
}

func (p PostgreUserRepo) Init() error {
	query, err := ioutil.ReadFile("docker/postgres_scripts/dump.sql")
	if err != nil {
		return err
	}
	strQuery := string(query)

	if _, err := p.conn.Exec(strQuery); err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) DeletePhoto(ctx context.Context, user models.User, photo string) error {
	return nil
}

func (p PostgreUserRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	query := `select id, name, email, password, date, description
		from profile
		where email = $1;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, email)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.GetTagsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}
	RespUser.Imgs, err = p.GetImgsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, nil
}

func (p PostgreUserRepo) GetUserByID(ctx context.Context, userID uint64) (models.User, error) {
	query := `select id, name, email, password, date, description
		from profile
		where id = $1;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, userID)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.GetTagsByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Imgs, err = p.GetImgsByID(ctx, userID)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, nil
}

func (p PostgreUserRepo) CreateUser(ctx context.Context, logUserData models.LoginUser) (models.User, error) {
	query := `INSERT into profile(
                  email,
                  password)
                  VALUES ($1,$2)
                  RETURNING id, email, password;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, logUserData.Email, logUserData.Password)
	return RespUser, err
}

func (p PostgreUserRepo) CreateUserAndProfile(ctx context.Context, user models.User) (models.User, error) {
	query := `insert into profile(name, email, password, date, description, imgs)
		values($1,$2,$3,$4,$5,$6)
		RETURNING id, name, email, password, email, password, date, description;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, user.Name, user.Email, user.Password, user.Date,
		user.Description, pq.Array(&user.Imgs))
	if err != nil {
		return models.User{}, err
	}

	err = p.InsertTags(ctx, RespUser.ID, user.Tags)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Imgs, err = p.GetImgsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	RespUser.Age, err = models.GetAgeFromDate(RespUser.Date)
	if err != nil {
		log.Fatal(err)
	}

	RespUser.Tags, err = p.GetTagsByID(ctx, RespUser.ID)
	if err != nil {
		return models.User{}, err
	}

	return RespUser, err
}

func (p PostgreUserRepo) UpdateUser(ctx context.Context, newUserData models.User) (models.User, error) {
	query := `update profile
		set name=$1, date=$3, description=$4, imgs=$5
		where email=$2
		RETURNING id, email, password, name, email, password, date, description;`

	var RespUser models.User
	err := p.conn.GetContext(ctx, &RespUser, query, newUserData.Name, newUserData.Email, newUserData.Date,
		newUserData.Description, pq.Array(&newUserData.Imgs))

	if len(newUserData.Tags) != 0 {
		err = p.DeleteTags(ctx, newUserData.ID)
		if err != nil {
			return models.User{}, err
		}
		err = p.InsertTags(ctx, newUserData.ID, newUserData.Tags)
		if err != nil {
			return models.User{}, err
		}
	}

	if len(newUserData.Imgs) != 0 {
		RespUser.Imgs, err = p.GetImgsByID(ctx, RespUser.ID)
		if err != nil {
			return models.User{}, err
		}
	}

	if len(newUserData.Tags) != 0 {
		RespUser.Tags, err = p.GetTagsByID(ctx, RespUser.ID)
		if err != nil {
			return models.User{}, err
		}
	}

	return RespUser, err
}

func (p PostgreUserRepo) DeleteTags(ctx context.Context, userId uint64) error {
	query := `delete from profile_tag where profile_id=$1`

	stmt, _ := p.conn.Prepare(query)
	_, err := stmt.Exec(userId)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) DeleteUser(ctx context.Context, user models.User) error {
	query := `delete from profile where id=$1`

	if err := p.conn.QueryRow(query, user.ID).Scan(); err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) GetTags(ctx context.Context) (map[uint64]string, error) {
	query := `select tag_name from tag;`

	var tags []models.Tag
	err := p.conn.Select(&tags, query)
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[uint64]string)

	var i uint64
	for i = 0; i < uint64(len(tags)); i++ {
		tagsMap[i] = tags[i].Tag_Name
	}

	return tagsMap, nil
}

func (p PostgreUserRepo) DropSwipes(ctx context.Context) error {
	query := `delete from reactions`

	if err := p.conn.QueryRow(query).Scan(); err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) DropUsers(ctx context.Context) error {
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

func (p PostgreUserRepo) GetTagsByID(ctx context.Context, id uint64) ([]string, error) {
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

func (p PostgreUserRepo) GetImgsByID(ctx context.Context, id uint64) ([]string, error) {
	sel := "SELECT imgs FROM profile WHERE id=$1;"

	var imgs []string
	if err := p.conn.QueryRow(sel, id).Scan(pq.Array(&imgs)); err != nil {
		return nil, err
	}

	return imgs, nil
}

func (p PostgreUserRepo) CreateTag(ctx context.Context, tag_name string) error {
	sel := "insert into tag(tag_name) values($1);"

	if err := p.conn.QueryRow(sel, tag_name).Scan(); err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) InsertTags(ctx context.Context, id uint64, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	query := "insert into profile_tag(profile_id, tag_id) values"

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

	stmt, _ := p.conn.Prepare(query)
	_, err := stmt.Exec(vals...)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) AddPhoto(ctx context.Context, user models.User, newPhoto io.Reader) error {
	photoPath := getPathUserPhoto(user) + "/" + user.GetNameToNewPhoto()

	savedPhoto, err := os.Create(photoPath)
	if err != nil {
		return err
	}
	defer savedPhoto.Close()

	_, err = io.Copy(savedPhoto, newPhoto)
	if err != nil {
		return err
	}

	user.SaveNewPhoto()

	return nil
}

func (p PostgreUserRepo) UpdateImgs(ctx context.Context, id uint64, imgs []string) error {
	query := `update profile set imgs=$2 where id=$1;`

	err := p.conn.QueryRow(query, id, pq.Array(&imgs))
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (p PostgreUserRepo) AddSwipedUsers(ctx context.Context, currentUserId uint64, swipedUserId uint64, type_name string) error {
	query := "insert into reactions(id1, id2, type) values ($1,$2,$3);"

	stmt, _ := p.conn.Prepare(query)
	_, err := stmt.Exec(currentUserId, swipedUserId, type_name)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) IsSwiped(ctx context.Context, userID, swipedUserID uint64) (bool, error) {
	query := `select exists(select id1, id2 from reactions where id1=$1 and id2=$2)`

	var resp bool
	err := p.conn.GetContext(ctx, &resp, query, userID, swipedUserID)
	if err != nil {
		return false, err
	}
	return resp, nil
}

func (p PostgreUserRepo) GetNextUserForSwipe(ctx context.Context, currentUserId uint64) ([]models.User, error) {
	query := `select
				op.id,
				op.name,
				op.email,
				op.password,
				op.date,
				op.description
			from
				(
				select
					r.id1 as rid1,
					r.id2 as rid2
				from
					reactions r
				where
					r.id1 = $1
				) prevReact
				right join profile op on prevReact.rid2 = op.id
			where
				prevReact.rid1 is null
				and op.id <> $1
				and op.name <> '' and op.date <> '' limit 5;`

	var notSwipedUser []models.User
	err := p.conn.Select(&notSwipedUser, query, currentUserId)
	if err != nil {
		return nil, err
	}

	for idx := range notSwipedUser {
		notSwipedUser[idx].Age, err = models.GetAgeFromDate(notSwipedUser[idx].Date)
		if err != nil {
			return nil, err
		}

		notSwipedUser[idx].Imgs, err = p.GetImgsByID(ctx, currentUserId)
		if err != nil {
			return nil, err
		}

		notSwipedUser[idx].Tags, err = p.GetTagsByID(ctx, currentUserId)
		if err != nil {
			return nil, err
		}
	}

	return notSwipedUser, nil
}

func (p PostgreUserRepo) GetUsersMatches(ctx context.Context, currentUserId uint64) ([]models.User, error) {
	query := `select
				op.id,
				op.name,
				op.email,
				op.password,
				op.date,
				op.description
			from profile p
			join matches m on (p.id = m.id1)
			join matches om on (om.id1 = m.id2 and om.id2 = m.id1)
			join profile op on (op.id = om.id1)
			where p.id = $1`

	var matchesUsers []models.User
	err := p.conn.Select(&matchesUsers, query, currentUserId)
	if err != nil {
		return nil, err
	}

	for idx := range matchesUsers {
		matchesUsers[idx].Age, err = models.GetAgeFromDate(matchesUsers[idx].Date)
		if err != nil {
			return nil, err
		}

		matchesUsers[idx].Imgs, err = p.GetImgsByID(ctx, currentUserId)
		if err != nil {
			return nil, err
		}

		matchesUsers[idx].Tags, err = p.GetTagsByID(ctx, currentUserId)
		if err != nil {
			return nil, err
		}
	}

	return matchesUsers, nil
}
