package routes

import (
	"net/http"

	"github.com/emicklei/go-restful/v3"
	pkg "github.com/vitorbgouveia/start-project-go/package"
	"github.com/vitorbgouveia/start-project-go/package/dto"
	"github.com/vitorbgouveia/start-project-go/package/services"
	"github.com/vitorbgouveia/start-project-go/package/web"
	"go.uber.org/zap"
)

var (
	ErrGetBalance     = "fail to get balance"
	ErrCreateWallet   = "fail to create wallet"
	ErrWithdrawWallet = "fail to withdraw wallet"
	ErrDepositWallet  = "fail to deposit wallet"

	ErrReadReqBody        = "could not read body of request"
	ErrDeserializeReqBody = "could not parse body of request"
)

type walletRoutes struct {
	Logger *zap.SugaredLogger
	svc    services.WalletService
}

func declareWalletRoutes(logger *zap.SugaredLogger, wsContainer *restful.Container, walletSvc services.WalletService) {
	ws := new(restful.WebService)
	defer wsContainer.Add(ws)

	routes := &walletRoutes{
		Logger: logger,
		svc:    walletSvc,
	}

	ws.Route(ws.POST("wallet").
		To(routes.CreateWallet).
		Doc("create new wallet to account").
		Consumes(restful.MIME_JSON).
		Reads(web.InputCreateWallet{}),
	)

	ws.Route(ws.GET("wallet/balance/{id}").
		Param(ws.PathParameter("id", "identifier of the account").DataType("string")).
		To(routes.BalanceWallet).
		Doc("return balance of account").
		Writes(web.OutputBalanceWalletUser{}),
	)

	ws.Route(ws.PUT("wallet/withdraw").
		To(routes.WithdrawWallet).
		Doc("perform withdraw in wallet of account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Reads(web.InputWithdrawWallet{}).
		Writes(web.OutputUpdateWalletBalance{}),
	)

	ws.Route(ws.PUT("wallet/deposit").
		To(routes.DepositWallet).
		Doc("perform deposit in wallet of account").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		Reads(web.InputDepositWallet{}).
		Writes(web.OutputUpdateWalletBalance{}),
	)
}

func (s *walletRoutes) CreateWallet(req *restful.Request, res *restful.Response) {
	input := new(web.InputCreateWallet)
	if err := parseRequestBody(req.Request.Body, input); err != nil {
		s.Logger.Errorw(ErrReadReqBody,
			zap.Any(pkg.ReqBodyKey, input), zap.String(pkg.RouteNameKey, "withdraw_wallet"), zap.Error(err))
		web.BadRequestResponse(res, ErrDeserializeReqBody)
		return
	}

	if err := s.svc.CreateWallet(req.Request.Context(), input.AccountId); err != nil {
		s.Logger.Errorw(ErrCreateWallet,
			zap.String(pkg.RouteNameKey, "create_wallet"), zap.String(pkg.AccountIDKey, input.AccountId), zap.Error(err))
		web.BadRequestResponse(res, ErrCreateWallet)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (s *walletRoutes) BalanceWallet(req *restful.Request, res *restful.Response) {
	accID := req.PathParameter("id")

	balance, err := s.svc.Balance(req.Request.Context(), accID)

	if err != nil {
		s.Logger.Errorw(ErrGetBalance,
			zap.String(pkg.RouteNameKey, "balance_wallet"), zap.String(pkg.AccountIDKey, accID), zap.Error(err))
		web.BadRequestResponse(res, ErrGetBalance)
		return
	}

	_ = res.WriteAsJson(web.OutputBalanceWalletUser{
		Balance: balance,
	})
}

func (s *walletRoutes) WithdrawWallet(req *restful.Request, res *restful.Response) {
	input := new(web.InputWithdrawWallet)
	if err := parseRequestBody(req.Request.Body, input); err != nil {
		s.Logger.Errorw(ErrReadReqBody,
			zap.Any(pkg.ReqBodyKey, input), zap.String(pkg.RouteNameKey, "withdraw_wallet"), zap.Error(err))
		web.BadRequestResponse(res, ErrReadReqBody)
		return
	}

	result, err := s.svc.WalletWithdraw(req.Request.Context(), dto.WalletWithdraw{
		AccountId: input.AccountId,
		Value:     input.Value,
	})
	if err != nil {
		s.Logger.Errorw(ErrWithdrawWallet,
			zap.String(pkg.RouteNameKey, "withdraw_wallet"), zap.String(pkg.AccountIDKey, input.AccountId), zap.Error(err))
		web.BadRequestResponse(res, ErrWithdrawWallet)
		return
	}

	res.WriteAsJson(web.OutputUpdateWalletBalance{
		OldBalance: result.OldBalance,
		NewBalance: result.NewBalance,
	})
}

func (s *walletRoutes) DepositWallet(req *restful.Request, res *restful.Response) {
	input := new(web.InputDepositWallet)
	if err := parseRequestBody(req.Request.Body, input); err != nil {
		s.Logger.Errorw(ErrReadReqBody,
			zap.Any(pkg.ReqBodyKey, input), zap.String(pkg.RouteNameKey, "deposit_wallet"), zap.Error(err))
		web.BadRequestResponse(res, ErrReadReqBody)
		return
	}

	result, err := s.svc.WalletDeposit(req.Request.Context(), dto.WalletDeposit{
		AccountId: input.AccountId,
		Value:     input.Value,
	})
	if err != nil {
		s.Logger.Errorw(ErrDepositWallet,
			zap.String(pkg.RouteNameKey, "deposit_wallet"), zap.String(pkg.AccountIDKey, input.AccountId), zap.Error(err))
		web.BadRequestResponse(res, ErrDepositWallet)
		return
	}

	res.WriteAsJson(web.OutputUpdateWalletBalance{
		OldBalance: result.OldBalance,
		NewBalance: result.NewBalance,
	})
}
