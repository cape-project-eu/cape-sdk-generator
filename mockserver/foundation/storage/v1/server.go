package v1

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"cape-project.eu/sdk-generator/provider/pulumi/secapi/models"
	"github.com/gin-gonic/gin"
)

type server struct {
	mu            sync.RWMutex
	blockStorages map[string]models.BlockStorage
}

func RegisterServer(router gin.IRouter) {
	RegisterHandlersWithOptions(router, &server{
		blockStorages: map[string]models.BlockStorage{},
	}, GinServerOptions{
		BaseURL: "/providers/seca.storage",
	})
}

func (s *server) ListImages(c *gin.Context, _tenant models.TenantPathParam, _params ListImagesParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) DeleteImage(c *gin.Context, _tenant models.TenantPathParam, _name models.ResourcePathParam, _params DeleteImageParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) GetImage(c *gin.Context, _tenant models.TenantPathParam, _name models.ResourcePathParam) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) CreateOrUpdateImage(c *gin.Context, _tenant models.TenantPathParam, _name models.ResourcePathParam, _params CreateOrUpdateImageParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) ListSkus(c *gin.Context, _tenant models.TenantPathParam, _params ListSkusParams) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) GetSku(c *gin.Context, _tenant models.TenantPathParam, _name models.ResourcePathParam) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}

func (s *server) ListBlockStorages(c *gin.Context, tenant models.TenantPathParam, workspace models.WorkspacePathParam, _params ListBlockStoragesParams) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	items := make([]models.BlockStorage, 0)
	for _, blockStorage := range s.blockStorages {
		if blockStorage.Metadata == nil {
			continue
		}
		if blockStorage.Metadata.Tenant == tenant && blockStorage.Metadata.Workspace == workspace {
			items = append(items, blockStorage)
		}
	}

	c.JSON(http.StatusOK, BlockStorageIterator{
		Items: items,
		Metadata: models.ResponseMetadata{
			Provider: "seca.storage/v1",
			Resource: fmt.Sprintf("tenants/%s/workspaces/%s/block-storages", tenant, workspace),
			Verb:     "list",
		},
	})
}

func (s *server) DeleteBlockStorage(c *gin.Context, tenant models.TenantPathParam, workspace models.WorkspacePathParam, name models.ResourcePathParam, _params DeleteBlockStorageParams) {
	s.mu.Lock()
	defer s.mu.Unlock()

	key := blockStorageKey(tenant, workspace, name)
	if _, ok := s.blockStorages[key]; !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "block-storage not found"})
		return
	}

	delete(s.blockStorages, key)
	c.JSON(http.StatusAccepted, gin.H{
		"deleted":   true,
		"tenant":    tenant,
		"workspace": workspace,
		"name":      name,
	})
}

func (s *server) GetBlockStorage(c *gin.Context, tenant models.TenantPathParam, workspace models.WorkspacePathParam, name models.ResourcePathParam) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	blockStorage, ok := s.blockStorages[blockStorageKey(tenant, workspace, name)]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "block-storage not found"})
		return
	}

	c.JSON(http.StatusOK, blockStorage)
}

func (s *server) CreateOrUpdateBlockStorage(c *gin.Context, tenant models.TenantPathParam, workspace models.WorkspacePathParam, name models.ResourcePathParam, _params CreateOrUpdateBlockStorageParams) {
	var blockStorage models.BlockStorage
	if err := c.ShouldBindJSON(&blockStorage); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now().UTC()

	s.mu.Lock()
	defer s.mu.Unlock()

	key := blockStorageKey(tenant, workspace, name)
	existing, exists := s.blockStorages[key]
	if !exists {
		blockStorage.Metadata = &models.RegionalWorkspaceResourceMetadata{
			ApiVersion:      "v1",
			CreatedAt:       now,
			Kind:            "block-storage",
			LastModifiedAt:  now,
			Name:            name,
			Provider:        "seca.storage",
			Region:          "global",
			Resource:        fmt.Sprintf("tenants/%s/workspaces/%s/block-storages/%s", tenant, workspace, name),
			ResourceVersion: 1,
			Tenant:          tenant,
			Verb:            "put",
			Workspace:       workspace,
		}
		setBlockStorageState(&blockStorage, models.ResourceStatePending)

		s.blockStorages[key] = blockStorage
		version := blockStorage.Metadata.ResourceVersion
		s.scheduleBlockStorageStateTransition(tenant, workspace, name, version, 100*time.Millisecond, models.ResourceStateCreating)
		s.scheduleBlockStorageStateTransition(tenant, workspace, name, version, 600*time.Millisecond, models.ResourceStateActive)
		c.JSON(http.StatusCreated, blockStorage)
		return
	}

	setBlockStorageState(&existing, models.ResourceStateActive)
	s.blockStorages[key] = existing

	if existing.Metadata != nil {
		blockStorage.Metadata = existing.Metadata
	} else {
		blockStorage.Metadata = &models.RegionalWorkspaceResourceMetadata{}
	}

	blockStorage.Metadata.ApiVersion = "v1"
	blockStorage.Metadata.Kind = "block-storage"
	blockStorage.Metadata.Name = name
	blockStorage.Metadata.Provider = "seca.storage"
	blockStorage.Metadata.Region = "global"
	blockStorage.Metadata.Resource = fmt.Sprintf("tenants/%s/workspaces/%s/block-storages/%s", tenant, workspace, name)
	blockStorage.Metadata.Tenant = tenant
	blockStorage.Metadata.Verb = "put"
	blockStorage.Metadata.Workspace = workspace

	if blockStorage.Metadata.CreatedAt.IsZero() {
		blockStorage.Metadata.CreatedAt = now
	}
	blockStorage.Metadata.LastModifiedAt = now
	blockStorage.Metadata.ResourceVersion++
	if blockStorage.Metadata.ResourceVersion == 0 {
		blockStorage.Metadata.ResourceVersion = 1
	}
	setBlockStorageState(&blockStorage, models.ResourceStateUpdating)

	s.blockStorages[key] = blockStorage
	version := blockStorage.Metadata.ResourceVersion
	s.scheduleBlockStorageStateTransition(tenant, workspace, name, version, 500*time.Millisecond, models.ResourceStateActive)
	c.JSON(http.StatusOK, blockStorage)
}

func (s *server) scheduleBlockStorageStateTransition(tenant models.TenantPathParam, workspace models.WorkspacePathParam, name models.ResourcePathParam, version int64, delay time.Duration, state models.ResourceState) {
	go func() {
		time.Sleep(delay)

		s.mu.Lock()
		defer s.mu.Unlock()

		key := blockStorageKey(tenant, workspace, name)
		blockStorage, ok := s.blockStorages[key]
		if !ok {
			return
		}

		if blockStorage.Metadata == nil || blockStorage.Metadata.ResourceVersion != version {
			return
		}

		setBlockStorageState(&blockStorage, state)
		s.blockStorages[key] = blockStorage
	}()
}

func setBlockStorageState(blockStorage *models.BlockStorage, state models.ResourceState) {
	if blockStorage.Status == nil {
		blockStorage.Status = &models.BlockStorageStatus{
			Conditions: []models.StatusCondition{},
			SizeGB:     blockStorage.Spec.SizeGB,
		}
	}
	if blockStorage.Status.Conditions == nil {
		blockStorage.Status.Conditions = []models.StatusCondition{}
	}
	if blockStorage.Status.State != nil && *blockStorage.Status.State == state {
		return
	}

	blockStorage.Status.SizeGB = blockStorage.Spec.SizeGB
	blockStorage.Status.State = &state

	msg := fmt.Sprintf("BlockStorage is now in %s state", state)
	reason := "stateChange"
	blockStorage.Status.Conditions = append(blockStorage.Status.Conditions, models.StatusCondition{
		LastTransitionAt: time.Now().UTC(),
		Message:          &msg,
		Reason:           &reason,
		State:            state,
	})
}

func blockStorageKey(tenant models.TenantPathParam, workspace models.WorkspacePathParam, name models.ResourcePathParam) string {
	return fmt.Sprintf("%s-%s-%s", tenant, workspace, name)
}
