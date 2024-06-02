package order

import (
	"beli-mang/internal/db/model"
	"beli-mang/internal/pkg/dto"
	"beli-mang/internal/pkg/errs"
	"beli-mang/internal/pkg/util"
	"fmt"
	"math"
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

	var merchants []model.Merchant

	for rows.Next() {
		var merchant model.Merchant
		if err := rows.Scan(
			&merchant.Location.Lat,
			&merchant.Location.Long,
		); err != nil {
			errs.NewInternalError(ctx, "Failed to scan merchants", err)
			return nil
		}
		merchants = append(merchants, merchant)
	}

	// calculate the total distance
	totalDistance, err := tspHeldKarp(
		merchants[0].Location.Lat,
		merchants[0].Location.Long,
		userLat,
		userLong,
		merchants,
	)
	if err != nil {
		errs.NewBadRequestError(ctx, "merchant too far", err)
		return nil
	}

	speed := 40.0
	estimatedTime := int((totalDistance / speed) * 60)

	response := &dto.OrderEstimateResponse{
		TotalPrice:           totalPrice,
		EstDelivTime:         estimatedTime,
		CalculatedEstimateId: calculatedEstimateId,
	}

	// cache the response if it's not already cached
	cacheData := dto.CacheItem{
		Request:  data,
		Response: *response,
		CachedAt: time.Now(),
	}
	cache.Set(calculatedEstimateId, &cacheData, 24*time.Hour)

	return response
}

func tspHeldKarp(startLat, startLon, endLat, endLon float64, merchants []model.Merchant) (float64, error) {
	n := len(merchants)
	allVisited := (1 << n) - 1
	dp := make([][]float64, n)
	for i := range dp {
		dp[i] = make([]float64, 1<<n)
		for j := range dp[i] {
			dp[i][j] = math.MaxFloat64
		}
	}
	dist := make([][]float64, n+1)
	for i := range dist {
		dist[i] = make([]float64, n+1)
	}

	// Precompute distances
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				dist[i][j] = util.Haversine(merchants[i].Location.Lat, merchants[i].Location.Long, merchants[j].Location.Lat, merchants[j].Location.Long)
			}
		}
		dist[n][i] = util.Haversine(startLat, startLon, merchants[i].Location.Lat, merchants[i].Location.Long)
		dist[i][n] = util.Haversine(merchants[i].Location.Lat, merchants[i].Location.Long, endLat, endLon)

		// check if any of the distance is more than 3km
		if dist[n][i] > 3.0 || dist[i][n] > 3.0 {
			return 0, errs.ErrMerchantTooFar
		}
	}

	var tsp func(last, visited int) float64
	tsp = func(last, visited int) float64 {
		if visited == allVisited {
			return dist[last][n] // Distance to end point
		}
		if dp[last][visited] != math.MaxFloat64 {
			return dp[last][visited]
		}
		for i := 0; i < n; i++ {
			if visited&(1<<i) == 0 {
				dp[last][visited] = math.Min(dp[last][visited], dist[last][i]+tsp(i, visited|(1<<i)))
			}
		}
		return dp[last][visited]
	}

	bestPath := []model.Merchant{}
	minDist := math.MaxFloat64

	// Compute optimal path
	for i := 0; i < n; i++ {
		currentDist := dist[n][i] + tsp(i, 1<<i)
		if currentDist < minDist {
			minDist = currentDist
			bestPath = []model.Merchant{merchants[i]}
			visited := 1 << i
			last := i
			for visited != allVisited {
				next := -1
				for j := 0; j < n; j++ {
					if visited&(1<<j) == 0 && (next == -1 || dist[last][j]+dp[j][visited|(1<<j)] < dist[last][next]+dp[next][visited|(1<<next)]) {
						next = j
					}
				}
				bestPath = append(bestPath, merchants[next])
				visited |= 1 << next
				last = next
			}
		}
	}

	bestPath = append(bestPath, model.Merchant{Location: model.Location{Lat: endLat, Long: endLon}})
	return minDist, nil
}
