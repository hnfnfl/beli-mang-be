package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	cache = dto.Cache{
		Data: make(map[string]dto.CacheItem),
	}
	// cacheMutex = &sync.RWMutex{}
	cacheTTL = 24 * time.Hour // Time-to-live for cache items
)

func (s *OrderService) EstimateOrder(ctx *gin.Context, data dto.OrderEstimateRequest) *dto.OrderEstimateResponse {
	db := s.db
	var (
		startingMerchant dto.OrderEstimateRequestMerchant
		checkItem        string
		checkMerchant    string
		stmt             strings.Builder
		calculateItems   []dto.OrderEstimateRequestItem
	)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	calculatedEstimateId := string(dataJSON)

	cache.RLock()
	cachedItem, exists := cache.Data[calculatedEstimateId]
	if exists && time.Since(cachedItem.CachedAt) < cacheTTL {
		cache.Unlock()
		return &cachedItem.Response
	}
	cache.Unlock()

	totalPrice := 0.0
	userLat, userLong := data.UserLocation.Lat, data.UserLocation.Long

	for _, order := range data.Orders {
		if order.IsStartingPoint {
			startingMerchant.MerchantId = order.MerchantId
			break
		}
	}

	checkItem = "("
	checkMerchant = "("

	for _, order := range data.Orders {
		checkMerchant = fmt.Sprintf("'%s', ", order.MerchantId)
		for _, item := range order.Items {
			checkItem = fmt.Sprintf("'%s', ", item.ItemId)
			calculateItems = append(calculateItems, dto.OrderEstimateRequestItem{
				ItemId:   item.ItemId,
				Quantity: item.Quantity,
			})
		}
	}

	checkItem = string(checkItem)[:len(checkItem)-1]
	checkMerchant = string(checkMerchant[:len(checkMerchant)-1])

	checkItem = fmt.Sprintf("%s)", checkItem)
	checkMerchant = fmt.Sprintf("%s)", checkMerchant)

	stmt.WriteString(fmt.Sprintf("SELECT item_id, price FROM merchant_items WHERE merchant_id IN %s and item_id IN %s", checkMerchant, checkItem))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		return nil
	}
	defer rows.Close()

	var items []dto.OrderEstimateItemPrice
	for rows.Next() {
		var item dto.OrderEstimateItemPrice
		if err := rows.Scan(&item.ItemId, &item.Price); err != nil {
			return nil
		}
		items = append(items, item)
	}

	// Check if all item IDs exist
	if len(items) != len(calculateItems) {
		errs.NewNotFoundError(ctx, errs.ErrItemNotFound)
		return nil // Some item IDs do not exist
	}

	for _, item := range items {
		for _, calculateItem := range calculateItems {
			if calculateItem.ItemId == item.ItemId {
				totalPrice += float64(item.Price) * float64(calculateItem.Quantity)
			}
		}
	}

	stmt.Reset()
	stmt.WriteString(fmt.Sprintf("SELECT lat, long FROM (SELECT lat, long, CASE WHEN merchant_id = %s THEN 1 ELSE 0 END AS start_merchant FROM merchants WHERE merchant_id IN %s) ORDER BY start_merchant DESC", startingMerchant.MerchantId, checkMerchant))

	rows, err = db.Query(ctx, stmt.String())
	if err != nil {
		return nil
	}
	defer rows.Close()

	var totalDistance float64
	var prevLat, prevLong = userLat, userLong
	merchantExists := false

	for rows.Next() {
		var merchantLat, merchantLong float64
		if err := rows.Scan(
			&merchantLat,
			&merchantLong,
		); err != nil {
			errs.NewInternalError(ctx, "Failed to scan merchants", err)
			return nil
		}

		// Check if the merchant location is within 3 kilometers from the user's location
		userLocationDistance := util.Haversine(userLat, userLong, merchantLat, merchantLong)
		if userLocationDistance > 3 {
			errs.NewBadRequestError(ctx, fmt.Sprintf("Merchant %v, %v coordinate is too far from user location %v, %v(> 3km² in Cartesian coordinate system)", merchantLat, merchantLong, userLat, userLong), errs.ErrBadParam)
			return nil // Merchant location is too far
		}

		totalDistance += util.Haversine(prevLat, prevLong, merchantLat, merchantLong)
		prevLat, prevLong = merchantLat, merchantLong
		merchantExists = true
	}

	// Check if all merchant IDs exist
	if !merchantExists {
		errs.NewNotFoundError(ctx, errs.ErrMerchantNotFound)
		return nil // Some merchant IDs do not exist
	}

	speed := 40.0
	estimatedTime := (totalDistance / speed) * 60

	response := dto.OrderEstimateResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: estimatedTime,
		CalculatedEstimateId:           calculatedEstimateId,
	}

	cache.Lock()
	cache.Data[calculatedEstimateId] = dto.CacheItem{
		Request:  data,
		Response: response,
		CachedAt: time.Now(),
	}
	cache.Unlock()

	return &response
}
