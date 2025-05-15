package repository

//
//import (
//	"context"
//	"github.com/rwpp/RzWeLook/bff/service"
//
//	//userv1 "github.com/rwpp/RzWeLook/api/proto/gen/user/v1"
//	"github.com/rwpp/RzWeLook/bff/domain"
//	"github.com/rwpp/RzWeLook/bff/repository/dao"
//)
//
//// AuthorRepository 封装user的client用于获取用户信息
//type AuthorRepository interface {
//	// FindAuthor id为文章id
//	FindAuthor(ctx context.Context, id int64) (domain.Author, error)
//}
//
//type GrpcAuthorRepository struct {
//	userSvc service.UserService
//	dao     dao.ArticleDAO
//}
//
//func NewGrpcAuthorRepository(articleDao dao.ArticleDAO, userSvc service.UserService) AuthorRepository {
//	return &GrpcAuthorRepository{
//		userSvc: userSvc,
//		dao:     articleDao,
//	}
//}
//
//func (g *GrpcAuthorRepository) FindAuthor(ctx context.Context, id int64) (domain.Author, error) {
//	art, err := g.dao.GetPubById(ctx, id)
//	if err != nil {
//		return domain.Author{}, nil
//	}
//	u, err := g.userSvc.Profile(ctx, &userv1.ProfileRequest{
//		Id: art.AuthorId,
//	})
//	if err != nil {
//		return domain.Author{}, err
//	}
//	return domain.Author{
//		Id:   u.User.Id,
//		Name: u.User.Nickname,
//	}, nil
//}
