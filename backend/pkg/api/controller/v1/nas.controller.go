package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/whyxn/easynas/backend/pkg/context"
	"github.com/whyxn/easynas/backend/pkg/db"
	"github.com/whyxn/easynas/backend/pkg/db/model"
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

// GetDataset
func (ctrl *nasController) GetDataset(ctx *gin.Context) {
	requester := context.GetRequesterFromContext(ctx)
	if requester == nil {
		returnErrorResponse(ctx, "unauthorized request", http.StatusUnauthorized)
		return
	}

	dataset := ctx.Param("dataset")
	dataset = util.Base64Decode(dataset)
	if dataset == "" {
		returnErrorResponse(ctx, "invalid dataset", http.StatusBadRequest)
		return
	}

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"pool": DefaultPool, "name": dataset})

	datasets, err := nas.ListZFSDatasets()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zfs datasets list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ds := range datasets {
		if dataset == ds.Name {
			if nfsShare != nil {
				ds.ShareEnabled = true
			}
			ctx.JSON(http.StatusOK, gin.H{
				"status": "success",
				"data":   ds,
			})
			return
		}
	}

	returnErrorResponse(ctx, "dataset not found", http.StatusNotFound)
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

	nfsShare, _ := db.Get[model.NfsShare](db.GetDb(), map[string]interface{}{"pool": DefaultPool, "name": dsName})

	datasets, err := nas.ListZFSDatasets()
	if err != nil {
		log.Logger.Errorw("Failed to fetch zfs datasets list", "err", err)
		returnErrorResponse(ctx, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataset = new(nas.ZFSDataset)
	for _, ds := range datasets {
		if dsName == ds.Name {
			if nfsShare != nil {
				ds.ShareEnabled = true
			}
			dataset = &ds
			break
		}
	}

	if dataset == nil {
		returnErrorResponse(ctx, "dataset not found", http.StatusNotFound)
		return
	}

	fileList, err := nas.ListAndSortFilesFolders(dataset.Name)
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
