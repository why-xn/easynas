package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"github.com/whyxn/easynas/backend/pkg/dto"
	"github.com/whyxn/easynas/backend/pkg/enum"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/nas"
	"github.com/whyxn/easynas/backend/pkg/util"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultPool     string = "naspool"
	DefaultClientIP string = "10.0.0.1"
)

type NasControllerInterface interface {
	GetPool(c *gin.Context)
	GetPoolList(c *gin.Context)
	GetDataset(c *gin.Context)
	GetDatasetList(c *gin.Context)
	GetDatasetFileSystem(c *gin.Context)
	CreateDataset(c *gin.Context)
	DeleteDataset(c *gin.Context)
	CreateNfsShare(c *gin.Context)
	DeleteNfsShare(c *gin.Context)
	GetNfsShareUserPermissions(c *gin.Context)
	AddUserPermissionToNfsShare(c *gin.Context)
	RemoveUserPermissionFromNfsShare(ctx *gin.Context)
	UploadFileToDataset(ctx *gin.Context)
	DeleteFileFromDataset(ctx *gin.Context)
}

type nasController struct{}

var nc nasController

func NasController() *nasController {
	return &nc
}

// GetPool
func (ctrl *nasController) GetPool(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	zpools, err := nas.ListZPools()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zpool list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	if len(zpools) > 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "success",
			"data":   zpools[0],
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   nil,
	})
}

// GetPoolList
func (ctrl *nasController) GetPoolList(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	zpools, err := nas.ListZPools()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zpool list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   zpools,
	})
}

func findDataset(dsName string) (*nas.ZFSDataset, error) {
	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"name": dsName})

	datasets, err := nas.ListZFSDatasets()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zfs datasets list", "err", err)
		return nil, err
	}

	var dataset *nas.ZFSDataset
	for _, ds := range datasets {
		if dsName == ds.Name {
			if nfsShare != nil {
				ds.ShareEnabled = true
			}
			dataset = &ds
			break
		}
	}
	return dataset, nil
}

// GetDataset
func (ctrl *nasController) GetDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	dsName := ctx.Param("dataset")
	dsName = util.Base64Decode(dsName)
	if dsName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(dsName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   dataset,
	})
}

// GetDatasetList
func (ctrl *nasController) GetDatasetList(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	nfsShareList, err := db.GetList[model.NfsShare](db.GetDb(), map[string]interface{}{"pool": DefaultPool})
	if err != nil {
		log.Logger.Errorw("Failed to fetch nfs share list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	nfsShareMap := map[string]model.NfsShare{}
	for _, nsl := range nfsShareList {
		nfsShareMap[nsl.Dataset] = nsl
	}

	datasets, err := nas.ListZFSDatasets()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zfs datasets list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var filteredDatasets []nas.ZFSDataset
	for _, ds := range datasets {
		if strings.HasPrefix(ds.Name, fmt.Sprintf("%s/", DefaultPool)) {
			if _, exists := nfsShareMap[ds.Name]; exists {
				ds.ShareEnabled = true
			}
			filteredDatasets = append(filteredDatasets, ds)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   filteredDatasets,
	})
}

// GetDatasetFileSystem
func (ctrl *nasController) GetDatasetFileSystem(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	dsName := ctx.Param("dataset")
	dsName = util.Base64Decode(dsName)
	if dsName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	path := ctx.Param("path")
	path = util.Base64Decode(path)

	dataset, err := findDataset(dsName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusNotFound)
		return
	}

	fileList, err := nas.ListAndSortFilesFolders(dataset.Name + path)
	if err != nil {
		log.Logger.Errorw("Failed to fetch filesystem", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   fileList,
	})
}

// CreateDataset
func (ctrl *nasController) CreateDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	var input dto.CreateZfsDatasetInputDTO

	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if input.Pool == "" {
		input.Pool = DefaultPool
	}

	dataset, err := findDataset(input.DatasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset != nil {
		returnErrorResponse(ctx, "dataset already exists", http.StatusBadRequest)
		return
	}

	err = nas.CreateZFSVolume(fmt.Sprintf("%s/%s", input.Pool, input.DatasetName), input.Quota)
	if err != nil {
		log.Logger.Errorw("Failed create zfs dataset", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteDataset
func (ctrl *nasController) DeleteDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	dsName := ctx.Param("dataset")
	dsName = util.Base64Decode(dsName)
	if dsName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(dsName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusNotFound)
		return
	}

	err = nas.DeleteZFSVolume(dsName)
	if err != nil {
		log.Logger.Errorw("Failed delete zfs dataset", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"name": dsName})

	if nfsShare != nil {
		if err = db.GetDb().Delete(&model.NfsShare{}, map[string]interface{}{"dataset": dsName}); err != nil {
			log.Logger.Warn("Failed delete nfs share record from db", "err", err)
		}

		if err = db.GetDb().Delete(&model.NfsSharePermission{}, map[string]interface{}{"nfsShareId": nfsShare.ID}); err != nil {
			log.Logger.Warn("Failed delete nfs share permission records from db", "err", err)
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// CreateNfsShare
func (ctrl *nasController) CreateNfsShare(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	var input dto.CreateNfsShareInputDTO

	input.Pool = ctx.Param("pool")
	if input.Pool == "" {
		input.Pool = DefaultPool
	}

	input.DatasetName = ctx.Param("dataset")
	input.DatasetName = util.Base64Decode(input.DatasetName)
	if input.DatasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(input.DatasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	err = nas.CreateNFSShare(input.DatasetName, []string{DefaultClientIP}, []string{})
	if err != nil {
		log.Logger.Errorw("Failed to create nfs share", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	nfsShare := model.NfsShare{
		Pool:    input.Pool,
		Dataset: input.DatasetName,
	}
	if err = db.GetDb().Insert(&nfsShare); err != nil {
		log.Logger.Fatalw("Failed to insert nfs share record in db", "err", err.Error())
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteNfsShare
func (ctrl *nasController) DeleteNfsShare(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	var input dto.DeleteNfsShareInputDTO

	input.Pool = ctx.Param("pool")
	if input.Pool == "" {
		input.Pool = DefaultPool
	}

	input.DatasetName = ctx.Param("dataset")
	input.DatasetName = util.Base64Decode(input.DatasetName)
	if input.DatasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(input.DatasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"dataset": input.DatasetName})
	if nfsShare == nil {
		returnErrorResponse(ctx, "nfs share not found", http.StatusBadRequest)
		return
	}

	err = nas.DeleteNFSShare(input.DatasetName)
	if err != nil {
		log.Logger.Errorw("Failed to delete nfs share", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.GetDb().Delete(&model.NfsShare{}, map[string]interface{}{"pool": input.Pool, "dataset": input.DatasetName}); err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if err := db.GetDb().Delete(&model.NfsSharePermission{}, map[string]interface{}{"nfs_share_id": 1}); err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// AddUserPermissionToNfsShare
func (ctrl *nasController) AddUserPermissionToNfsShare(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	var input dto.AddUserPermissionToNfsShareInputDTO
	err := ctx.BindJSON(&input)
	if err != nil {
		log.Logger.Errorw("Failed to bind JSON", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	input.DatasetName = ctx.Param("dataset")
	input.DatasetName = util.Base64Decode(input.DatasetName)
	if input.DatasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(input.DatasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"dataset": input.DatasetName})
	if nfsShare == nil {
		returnErrorResponse(ctx, "nfs share not found", http.StatusBadRequest)
		return
	}

	// Insert Nfs Share Permission data in db
	nfsSharePermission := model.NfsSharePermission{
		NfsShareId: nfsShare.ID,
		UserId:     input.UserId,
		Permission: input.Permission,
	}
	if err = db.GetDb().Insert(&nfsSharePermission); err != nil {
		log.Logger.Fatalw("Failed to insert nfs share permission record in db", "err", err.Error())
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch all Nfs share permissions from db
	permissionList, err := db.GetList[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"nfs_share_id": nfsShare.ID}, "NfsShare", "User")
	if err != nil {
		log.Logger.Errorw("Failed to fetch nfs share permission list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	var rPermissions []string
	var rwPermissions = []string{DefaultClientIP}

	for _, p := range permissionList {
		if p.Permission == enum.ReadOnly {
			rPermissions = append(rPermissions, p.User.NasClientIP)
		} else if p.Permission == enum.ReadWrite {
			rwPermissions = append(rwPermissions, p.User.NasClientIP)
		}
	}

	// Recreate NFS Share with updated permission
	err = nas.CreateNFSShare(input.DatasetName, rwPermissions, rwPermissions)
	if err != nil {
		log.Logger.Errorw("Failed to re-create nfs share with update permissions", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// RemoveUserPermissionFromNfsShare
func (ctrl *nasController) RemoveUserPermissionFromNfsShare(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	permissionId := ctx.Param("id")

	nfsSharePermission, _ := db.Get[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"ID": permissionId}, "NfsShare", "User")
	if nfsSharePermission == nil {
		returnErrorResponse(ctx, "nfs share permission not found", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(nfsSharePermission.NfsShare.Dataset)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	// Delete Nfs Share Permission data from db
	if err := db.GetDb().Delete(&model.NfsSharePermission{}, map[string]interface{}{"ID": permissionId}); err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	// Fetch all Nfs share permissions from db
	permissionList, err := db.GetList[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"nfs_share_id": nfsSharePermission.NfsShare.ID}, "NfsShare", "User")
	if err != nil {
		log.Logger.Errorw("Failed to fetch nfs share permission list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	var rPermissions []string
	var rwPermissions = []string{DefaultClientIP}

	for _, p := range permissionList {
		if p.Permission == enum.ReadOnly {
			rPermissions = append(rPermissions, p.User.NasClientIP)
		} else if p.Permission == enum.ReadWrite {
			rwPermissions = append(rwPermissions, p.User.NasClientIP)
		}
	}

	// Recreate NFS Share with updated permission
	err = nas.CreateNFSShare(nfsSharePermission.NfsShare.Dataset, rwPermissions, rwPermissions)
	if err != nil {
		log.Logger.Errorw("Failed to re-create nfs share with update permissions", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// GetNfsShareUserPermissions
func (ctrl *nasController) GetNfsShareUserPermissions(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	} else if !isAdmin(requester) {
		returnErrorResponse(ctx, "permission denied", http.StatusUnauthorized)
		return
	}

	pool := ctx.Param("pool")
	if pool == "" {
		pool = DefaultPool
	}

	datasetName := ctx.Param("dataset")
	datasetName = util.Base64Decode(datasetName)
	if datasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	dataset, err := findDataset(datasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"dataset": datasetName})
	if nfsShare == nil {
		returnErrorResponse(ctx, "nfs share not found", http.StatusBadRequest)
		return
	}

	// Fetch all Nfs share permissions from db
	permissionList, err := db.GetList[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"nfs_share_id": nfsShare.ID}, "NfsShare", "User")
	if err != nil {
		log.Logger.Errorw("Failed to fetch nfs share permission list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   permissionList,
	})
}

// UploadFileToDataset
func (ctrl *nasController) UploadFileToDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	datasetName := ctx.Param("dataset")
	datasetName = util.Base64Decode(datasetName)
	if datasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	relativePath := ctx.Param("path")
	relativePath = util.Base64Decode(relativePath)

	dataset, err := findDataset(datasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	if !isAdmin(requester) {
		nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"dataset": datasetName})
		if nfsShare == nil {
			returnErrorResponse(ctx, "nfs share not found", http.StatusBadRequest)
			return
		}

		// Fetch User's Nfs share permission
		userPermission, _ := db.Get[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"nfs_share_id": nfsShare.ID, "user_id": requester.ID}, "NfsShare", "User")
		if userPermission == nil {
			log.Logger.Errorw("user don't have any permission on this dataset", "err", err)
			returnErrorResponse(ctx, "you don't have any read/write permission on this dataset", http.StatusBadRequest)
			return
		}

		if userPermission.Permission != enum.ReadWrite {
			returnErrorResponse(ctx, "you don't have any write permission on this dataset", http.StatusBadRequest)
			return
		}
	}

	// Get the file from the request
	file, err := ctx.FormFile("file")
	if err != nil {
		returnErrorResponse(ctx, "file not found in the request", http.StatusBadRequest)
		return
	}

	// Specify the directory to save the uploaded file
	uploadDir := fmt.Sprintf("/%s/%s", datasetName, relativePath)

	// Create the directory if it doesn't exist
	if err = os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		log.Logger.Errorw("failed to create upload directory", err, err.Error())
		returnErrorResponse(ctx, "Failed to create upload directory", http.StatusBadRequest)
		return
	}

	// Save the file to the specified directory
	filePath := filepath.Join(uploadDir, file.Filename)
	if err = ctx.SaveUploadedFile(file, filePath); err != nil {
		log.Logger.Errorw("failed to save file", err, err.Error())
		returnErrorResponse(ctx, "Failed to save file", http.StatusBadRequest)
		return
	}

	// Respond with a success message
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}

// DeleteFileFromDataset
func (ctrl *nasController) DeleteFileFromDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	datasetName := ctx.Param("dataset")
	datasetName = util.Base64Decode(datasetName)
	if datasetName == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	relativePath := ctx.Param("path")
	relativePath = util.Base64Decode(relativePath)

	dataset, err := findDataset(datasetName)
	if err != nil {
		returnErrorResponse(ctx, err.Error(), http.StatusBadRequest)
		return
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusBadRequest)
		return
	}

	if !isAdmin(requester) {
		nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"dataset": datasetName})
		if nfsShare == nil {
			returnErrorResponse(ctx, "nfs share not found", http.StatusBadRequest)
			return
		}

		// Fetch User's Nfs share permission
		userPermission, _ := db.Get[model.NfsSharePermission](db.GetDb(), map[string]interface{}{"nfs_share_id": nfsShare.ID, "user_id": requester.ID}, "NfsShare", "User")
		if userPermission == nil {
			log.Logger.Errorw("user don't have any permission on this dataset", "err", err)
			returnErrorResponse(ctx, "you don't have any read/write permission on this dataset", http.StatusBadRequest)
			return
		}

		if userPermission.Permission != enum.ReadWrite {
			returnErrorResponse(ctx, "you don't have any write permission on this dataset", http.StatusBadRequest)
			return
		}
	}

	// Construct the full file path (relative to the base directory)
	filePath := filepath.Join(fmt.Sprintf("/%s/%s", datasetName, relativePath))

	// Check if the file exists
	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		log.Logger.Errorw("file not found", "err", err)
		returnErrorResponse(ctx, "file not found", http.StatusBadRequest)
		return
	}

	// Attempt to delete the file
	if err := os.Remove(filePath); err != nil {
		log.Logger.Errorw("failed to delete file", "err", err)
		returnErrorResponse(ctx, "failed to delete file", http.StatusBadRequest)
		return
	}

	// Respond with a success message
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
	})
}
