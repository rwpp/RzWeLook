package web

import (
	"fmt"
	"github.com/ecodeclub/ekit/slice"
	"github.com/gin-gonic/gin"
	articlev1 "github.com/rwpp/RzWeLook/api/proto/gen/article/v1"
	intrv1 "github.com/rwpp/RzWeLook/api/proto/gen/intr/v1"
	rewardv1 "github.com/rwpp/RzWeLook/api/proto/gen/reward/v1"

	"github.com/rwpp/RzWeLook/bff/web/jwt"
	"github.com/rwpp/RzWeLook/pkg/ginx"
	"github.com/rwpp/RzWeLook/pkg/logger"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

var _ handler = (*ArticleHandler)(nil)

type ArticleHandler struct {
	svc     articlev1.ArticleServiceClient
	intrSvc intrv1.InteractiveServiceClient
	reward  rewardv1.RewardServiceClient
	l       logger.LoggerV1
	biz     string
}

func NewArticleHandler(svc articlev1.ArticleServiceClient, intrSvc intrv1.InteractiveServiceClient, l logger.LoggerV1, reward rewardv1.RewardServiceClient) *ArticleHandler {
	return &ArticleHandler{
		svc:     svc,
		l:       l,
		reward:  reward,
		biz:     "article",
		intrSvc: intrSvc,
	}
}

func (a *ArticleHandler) RegisterRoutes(s *gin.Engine) {
	g := s.Group("/articles")
	// 在有 list 等路由的时候，无法这样注册
	// g.GET("/:id", a.Detail)
	g.GET("/detail/:id", a.Detail)
	// 理论上来说应该用 GET的，但是我实在不耐烦处理类型转化
	// 直接 POST，JSON 转一了百了。
	g.POST("/list", ginx.WrapBodyAndToken[ListReq, jwt.UserClaims](a.List))

	g.POST("/edit", a.Edit)
	g.POST("/publish", a.Publish)
	g.POST("/withdraw", a.Withdraw)

	pub := g.Group("/pub")
	//pub.GET("/pub", a.PubList)
	pub.GET("/detail/:id",
		ginx.WrapToken[jwt.UserClaims](a.PubDetail))
	pub.POST("/like", ginx.WrapBodyAndToken[LikeReq](a.Like))
	pub.POST("/collect", ginx.WrapBodyAndToken[CollectReq](a.Collect))
	// 打赏
	pub.POST("/reward", ginx.WrapBodyAndToken[RewardReq](a.Reward))
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	type Req struct {
		Id int64 `json:"id"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("users").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	_, err := a.svc.Withdraw(ctx, &articlev1.WithdrawRequest{
		Article: &articlev1.Article{
			Id: req.Id,
			Author: &articlev1.Author{
				Id: usr.Uid,
			},
		},
	})
	if err != nil {
		a.l.Error("设置为尽自己可见失败", logger.Error(err),
			logger.Field{Key: "id", Value: req.Id})
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "OK",
	})
}

func (a *ArticleHandler) List(ctx *gin.Context, req ListReq, usr jwt.UserClaims) (ginx.Result, error) {

	// 对于批量接口来说，要小心批次大小
	if req.Limit > 100 {
		a.l.Error("获得用户会话信息失败，LIMIT过大")
		return ginx.Result{
			Code: 4,
			// 我会倾向于不告诉攻击者批次太大
			// 因为一般你和前端一起完成任务的时候
			// 你们是协商好了的，所以会进来这个分支
			// 就表明是有人跟你过不去
			Msg: "请求有误",
		}, nil
	}
	arts, err := a.svc.List(ctx, &articlev1.ListRequest{Author: usr.Uid,
		Offset: req.Offset, Limit: req.Limit})
	if err != nil {
		a.l.Error("获得用户会话信息失败")
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, nil
	}
	return ginx.Result{
		Data: slice.Map[*articlev1.Article, ArticleVo](arts.Articles,
			func(idx int, src *articlev1.Article) ArticleVo {
				return ArticleVo{
					Id:       src.Id,
					Title:    src.Title,
					Abstract: src.Abstract,
					Status:   src.Status,
					// 这个列表请求，不需要返回内容
					//Content: src.Content,
					// 这个是创作者看自己的文章列表，也不需要这个字段
					//Author: src.Author
					Ctime: src.Ctime.AsTime().Format(time.DateTime),
					Utime: src.Utime.AsTime().Format(time.DateTime),
				}
			}),
	}, nil
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "参数错误",
		})
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("users").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	resp, err := a.svc.GetById(ctx, &articlev1.GetByIdRequest{Id: id})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得文章信息失败", logger.Error(err))
		return
	}
	art := resp.GetArticle()
	// 这是不借助数据库查询来判定的方法
	if art.Author.Id != usr.Uid {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			// 也不需要告诉前端究竟发生了什么
			Msg: "输入有误",
		})
		// 如果公司有风控系统，这个时候就要上报这种非法访问的用户了。
		a.l.Error("非法访问文章，创作者 ID 不匹配",
			logger.Int64("uid", usr.Uid))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: ArticleVo{
			Id:    art.Id,
			Title: art.Title,
			// 不需要这个摘要信息
			//Abstract: art.Abstract(),
			Status:  art.Status,
			Content: art.Content,
			// 这个是创作者看自己的文章列表，也不需要这个字段
			//Author: art.Author
			Ctime: art.Ctime.AsTime().Format(time.DateTime),
			Utime: art.Utime.AsTime().Format(time.DateTime),
		},
	})
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("users").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Publish(ctx, &articlev1.PublishRequest{
		Article: &articlev1.Article{
			Id:      req.Id,
			Title:   req.Title,
			Content: req.Content,
			Author: &articlev1.Author{
				Id: usr.Uid,
			},
		},
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("发表文章失败", logger.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (a *ArticleHandler) Edit(ctx *gin.Context) {
	var req ArticleReq
	if err := ctx.Bind(&req); err != nil {
		a.l.Error("反序列化请求失败", logger.Error(err))
		return
	}
	usr, ok := ctx.MustGet("users").(jwt.UserClaims)
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("获得用户会话信息失败")
		return
	}
	id, err := a.svc.Save(ctx, &articlev1.SaveRequest{Article: req.toDTO(usr.Uid)})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		a.l.Error("保存数据失败", logger.Field{Key: "error", Value: err})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg:  "OK",
		Data: id,
	})
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context, usr jwt.UserClaims) (ginx.Result, error) {
	idstr := ctx.Param("id")
	id, err := strconv.ParseInt(idstr, 10, 64)
	if err != nil {
		a.l.Error("前端输入的 ID 不对", logger.Error(err))
		return ginx.Result{
			Code: 4,
			Msg:  "参数错误",
		}, fmt.Errorf("查询文章详情的 ID %s 不正确, %w", idstr, err)
	}
	// 使用 error group 来同时查询数据
	var (
		eg           errgroup.Group
		artResp      *articlev1.GetPublishedByIdResponse
		intrResponse *intrv1.GetResponse
	)
	eg.Go(func() error {
		artResp, err = a.svc.GetPublishedById(ctx, &articlev1.GetPublishedByIdRequest{
			Id: id, Uid: usr.Uid,
		})
		return err
		//if art.Author.Id != usr.Uid {
		//	return ginx.Result{
		//		Code: 4,
		//		Msg:  "输入有误",
		//	}, fmt.Errorf("输入的 ID 不正确")
		//}
	})

	eg.Go(func() error {
		var er error
		intrResponse, er = a.intrSvc.Get(ctx, &intrv1.GetRequest{
			Biz:   a.biz,
			BizId: id,
			Uid:   usr.Uid,
		})
		return er
	})
	err = eg.Wait()
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, fmt.Errorf("获取文章信息失败 %w", err)
	}
	// 直接异步操作，在确定我们获取到了数据之后再来操作
	//go func() {
	//	_, er := a.intrSvc.IncrReadCnt(ctx, &intrv1.IncrReadCntRequest{
	//		Biz:   a.biz,
	//		BizId: art.Id,
	//	})
	//	if er != nil {
	//		a.l.Error("增加文章阅读数失败",
	//			logger.Error(err))
	//	}
	//}()
	art := artResp.GetArticle()
	intr := intrResponse.Intr
	return ginx.Result{
		Data: ArticleVo{
			Id:      art.Id,
			Title:   art.Title,
			Status:  art.Status,
			Content: art.Content,
			// 要把作者信息带出去
			Author:     art.Author.Name,
			Ctime:      art.Ctime.AsTime().Format(time.DateTime),
			Utime:      art.Utime.AsTime().Format(time.DateTime),
			ReadCnt:    intr.ReadCnt,
			CollectCnt: intr.CollectCnt,
			LikeCnt:    intr.LikeCnt,
			Liked:      intr.Liked,
			Collected:  intr.Collected,
		},
	}, nil
}

func (a *ArticleHandler) Like(ctx *gin.Context, req LikeReq, uc jwt.UserClaims) (ginx.Result, error) {
	var err error
	if req.Like {
		_, err = a.intrSvc.Like(ctx, &intrv1.LikeRequest{
			Biz:   a.biz,
			BizId: req.Id,
			Uid:   uc.Uid,
		})
	} else {
		_, err = a.intrSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz:   a.biz,
			BizId: req.Id,
			Uid:   uc.Uid,
		})
	}

	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}

func (a *ArticleHandler) Reward(ctx *gin.Context, req RewardReq, uc jwt.UserClaims) (ginx.Result, error) {
	artResp, err := a.svc.GetPublishedById(ctx.Request.Context(), &articlev1.GetPublishedByIdRequest{
		Id: req.Id,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	art := artResp.GetArticle()
	resp, err := a.reward.PreReward(ctx.Request.Context(), &rewardv1.PreRewardRequest{
		Biz:       "article",
		BizId:     art.Id,
		BizName:   art.Title,
		TargetUid: art.Author.GetId(),
		Uid:       uc.Uid,
		Amt:       req.Amt,
	})
	if err != nil {
		return ginx.Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return ginx.Result{
		Data: map[string]any{
			"codeURL": resp.CodeUrl,
			"rid":     resp.Rid,
		},
	}, nil
}

func (a *ArticleHandler) Collect(ctx *gin.Context, req CollectReq,
	uc jwt.UserClaims) (Result, error) {
	_, err := a.intrSvc.Collect(ctx, &intrv1.CollectRequest{
		Biz:   a.biz,
		BizId: req.Id,
		Cid:   req.Cid,
		Uid:   uc.Uid,
	})
	if err != nil {
		return Result{
			Code: 5,
			Msg:  "系统错误",
		}, err
	}
	return Result{Msg: "OK"}, nil
}
