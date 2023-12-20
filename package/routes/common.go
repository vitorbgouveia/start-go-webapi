package routes

import (
	"encoding/json"
	"io"

	"github.com/emicklei/go-restful/v3"
	"github.com/vitorbgouveia/start-project-go/package/services"
	"github.com/vitorbgouveia/start-project-go/package/web"
	"go.uber.org/zap"
)

func DeclareAllRoutes(logger *zap.SugaredLogger, wsContainer *restful.Container, walletSvc services.WalletService) {
	declareWalletRoutes(logger, wsContainer, walletSvc)
}

func parseRequestBody(reader io.ReadCloser, input web.InputJsonRequest) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &input); err != nil {
		return err
	}

	if err := input.Validate(); err != nil {
		return err
	}

	return nil
}
