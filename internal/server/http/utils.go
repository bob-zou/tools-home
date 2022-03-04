package http

import (
	"net/http"
	"strconv"
	"tools-home/internal/server/http/common"

	"github.com/gin-gonic/gin"
)

const (
	_defaultOffset = 0
	_defaultLimit  = 20
)

func paramInt(c *gin.Context, key string, def int) (id int) {
	val := c.Param(key)
	if val == "" {
		return def
	}

	id, err := strconv.Atoi(val)
	if err != nil {
		return def
	}

	return
}

// nolint
func paramString(c *gin.Context, key string, def string) (s string) {
	s = c.Param(key)
	if s == "" {
		return def
	}

	return
}

func queryInt(c *gin.Context, key string) (vInt int) {
	val := c.Query(key)
	if val == "" {
		return
	}

	vInt, err := strconv.Atoi(val)
	if err != nil {
		return
	}

	return
}

func queryArrayInt(c *gin.Context, key string) (ids []int) {
	values := c.QueryArray(key)

	for _, v := range values {
		id, _ := strconv.Atoi(v)
		if id == 0 {
			continue
		}
		ids = append(ids, id)
	}
	return
}

func handleQueryInput(c *gin.Context) (search string, offset int, limit int) {
	var (
		queryOffset   = c.DefaultQuery("offset", "")
		queryLimit    = c.DefaultQuery("limit", "")
		queryPage     = c.DefaultQuery("pageIndex", "")
		queryPageSize = c.DefaultQuery("pageSize", "")
	)

	search = c.DefaultQuery("search", "")
	offset = _defaultOffset
	limit = _defaultLimit
	if queryOffset == "" && queryLimit == "" && queryPage == "" && queryPageSize == "" {
		return
	}

	if queryOffset != "" || queryLimit != "" {
		if queryOffset != "" {
			if tmpOffset, _ := strconv.Atoi(queryOffset); tmpOffset > 0 {
				offset = tmpOffset
			}
		}
		if queryLimit != "" {
			if tmpLimit, _ := strconv.Atoi(queryLimit); tmpLimit > 0 {
				limit = tmpLimit
			}
		}

		return
	}

	if queryPageSize != "" {
		if tmpPageSize, _ := strconv.Atoi(queryPageSize); tmpPageSize > 0 {
			limit = tmpPageSize
		}
		if limit == 0 {
			limit = _defaultLimit
		}
	}

	if queryPage != "" {
		if tmpPage, _ := strconv.Atoi(queryPage); tmpPage > 0 {
			offset = (tmpPage - 1) * limit
		}
	}

	return
}

func renderOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, common.Reply{Data: data})
}

func renderBadRequest(c *gin.Context, err error) {
	c.JSON(http.StatusBadRequest, common.Reply{Code: common.ReplyCodeErr, Msg: err.Error()})
}
