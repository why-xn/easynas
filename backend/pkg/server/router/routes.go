package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/whyxn/easynas/backend/pkg/api/controller/v1"
)

func AddApiRoutes(httpRg *gin.RouterGroup) {
	httpRg.GET("health", v1.HealthController().Check)
	httpRg.GET("health/secured", v1.HealthController().SecuredCheck)

	httpRg.POST("api/v1/auth/login", v1.AuthController().Login)

	httpRg.POST("api/v1/users", v1.UserController().Create)
	httpRg.GET("api/v1/users/:id", v1.UserController().Get)
	httpRg.GET("api/v1/users", v1.UserController().GetList)
	httpRg.DELETE("api/v1/users/:id", v1.UserController().Delete)

	httpRg.GET("api/v1/nas/pools/main", v1.NasController().GetPool)
	httpRg.GET("api/v1/nas/pools", v1.NasController().GetPoolList)

	httpRg.GET("api/v1/nas/pools/:pool/datasets/:dataset", v1.NasController().GetDataset)
	httpRg.GET("api/v1/nas/pools/:pool/datasets", v1.NasController().GetDatasetList)
	httpRg.POST("api/v1/nas/pools/:pool/datasets", v1.NasController().CreateDataset)
	httpRg.DELETE("api/v1/nas/pools/:pool/datasets/:dataset", v1.NasController().DeleteDataset)

	httpRg.POST("api/v1/nas/pools/:pool/datasets/:dataset/nfs-share", v1.NasController().CreateNfsShare)
	httpRg.DELETE("api/v1/nas/pools/:pool/datasets/:dataset/nfs-share", v1.NasController().DeleteNfsShare)
	httpRg.GET("api/v1/nas/pools/:pool/datasets/:dataset/nfs-share/permissions", v1.NasController().GetNfsShareUserPermissions)
	httpRg.POST("api/v1/nas/pools/:pool/datasets/:dataset/nfs-share/permissions", v1.NasController().AddUserPermissionToNfsShare)
	httpRg.DELETE("api/v1/nas/pools/:pool/datasets/:dataset/nfs-share/permissions/:id", v1.NasController().RemoveUserPermissionFromNfsShare)

	httpRg.GET("api/v1/nas/pools/:pool/datasets/:dataset/files/:path", v1.NasController().GetDatasetFileSystem)
	httpRg.POST("api/v1/nas/pools/:pool/datasets/:dataset/files/:path", v1.NasController().UploadFileToDataset)
	httpRg.DELETE("api/v1/nas/pools/:pool/datasets/:dataset/files/:path", v1.NasController().DeleteFileFromDataset)

	httpRg.GET("api/v1/nas/pools/:pool/datasets/:dataset/snapshots", v1.NasController().GetSnapshotList)
	httpRg.POST("api/v1/nas/pools/:pool/datasets/:dataset/snapshots", v1.NasController().CreateSnapshot)
	httpRg.POST("api/v1/nas/pools/:pool/datasets/:dataset/snapshots/restore", v1.NasController().RestoreFromSnapshot)
	httpRg.DELETE("api/v1/nas/pools/:pool/datasets/:dataset/snapshots/:snapshotName", v1.NasController().DeleteSnapshot)

	httpRg.GET("api/v1/metrics/system", v1.MetricsController().GetSystemMetrics)
}
