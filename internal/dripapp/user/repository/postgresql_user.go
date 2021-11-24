package repository

import (
	"context"
	"database/sql"
	"dripapp/configs"
	"dripapp/internal/dripapp/models"
	"dripapp/internal/pkg/logger"
	"fmt"
	"log"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

const success = "Connection success (postgre) on: "

type PostgreUserRepo struct {
	Conn sqlx.DB
}

func NewPostgresUserRepository(config configs.PostgresConfig) (models.UserRepository, error) {
	ConnStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s sslmode=disable",
		config.User,
		config.DBName,
		config.Password,
		config.Host)

	Conn, err := sqlx.Open("postgres", ConnStr)
	if err != nil {
		return nil, err
	}

	log.Printf("%s%s", success, ConnStr)
	return &PostgreUserRepo{*Conn}, nil
}

func (p PostgreUserRepo) GetUser(ctx context.Context, email string) (models.User, error) {
	var RespUser models.User
	err := p.Conn.QueryRow(GetUserQuery, email).
		Scan(&RespUser.ID, &RespUser.Email, &RespUser.Password, &RespUser.Name, &RespUser.Gender, &RespUser.Prefer,
			&RespUser.FromAge, &RespUser.ToAge, &RespUser.Date, &RespUser.Age, &RespUser.Description, pq.Array(&RespUser.Imgs))
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.getTagsByID(ctx, RespUser.ID)
	if err != nil {
		if err != sql.ErrNoRows {
			return models.User{}, err
		}
	}

	return RespUser, nil
}

func (p PostgreUserRepo) GetUserByID(ctx context.Context, userID uint64) (models.User, error) {
	var RespUser models.User
	err := p.Conn.QueryRow(GetUserByIdAQuery, userID).
		Scan(&RespUser.ID, &RespUser.Email, &RespUser.Password, &RespUser.Name, &RespUser.Gender, &RespUser.Prefer,
			&RespUser.FromAge, &RespUser.ToAge, &RespUser.Date, &RespUser.Age, &RespUser.Description, pq.Array(&RespUser.Imgs))
	if err != nil {
		return models.User{}, err
	}

	RespUser.Tags, err = p.getTagsByID(ctx, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			return models.User{}, err
		}
	}

	return RespUser, nil
}

func (p PostgreUserRepo) CreateUser(ctx context.Context, logUserData models.LoginUser) (models.User, error) {
	var RespUser models.User
	err := p.Conn.GetContext(ctx, &RespUser, CreateUserQuery, logUserData.Email, logUserData.Password)
	return RespUser, err
}

func (p PostgreUserRepo) UpdateUser(ctx context.Context, newUserData models.User) (models.User, error) {

	if newUserData.FromAge == 0 {
		newUserData.FromAge = 18
	}
	if newUserData.ToAge == 0 {
		newUserData.ToAge = 100
	}

	var RespUser models.User
	err := p.Conn.QueryRow(UpdateUserQuery, newUserData.Email, newUserData.Name, newUserData.Gender, newUserData.Prefer,
		newUserData.FromAge, newUserData.ToAge, newUserData.Date, newUserData.Description, pq.Array(&newUserData.Imgs)).
		Scan(&RespUser.ID, &RespUser.Email, &RespUser.Password, &RespUser.Name, &RespUser.Gender, &RespUser.Prefer,
			&RespUser.FromAge, &RespUser.ToAge, &RespUser.Date, &RespUser.Age, &RespUser.Description, pq.Array(&RespUser.Imgs))
	if err != nil {
		logger.DripLogger.DebugLogging("update error")
		return models.User{}, err
	}

	err = p.deleteTags(ctx, newUserData.ID)
	if err != nil && err != sql.ErrNoRows {
		logger.DripLogger.DebugLogging("delete error")
		return models.User{}, err
	}

	if len(newUserData.Tags) != 0 {
		err = p.insertTags(ctx, newUserData.ID, newUserData.Tags)
		if err != nil {
			logger.DripLogger.DebugLogging("insert error")
			return models.User{}, err
		}
	}

	RespUser.Tags, err = p.getTagsByID(ctx, RespUser.ID)
	if err != nil && err != sql.ErrNoRows {
		logger.DripLogger.DebugLogging("get tags by id")
		return models.User{}, err
	}

	return RespUser, nil
}

func (p PostgreUserRepo) deleteTags(ctx context.Context, userId uint64) error {
	var id uint64
	err := p.Conn.QueryRow(DeleteTagsQuery, userId).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) GetTags(ctx context.Context) (map[uint64]string, error) {
	var tags []models.Tag
	err := p.Conn.Select(&tags, GetTagsQuery)
	if err != nil {
		return nil, err
	}

	tagsMap := make(map[uint64]string)

	var i uint64
	for i = 0; i < uint64(len(tags)); i++ {
		tagsMap[i] = tags[i].TagName
	}

	return tagsMap, nil
}

func (p PostgreUserRepo) getTagsByID(ctx context.Context, id uint64) ([]string, error) {
	var tags []string
	err := p.Conn.Select(&tags, GetTagsByIdQuery, id)
	if err != nil {
		return nil, err
	}

	return tags, nil
}

func (p PostgreUserRepo) getImgsByID(ctx context.Context, id uint64) ([]string, error) {
	var imgs []string
	if err := p.Conn.QueryRow(GetImgsByIDQuery, id).Scan(pq.Array(&imgs)); err != nil {
		return nil, err
	}

	return imgs, nil
}

func (p PostgreUserRepo) insertTags(ctx context.Context, id uint64, tags []string) error {
	if len(tags) == 0 {
		return nil
	}

	vals := []interface{}{}
	vals = append(vals, id)
	for _, val := range tags {
		vals = append(vals, val)
	}

	var sb strings.Builder
	sb.WriteString(InsertTagsQueryFirstPart)
	var inserts []string
	for idx := range tags {
		str := fmt.Sprintf(InsertTagsQueryParts, idx+2)
		inserts = append(inserts, str)
	}
	sb.WriteString(strings.Join(inserts, ",\n"))
	sb.WriteString(" returning id;")
	insertTagsQuery := sb.String()

	var respId uint64
	err := p.Conn.QueryRow(insertTagsQuery, vals...).Scan(&respId)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}

func (p PostgreUserRepo) UpdateImgs(ctx context.Context, id uint64, imgs []string) error {
	var user_id uint64
	err := p.Conn.QueryRow(UpdateImgsQuery, id, pq.Array(&imgs)).Scan(&user_id)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) AddReaction(ctx context.Context, currentUserId uint64, swipedUserId uint64, reactionType uint64) error {
	var id uint64
	err := p.Conn.QueryRow(AddReactionQuery, currentUserId, swipedUserId, reactionType).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) GetNextUserForSwipe(ctx context.Context, currentUser models.User) (notSwipedUsers []models.User, err error) {
	var sb strings.Builder
	sb.WriteString(GetNextUserForSwipeQuery1)
	if len(currentUser.Prefer) != 0 {
		sb.WriteString(GetNextUserForSwipeQueryPrefer)
	}
	sb.WriteString(Limit)
	GetNextUserForSwipeQuery := sb.String()

	if len(currentUser.Prefer) != 0 {
		err = p.Conn.Select(&notSwipedUsers, GetNextUserForSwipeQuery, currentUser.ID, currentUser.FromAge, currentUser.ToAge, currentUser.Prefer)
	} else {
		err = p.Conn.Select(&notSwipedUsers, GetNextUserForSwipeQuery, currentUser.ID, currentUser.FromAge, currentUser.ToAge)
	}
	if err != nil {
		return nil, err
	}

	for idx := range notSwipedUsers {
		notSwipedUsers[idx].Imgs, err = p.getImgsByID(ctx, notSwipedUsers[idx].ID)
		if err != nil {
			return nil, err
		}

		notSwipedUsers[idx].Tags, err = p.getTagsByID(ctx, notSwipedUsers[idx].ID)
		if err != nil {
			return nil, err
		}
	}

	return notSwipedUsers, nil
}

func (p PostgreUserRepo) GetUsersMatches(ctx context.Context, currentUserId uint64) ([]models.User, error) {
	var matchesUsers []models.User
	err := p.Conn.Select(&matchesUsers, GetUsersForMatchesQuery, currentUserId)
	if err != nil {
		return nil, err
	}

	for idx := range matchesUsers {
		matchesUsers[idx].Imgs, err = p.getImgsByID(ctx, matchesUsers[idx].ID)
		if err != nil {
			return nil, err
		}

		matchesUsers[idx].Tags, err = p.getTagsByID(ctx, matchesUsers[idx].ID)
		if err != nil {
			return nil, err
		}
	}

	return matchesUsers, nil
}

func (p PostgreUserRepo) GetLikes(ctx context.Context, currentUserId uint64) ([]uint64, error) {
	// type = 1 is like (dislike - 2)

	var likes []uint64
	err := p.Conn.Select(&likes, GetLikesQuery, currentUserId)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func (p PostgreUserRepo) DeleteReaction(ctx context.Context, firstUser uint64, secondUser uint64) error {
	var id uint64
	err := p.Conn.QueryRow(DeleteReactionQuery, firstUser, secondUser).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (p PostgreUserRepo) DeleteMatches(ctx context.Context, firstUser uint64, secondUser uint64) error {
	var id uint64
	err := p.Conn.QueryRow(DeleteMatchQuery, firstUser, secondUser).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

func (p PostgreUserRepo) AddMatch(ctx context.Context, firstUser uint64, secondUser uint64) error {
	var id uint64
	err := p.Conn.QueryRow(AddMatchQuery, firstUser, secondUser).Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgreUserRepo) GetUsersLikes(ctx context.Context, currentUserId uint64) ([]models.User, error) {
	var likesUsers []models.User
	err := p.Conn.Select(&likesUsers, GetUserLikesQuery, currentUserId)
	if err != nil {
		return nil, err
	}

	for idx := range likesUsers {
		likesUsers[idx].Imgs, err = p.getImgsByID(ctx, likesUsers[idx].ID)
		if err != nil {
			return nil, err
		}

		likesUsers[idx].Tags, err = p.getTagsByID(ctx, likesUsers[idx].ID)
		if err != nil {
			return nil, err
		}
	}

	return likesUsers, nil
}

func (p PostgreUserRepo) GetUsersMatchesWithSearching(ctx context.Context, currentUserId uint64, searchTmpl string) ([]models.User, error) {
	var matchesUsers []models.User
	err := p.Conn.Select(&matchesUsers, GetUsersForMatchesWithSearchingQuery, currentUserId, searchTmpl)
	if err != nil {
		return nil, err
	}

	for idx := range matchesUsers {
		matchesUsers[idx].Imgs, err = p.getImgsByID(ctx, matchesUsers[idx].ID)
		if err != nil {
			return nil, err
		}

		matchesUsers[idx].Tags, err = p.getTagsByID(ctx, matchesUsers[idx].ID)
		if err != nil {
			return nil, err
		}
	}

	return matchesUsers, nil
}

func (p PostgreUserRepo) GetReports(ctx context.Context) (map[uint64]string, error) {
	var reports []models.Report
	err := p.Conn.Select(&reports, GetReportsQuery)
	fmt.Println(352, err)
	if err != nil {
		return nil, err
	}

	reportsMap := make(map[uint64]string)

	var i uint64
	for i = 0; i < uint64(len(reports)); i++ {
		reportsMap[i] = reports[i].ReportDesc
	}

	return reportsMap, nil
}

func (p PostgreUserRepo) AddReport(ctx context.Context, report models.NewReport) error {
	var reportId uint64
	if err := p.Conn.QueryRow(GetReportIdFromDescQuery, report.ReportDesc).Scan(&reportId); err != nil {
		return err
	}

	var respId uint64
	err := p.Conn.QueryRow(AddReportToProfileQuery, report.ToId, reportId).Scan(&respId)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}

func (p PostgreUserRepo) GetReportsCount(ctx context.Context, userId uint64) (uint64, error) {
	var curCount uint64
	if err := p.Conn.QueryRow(GetReportsCountQuery, userId).Scan(&curCount); err != nil {
		return curCount, err
	}

	return curCount, nil
}

func (p PostgreUserRepo) GetReportsWithMaxCount(ctx context.Context, userId uint64) (uint64, error) {
	var reportId uint64
	if err := p.Conn.QueryRow(GetReportsIdWithMaxCountQuery, userId).Scan(&reportId); err != nil {
		return reportId, err
	}

	return reportId, nil
}

func (p PostgreUserRepo) GetReportDesc(ctx context.Context, reportId uint64) (string, error) {
	var reportDesc string
	if err := p.Conn.QueryRow(GetReportDescFromIdQuery, reportId).Scan(&reportDesc); err != nil {
		return reportDesc, err
	}

	return reportDesc, nil
}

func (p PostgreUserRepo) UpdateReportStatus(ctx context.Context, userId uint64, reportStatus string) error {
	var respId uint64
	err := p.Conn.QueryRow(UpdateProfilesReportStatusQuery, userId, reportStatus).Scan(&respId)

	if err != nil {
		if err != sql.ErrNoRows {
			return err
		}
	}

	return nil
}
