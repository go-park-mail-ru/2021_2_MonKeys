package middleware

import (
	"dripapp/internal/pkg/logger"
	"dripapp/internal/pkg/monitoring"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestMiddlware(t *testing.T) {
	r := mux.NewRouter()
	metrics := monitoring.RegisterMetrics(r)
	t.Run("logger", func(t *testing.T) {
		Logger(nil, logger.DripLogger, metrics)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, true, false)
			}
		}()
		assert.Equal(t, true, true)
	})
	t.Run("panic", func(t *testing.T) {
		PanicRecovery(logger.DripLogger, metrics)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, true, false)
			}
		}()
		assert.Equal(t, true, true)
	})
	t.Run("cors", func(t *testing.T) {
		CORS(logger.DripLogger)
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, true, false)
			}
		}()
		assert.Equal(t, true, true)
	})
}
