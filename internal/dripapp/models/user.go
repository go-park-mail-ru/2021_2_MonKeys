package models

import (
	"context"
  "dripapp/internal/pkg/hasher"
	"io"
)

type User struct {
	ID          uint64   `json:"id,omitempty"`
	Email       string   `json:"email,omitempty"`
	Password    string   `json:"-"`
	Name        string   `json:"name,omitempty"`
	Gender      string   `json:"gender,omitempty"`
	Prefer      string   `json:"prefer,omitempty"`
	FromAge     uint8    `json:"fromage,omitempty"`
	ToAge       uint8    `json:"toage,omitempty"`
	Date        string   `json:"date,omitempty"`
	Age         string   `json:"age,omitempty"`
	Description string   `json:"description,omitempty"`
	Imgs        []string `json:"imgs,omitempty"`
	Tags        []string `json:"tags,omitempty"`
  ReportStatus string   `json:"reportStatus,omitempty"`
)

const (
	LikeResction    = 1
	DislikeReaction = 2
)

const (
	ReportLimit      = 3
	FakeReport       = "Фалишивый профиль/спам"
	AggressionReport = "Непристойное общение"
	SkamReport       = "Скам"
	UnderageReport   = "Несовершеннолетний пользователь"
)

}

type LoginUser struct {
	ID       uint64 `json:"id,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserReaction struct {
	Id       uint64 `json:"id"`
	Reaction uint64 `json:"reaction"`
}

type Match struct {
	Match bool `json:"match"`
}

type Tag struct {
	TagName string `json:"tagText"`
}

type Tags struct {
	AllTags map[uint64]Tag `json:"allTags"`
	Count   uint64         `json:"tagsCount"`
}

type Matches struct {
	AllUsers map[uint64]User `json:"allUsers"`
	Count    string          `json:"matchesCount"`
}

type Likes struct {
	AllUsers map[uint64]User `json:"allUsers"`
	Count    string          `json:"likesCount"`
}

type Search struct {
	SearchingTmpl string `json:"searchTmpl"`
}

type Report struct {
	ReportDesc string `json:"reportDesc"`
}

type Reports struct {
	AllReports map[uint64]Report `json:"allReports"`
	Count      uint64            `json:"reportsCount"`
}

type NewReport struct {
	ToId       uint64 `json:"toId"`
	ReportDesc string `json:"reportDesc"`
}

type UserReportsCount struct {
	Count uint64 `json:"userReportsCount"`
}


// ArticleUsecase represent the article's usecases
type UserUsecase interface {
	CurrentUser(c context.Context) (User, error)
	EditProfile(c context.Context, newUserData User) (User, error)
	AddPhoto(c context.Context, photo io.Reader, fileName string) (Photo, error)
	DeletePhoto(c context.Context, photo Photo) error
	Login(c context.Context, logUserData LoginUser) (User, error)
	Signup(c context.Context, logUserData LoginUser) (User, error)
	NextUser(c context.Context) ([]User, error)
	GetAllTags(c context.Context) (Tags, error)
	UsersMatches(c context.Context) (Matches, error)
	Reaction(c context.Context, reactionData UserReaction) (Match, error)
	UserLikes(c context.Context) (Likes, error)
	UsersMatchesWithSearching(c context.Context, searchData Search) (Matches, error)
	GetAllReports(c context.Context) (Reports, error)
	AddReport(c context.Context, report NewReport) error
}

// ArticleRepository represent the article's repository contract
type UserRepository interface {
	GetUser(ctx context.Context, email string) (User, error)
	GetUserByID(ctx context.Context, userID uint64) (User, error)
	CreateUser(ctx context.Context, logUserData LoginUser) (User, error)
	UpdateUser(ctx context.Context, newUserData User) (User, error)
	GetTags(ctx context.Context) (map[uint64]string, error)
	UpdateImgs(ctx context.Context, id uint64, imgs []string) error
	AddReaction(ctx context.Context, currentUserId uint64, swipedUserId uint64, reactionType uint64) error
	GetNextUserForSwipe(ctx context.Context, currentUser User) ([]User, error)
	GetUsersMatches(ctx context.Context, currentUserId uint64) ([]User, error)
	GetLikes(ctx context.Context, currentUserId uint64) ([]uint64, error)
	DeleteReaction(ctx context.Context, firstUser uint64, secondUser uint64) error
	DeleteMatches(ctx context.Context, firstUser uint64, secondUser uint64) error
	AddMatch(ctx context.Context, firstUser uint64, secondUser uint64) error
	GetUsersLikes(ctx context.Context, currentUserId uint64) ([]User, error)
	GetUsersMatchesWithSearching(ctx context.Context, currentUserId uint64, searchTmpl string) ([]User, error)
	GetReports(ctx context.Context) (map[uint64]string, error)
	AddReport(ctx context.Context, report NewReport) error
	GetReportsCount(ctx context.Context, userId uint64) (uint64, error)
	GetReportsWithMaxCountCount(ctx context.Context, userId uint64) (uint64, error)
	GetReportDesc(ctx context.Context, reportId uint64) (string, error)
	UpdateReportStatus(ctx context.Context, userId uint64, reportStatus string) error
}
