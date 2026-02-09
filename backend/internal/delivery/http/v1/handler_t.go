package v1

import (
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initTestRoutes(api *gin.RouterGroup) {
	test := api.Group("/test")
	{
		test.GET("/owner", h.TestOwner)
		test.GET("/doctor", h.TestDoctor)
		test.GET("/nurse", h.TestNurse)
	}
}

// TestOwner godoc
// @Summary Test Owner Access
// @Description Verify if user has owner role
// @Tags test
// @Response 200 {object} response.Response
// @Router /test/owner [get]
// @Security BearerAuth
func (h *Handler) TestOwner(c *gin.Context) {
	response.Success(c, codes.Ok, "You are an Owner!")
}

// TestDoctor godoc
// @Summary Test Doctor Access
// @Description Verify if user has doctor role
// @Tags test
// @Response 200 {object} response.Response
// @Router /test/doctor [get]
// @Security BearerAuth
func (h *Handler) TestDoctor(c *gin.Context) {
	response.Success(c, codes.Ok, "You are a Doctor!")
}

// TestNurse godoc
// @Summary Test Nurse Access
// @Description Verify if user has nurse role
// @Tags test
// @Response 200 {object} response.Response
// @Router /test/nurse [get]
// @Security BearerAuth
func (h *Handler) TestNurse(c *gin.Context) {
	response.Success(c, codes.Ok, "You are a Nurse!")
}
