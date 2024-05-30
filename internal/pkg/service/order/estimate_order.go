package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func (s *OrderService) EstimateOrder(ctx *gin.Context, data dto.OrderEstimateRequest) *dto.OrderEstimateResponse {
	db := s.db
	var (
		startingMerchant dto.Orders
		checkItem        string
		checkMerchant    string
		stmt             strings.Builder
		calculateItems   []dto.OrdersItems
	)

	var (
		cache      = make(map[string]dto.CacheItem)
		cacheMutex = &sync.Mutex{}
		cacheTTL   = 24 * time.Hour // Time-to-live for cache items

	)

	dataJSON, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	calculatedEstimateId := string(dataJSON)

	cacheMutex.Lock()
	cachedItem, exists := cache[calculatedEstimateId]
	if exists && time.Since(cachedItem.CachedAt) < cacheTTL {
		cacheMutex.Unlock()
		return &cachedItem.Response
	}
	cacheMutex.Unlock()

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
			calculateItems = append(calculateItems, dto.OrdersItems{
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

	for rows.Next() {
		var merchantLat, merchantLong float64
		if err := rows.Scan(
			&merchantLat,
			&merchantLong,
		); err != nil {
			errs.NewInternalError(ctx, "Failed to scan merchants", err)
			return nil
		}
		totalDistance += util.Haversine(prevLat, prevLong, merchantLat, merchantLong)
		prevLat, prevLong = merchantLat, merchantLong
	}

	speed := 40.0
	estimatedTime := (totalDistance / speed) * 60

	response := dto.OrderEstimateResponse{
		TotalPrice:                     totalPrice,
		EstimatedDeliveryTimeInMinutes: estimatedTime,
		CalculatedEstimateId:           calculatedEstimateId,
	}

	cacheMutex.Lock()
	cache[calculatedEstimateId] = dto.CacheItem{
		Request:  data,
		Response: response,
		CachedAt: time.Now(),
	}
	cacheMutex.Unlock()

	return &response
}
