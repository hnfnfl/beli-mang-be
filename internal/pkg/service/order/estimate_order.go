package order

import (
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var cache = util.NewCache()

func (s *OrderService) EstimateOrder(ctx *gin.Context, data dto.OrderEstimateRequest) *dto.OrderEstimateResponse {
	db := s.db
	var (
		startingMerchant dto.Orders
		stmt             strings.Builder
		calculateItems   []dto.OrdersItems

		checkMerchantIds = make([]string, 0)
		checkItemIds     = make([]string, 0)
	)

	username := ctx.Value("username").(string)
	timeStamp := time.Now().UnixNano()
	calculatedEstimateId := fmt.Sprintf("%s-%d", username, timeStamp)

	// check if the request is already cached
	if cachedResponse, exist := cache.Get(calculatedEstimateId); exist {
		return cachedResponse.(*dto.OrderEstimateResponse)
	}

	totalPrice := 0.0
	userLat, userLong := data.UserLocation.Lat, data.UserLocation.Long

	for _, order := range data.Orders {
		if order.IsStartingPoint {
			startingMerchant.MerchantId = order.MerchantId
			break
		}
	}

	for _, order := range data.Orders {
		checkMerchantIds = append(checkMerchantIds, fmt.Sprintf("'%s'", order.MerchantId))
		for _, item := range order.Items {
			checkItemIds = append(checkItemIds, fmt.Sprintf("'%s'", item.ItemId))
			calculateItems = append(calculateItems, dto.OrdersItems{
				ItemId:   item.ItemId,
				Quantity: item.Quantity,
			})
		}
	}

	checkMerchant := fmt.Sprintf("(%s)", strings.Join(checkMerchantIds, ","))
	checkItem := fmt.Sprintf("(%s)", strings.Join(checkItemIds, ","))

	stmt.WriteString(fmt.Sprintf("SELECT item_id, price FROM merchant_items WHERE merchant_id IN %s and item_id IN %s", checkMerchant, checkItem))

	rows, err := db.Query(ctx, stmt.String())
	if err != nil {
		errs.NewInternalError(ctx, "Failed to query merchant items", err)
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var item dto.OrderEstimateItemPrice
		if err := rows.Scan(&item.ItemId, &item.Price); err != nil {
			return nil
		}

		for _, calculateItem := range calculateItems {
			if calculateItem.ItemId == item.ItemId {
				totalPrice += float64(item.Price) * float64(calculateItem.Quantity)
			}
		}
	}

	stmt.Reset()
	stmt.WriteString(fmt.Sprintf("SELECT lat, long FROM (SELECT lat, long, CASE WHEN merchant_id = '%s' THEN 1 ELSE 0 END AS start_merchant FROM merchants WHERE merchant_id IN %s) AS subquery ORDER BY start_merchant DESC", startingMerchant.MerchantId, checkMerchant))

	rows, err = db.Query(ctx, stmt.String())
	if err != nil {
		errs.NewInternalError(ctx, "Failed to query merchants", err)
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
	estimatedTime := int((totalDistance / speed) * 60)

	response := dto.OrderEstimateResponse{
		TotalPrice:           totalPrice,
		EstDelivTime:         estimatedTime,
		CalculatedEstimateId: calculatedEstimateId,
	}

	// cache the response if it's not already cached
	cache.Set(calculatedEstimateId, &response, 24*time.Hour)

	return &response
}
