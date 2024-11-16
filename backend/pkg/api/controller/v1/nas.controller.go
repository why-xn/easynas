package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/db/model"
	"github.com/whyxn/easynas/backend/pkg/dto"
	"github.com/whyxn/easynas/backend/pkg/log"
	"github.com/whyxn/easynas/backend/pkg/nas"
	"github.com/whyxn/easynas/backend/pkg/util"
	"net/http"
	"strings"
)

const DefaultPool string = "naspool"

type NasControllerInterface interface {
	GetPool(c *gin.Context)
	GetPoolList(c *gin.Context)
	GetDataset(c *gin.Context)
	GetDatasetList(c *gin.Context)
	GetDatasetFileSystem(c *gin.Context)
	CreateDataset(c *gin.Context)
	DeleteDataset(c *gin.Context)
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
