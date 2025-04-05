package handlers_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/cheina97/timeserver/pkg/api"
	"github.com/cheina97/timeserver/pkg/handlers"
)

var (
	router = gin.New()
)

func TestHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handlers Suite")
}

var _ = BeforeSuite(func() {
	api.RegisterHandlers(router, handlers.NewServer())
})
