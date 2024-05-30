package util

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
	"net/url"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func UuidGenerator(prefix string, length int) string {
	const chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	randStr := make([]byte, length)
	for i := range randStr {
		randStr[i] = chars[rand.Intn(len(chars))]
	}

	return string(prefix) + string(randStr)
}

func JsonBinding(ctx *gin.Context, in interface{}) (string, error) {
	if err := ctx.ShouldBindJSON(in); err != nil {
		var errMsg string
		switch e := err.(type) {
		case *json.SyntaxError:
			errMsg = fmt.Sprintf("Invalid JSON syntax at position %d", e.Offset)
		case *json.UnmarshalTypeError:
			errMsg = fmt.Sprintf("Invalid type for JSON value: expected %s but got %s", e.Type, e.Value)
		default:
			errMsg = "JSON binding error"
		}

		return errMsg, err
	}

	return "", nil
}

func QueryBinding(ctx *gin.Context, in interface{}) (string, error) {
	if err := ctx.ShouldBindQuery(in); err != nil {
		return "Query binding error", err
	}

	return "", nil
}

func IsValidUrl(in string) bool {
	u, err := url.ParseRequestURI(in)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return false
	}

	re := regexp.MustCompile(`\.(jpg|jpeg)$`)
	return re.MatchString(u.Path)
}

func ExtractRole(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) >= 2 {
		return parts[1]
	}

	return ""
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180.0
	dLon := (lon2 - lon1) * math.Pi / 180.0
	lat1 = lat1 * math.Pi / 180.0
	lat2 = lat2 * math.Pi / 180.0

	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(lat1)*math.Cos(lat2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}
