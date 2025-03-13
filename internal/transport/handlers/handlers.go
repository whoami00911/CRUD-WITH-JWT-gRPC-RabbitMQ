package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"webPractice1/internal/domain"
	"webPractice1/internal/service"
	"webPractice1/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/whoami00911/gRPC-server/pkg/grpcPb"
)

// HandlerAssetsResponse wraps handler.AssetsResponse with a logger
type HandlerAssetsResponse struct {
	service *service.Service
	domain.AssetsResponse
	Logger *logger.Logger
}

// NewHandlerAssetsResponse initializes a new HandlerAssetsResponse
func NewHandlerAssetsResponse(log *logger.Logger, service *service.Service) *HandlerAssetsResponse {
	return &HandlerAssetsResponse{
		AssetsResponse: domain.AssetsResponse{
			Cache: make(map[*domain.AssetData]time.Time),
		},
		Logger:  log,
		service: service,
	}
}

// GetAllHandler godoc
// @Summary Get all entitys
// @Description Get all entitys from database
// @Tags CRUD
// @Produce json
// @Success 200 {array} domain.AssetData
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Security BearerAuth
// @Router /Abuseip/ [get]
func (har *HandlerAssetsResponse) GetAllHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}
	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	jsonData, err := json.Marshal(har.service.CRUDList.GetEntities())
	if err != nil {
		har.Logger.Error(fmt.Sprintf("Marshal method error: %s", err))
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, string(jsonData))
	//id, err := har.service.GetEntityById()
	if err = har.service.LogMq.Produce(grpcPb.LogItem{
		Entity: "ENTITY",
		Action: "GET",
		UserID: userID,
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}

// CreateHandler godoc
// @Summary Create a new entity
// @Description Create a new entity and store it in the database
// @Tags CRUD
// @Accept json
// @Produce json
// @Param asset body domain.AssetData true "Asset Data"
// @Success 201 {object} map[string]string "Created"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /Abuseip/ [post]
func (har *HandlerAssetsResponse) CreateHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}

	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	har.Mu.Lock()
	defer har.Mu.Unlock()

	ttl := time.Second * 15

	reqBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("ReadAll method error: %s", err))
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	var asset domain.AssetData
	if err = json.Unmarshal(reqBytes, &asset); err != nil {
		har.Logger.Error(fmt.Sprintf("Unmarshal method error: %s", err))
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	if asset.IPAddress == "" || asset.IPVersion == 0 {
		c.JSON(400, gin.H{"error": "ipAddress or ipVersion not set"})
		return
	}
	found := false
	for v := range har.Cache {
		if v.IPAddress == asset.IPAddress {
			found = true
			break
		}
	}
	if found {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	har.Cache[&asset] = time.Now().Add(ttl)

	go func() {
		time.Sleep(ttl)
		har.Mu.Lock()
		defer har.Mu.Unlock()
		delete(har.Cache, &asset)
	}()
	if err = har.service.CRUDList.AddEntity(asset); err != nil {
		har.Logger.Error(fmt.Sprintf("Add Entity Method error: %s", err))
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	c.JSON(201, gin.H{"message": "Created"})
	entityId, err := har.service.GetEntityById(asset.IPAddress)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("GetEntityById method error: %s", err))
	}
	if err = har.service.LogMq.Produce(grpcPb.LogItem{
		Entity:   "ENTITY",
		Action:   "CREATE",
		UserID:   userID,
		EntityID: int64(entityId),
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}

// UpdateHandler godoc
// @Summary Update existing entity
// @Description Update existing entity in database
// @Tags CRUD
// @Accept json
// @Produce json
// @Param asset body domain.AssetData true "Asset Data"
// @Success 201 {object} map[string]string "Updated"
// @Failure 400 {object} map[string]string "Bad Request"
// @Security BearerAuth
// @Router /Abuseip/ [put]
func (har *HandlerAssetsResponse) UpdateHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}

	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	har.Mu.Lock()
	defer har.Mu.Unlock()
	ttl := time.Second * 15
	reqBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("ReadAll method error: %s", err))
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	var asset domain.AssetData
	if err = json.Unmarshal(reqBytes, &asset); err != nil {
		har.Logger.Error(fmt.Sprintf("Unmarshal method error: %s", err))
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	if asset.IPAddress == "" || asset.IPVersion == 0 {
		c.JSON(400, gin.H{"error": "ipAddress or ipVersion not set"})
		return
	}
	find := false
	for v := range har.Cache {
		if v.IPAddress == asset.IPAddress {
			find = true
			delete(har.Cache, v)
			break
		}
	}
	if find {
		har.Cache[&asset] = time.Now().Add(ttl)
		go func() {
			time.Sleep(ttl)
			har.Mu.Lock()
			defer har.Mu.Unlock()
			delete(har.Cache, &asset)
		}()
	}
	if err = har.service.CRUDList.UpdateEntity(asset); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	c.JSON(201, gin.H{"message": "Updating"})

	entityId, err := har.service.GetEntityById(asset.IPAddress)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("GetEntityById method error: %s", err))
	}
	if err = har.service.LogMq.Produce(grpcPb.LogItem{
		Entity:   "ENTITY",
		Action:   "UPDATE",
		UserID:   userID,
		EntityID: int64(entityId),
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}

// DeleteAllHandler godoc
// @Summary Delete all entitys
// @Description Delete all entitys from cache and database
// @Tags CRUD
// @Success 200 {object} map[string]string "All entitys deleted"
// @Security BearerAuth
// @Router /Abuseip/ [delete]
func (har *HandlerAssetsResponse) DeleteAllHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}

	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	har.Mu.Lock()
	defer har.Mu.Unlock()

	for v := range har.Cache {
		delete(har.Cache, v)
	}

	har.service.CRUDList.DeleteAllEntitiesDB()
	c.JSON(200, gin.H{"message": "All assets deleted"})

	/*entityId, err := har.service.GetEntityById(asset.IPAddress)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("GetEntityById method error: %s", err))
	}*/
	if err := har.service.LogMq.Produce(grpcPb.LogItem{
		Entity: "ENTITY",
		Action: "DELETE",
		UserID: userID,
		//EntityID: int64(entityId),
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}

// DeleteHandler godoc
// @Summary Delete an entity by IP
// @Description Delete an entity by IP address from cache and database
// @Tags CRUD
// @Param ip path string true "IP Address"
// @Success 200 {object} map[string]string "entity deleted"
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /Abuseip/{ip} [delete]
func (har *HandlerAssetsResponse) DeleteHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}

	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ip := c.Param("ip")
	har.Mu.Lock()
	defer har.Mu.Unlock()

	for v := range har.Cache {
		if v.IPAddress == ip {
			delete(har.Cache, v)
			break
		}
	}

	entityId, err := har.service.GetEntityById(ip)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("GetEntityById method error: %s", err))
	}

	if err := har.service.CRUDList.DeleteEntityDB(ip); err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	c.JSON(200, gin.H{"message": "Asset deleted"})

	if err = har.service.LogMq.Produce(grpcPb.LogItem{
		Entity:   "ENTITY",
		Action:   "DELETE",
		UserID:   userID,
		EntityID: int64(entityId),
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}

// GetHandler godoc
// @Summary Get an entity by IP
// @Description Get an entity by IP address from cache or database
// @Tags CRUD
// @Param ip path string true "IP Address"
// @Produce json
// @Success 200 {object} domain.AssetData
// @Failure 404 {object} map[string]string "Not Found"
// @Security BearerAuth
// @Router /Abuseip/{ip} [get]
func (har *HandlerAssetsResponse) GetHandler(c *gin.Context) {
	id, ok := c.Get("UserId")
	if !ok {
		c.JSON(401, gin.H{"error": "bad token"})
		return
	}

	var userID int64
	switch v := id.(type) {
	case int:
		userID = int64(v)
	case int64:
		userID = v
	default:
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	ip := c.Param("ip")
	har.Mu.Lock()
	defer har.Mu.Unlock()
	for v := range har.Cache {
		if v.IPAddress == ip {
			v.IsCache = true
			v.IsDb = false
			jsonData, err := json.Marshal(v)
			if err != nil {
				har.Logger.Error(fmt.Sprintf("Marshal method error: %s", err))
				c.JSON(500, gin.H{"error": "Internal Server Error"})
				return
			}
			c.JSON(200, string(jsonData))
			return
		}
	}
	entity, err := har.service.CRUDList.GetEntity(ip)
	if entity == nil {
		c.JSON(404, gin.H{"error": "Not Found"})
		return
	}
	if err != nil {
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}
	jsonData, err := json.Marshal(entity)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("Marshal method error: %s", err))
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}
	c.JSON(200, string(jsonData))

	entityId, err := har.service.GetEntityById(ip)
	if err != nil {
		har.Logger.Error(fmt.Sprintf("GetEntityById method error: %s", err))
	}
	if err = har.service.LogMq.Produce(grpcPb.LogItem{
		Entity:   "ENTITY",
		Action:   "GET",
		UserID:   userID,
		EntityID: int64(entityId),
	}); err != nil {
		har.Logger.Error(fmt.Sprintf("Produce method error: %s", err))
	}
}
