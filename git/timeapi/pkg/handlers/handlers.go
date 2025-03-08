package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/cheina97/timeserver/pkg/api"
	"github.com/cheina97/timeserver/pkg/service"
)

var _ api.ServerInterface = &Server{}

type Server struct{}

// GetTime implements api.ServerInterface.
func (s *Server) GetTime(c *gin.Context, params api.GetTimeParams) {
	var tz string
	if params.Timezone != nil {
		tz = *params.Timezone
	}

	t, err := service.GetTimeWithTimezone(tz)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, api.TimeResponse{CurrentTime: &t})
}

func NewServer() *Server {
	return &Server{}
}
