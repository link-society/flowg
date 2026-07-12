package operations

import (
	"context"
	"errors"
	"log/slog"

	"net/http"

	"go.uber.org/fx"

	"github.com/swaggest/openapi-go"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"

	"link-society.com/flowg/api/auth"
	"link-society.com/flowg/api/logging"
	"link-society.com/flowg/api/routing"
	"link-society.com/flowg/internal/models"

	"link-society.com/flowg/internal/storage/generic/kv"
	storage "link-society.com/flowg/internal/storage/interfaces"
)

// BackupConfigDeps lists the dependencies of [NewBackupConfigUsecase].
type BackupConfigDeps struct {
	fx.In

	AuthStorage   storage.AuthStorage
	ConfigStorage storage.ConfigStorage
}

// BackupConfigRequest is empty: the whole configuration database is exported.
type BackupConfigRequest struct{}

// BackupConfigResponse streams the configuration database snapshot to the
// client.
//
// It embeds the writer so the snapshot is streamed as a file download rather
// than buffered in memory.
type BackupConfigResponse struct {
	usecase.OutputWithEmbeddedWriter
}

// NewBackupConfigUsecase streams a full snapshot of the configuration database
// (pipelines, transformers and forwarders) as a downloadable file.
//
// It is the export half of the configuration backup story; the snapshot can
// later be reloaded with [NewRestoreConfigUsecase]. Callers must have
// read access to pipelines, transformers and forwarders.
func NewBackupConfigUsecase(deps BackupConfigDeps) usecase.Interactor {
	logger := logging.Logger()

	u := usecase.NewInteractor(
		auth.RequireScopesApiDecorator(
			deps.AuthStorage,
			[]models.Scope{
				models.SCOPE_READ_PIPELINES,
				models.SCOPE_READ_TRANSFORMERS,
				models.SCOPE_READ_FORWARDERS,
			},
			func(
				ctx context.Context,
				req BackupConfigRequest,
				resp *BackupConfigResponse,
			) error {
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Type", "application/octet-stream")
				resp.Writer.(http.ResponseWriter).Header().Set("Content-Disposition", "attachment; filename=config.db")
				resp.Writer.(http.ResponseWriter).Header().Set("Cache-Control", "no-cache")

				_, err := deps.ConfigStorage.Dump(ctx, resp.Writer, 0)
				if err != nil {
					if errors.Is(err, kv.ErrNotSupported) {
						return status.Wrap(err, status.Unimplemented)
					}

					logger.ErrorContext(
						ctx,
						"Failed to backup configuration database",
						slog.String("error", err.Error()),
					)

					return status.Wrap(err, status.Internal)
				}

				resp.Writer.(http.Flusher).Flush()

				return nil
			},
		),
	)

	u.SetName("backup_config")
	u.SetTitle("Backup Configuration")
	u.SetDescription("Download a full snapshot of the configuration database.")
	u.SetTags("backup")

	u.SetExpectedErrors(status.Unauthenticated, status.PermissionDenied, status.Unimplemented, status.Internal)

	return u
}

// annotateBackupConfig documents the backup response as a binary file download.
func annotateBackupConfig(oc openapi.OperationContext) error {
	contentUnits := oc.Response()
	for i, cu := range contentUnits {
		if cu.HTTPStatus == 200 {
			cu.ContentType = "application/octet-stream"
			cu.Description = "Binary file"
			cu.Format = "Binary file"
		}

		contentUnits[i] = cu
	}

	return nil
}

func init() {
	routing.RegisterOperation(
		NewBackupConfigUsecase,
		http.MethodGet,
		"/api/v1/backup/config",
		routing.Annotated(annotateBackupConfig),
	)
}
