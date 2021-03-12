package router

import (
	"net/http"
	"strconv"

	"github.com/oscarSantoyo/academy-go-q12021/controller"
	"github.com/oscarSantoyo/academy-go-q12021/service"

	"github.com/golobby/container"
	"github.com/labstack/echo/v4"
)

type errorMessage struct {
	Message string `json:"message"`
	ErrorID int    `json:"errorId"`
}

// SetRouter Sets routes for server
func SetRouter(e *echo.Echo) {
	var configService service.ConfigService
	container.Make(&configService)
	config := configService.GetConfig()
	port := config.Server.Port

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusAccepted, "I'm Up")
	})

	e.GET("/filter", func(c echo.Context) error {
		// id := c.QueryParam("id")
		params := c.QueryParams()

		conditions := make(map[string] string)

		for k, v := range params {
			conditions[k] = v[0]
		}

		if len(conditions) == 0 {
			return reportMessage(c, "No filters were provided")
		}

		result, err := controller.SearchByConditions(conditions)
		if err != nil {
			return reportMessage(c, "Could not perform search")
		}

		if len(result) == 0 {
			return reportMessage(c, "No records found")
		}
		return c.JSON(http.StatusAccepted, result)
	})

	e.GET("/loadData", func(c echo.Context) error {
		topic := c.QueryParam("topic")
		if topic == "" {
			return reportMessage(c, "Missing parameter 'topic'")
		}
		err := controller.LoadCsvData(topic)
		if err != nil {
			return reportMessage(c, "Data was not able to load")
		}
		return c.JSON(http.StatusAccepted, "loadded correctly")
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(port)))
}

func reportMessage(c echo.Context, msg string) error {
	return c.JSON(http.StatusOK, &errorMessage{
		Message: msg,
		ErrorID: http.StatusOK,
	})
}

func reportInternalError(c echo.Context, msg string) error {
	return c.JSON(http.StatusInternalServerError, &errorMessage{
		Message: msg,
		ErrorID: http.StatusInternalServerError,
	})
}
