package benchmarks

import (
	"net/http/httptest"
	"testing"

	"demo/internal/handlers"

	"github.com/gin-gonic/gin"
)

func BenchmarkGetOrder(b *testing.B) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()

	handler := handlers.NewOrderHandler(nil)
	router.GET("/order/:id", handler.GetOrderById)

	req := httptest.NewRequest("GET", "/order/test-id", nil)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)
	}
}
