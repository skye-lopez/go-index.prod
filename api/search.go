package api

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	goquery "github.com/skye-lopez/go-query"
)

type SearchParams struct {
	Search string
	Page   int
	Limit  int
	Suffix bool
}

func Search(c *gin.Context, db *goquery.GoQuery) {
	params, code, message := ValidateAndFormParams(c)
	if code != 200 {
		c.JSON(code, gin.H{"message": message})
		return
	}

	offset := params.Page * params.Limit
	query := `SELECT url FROM packages WHERE url LIKE $1 ORDER BY url LIMIT $2 OFFSET $3`

	var searchValue string
	if params.Suffix {
		searchValue = params.Search + "%"
	} else {
		searchValue = "%" + params.Search + "%"
	}

	packages, err := db.QueryString(query, searchValue, params.Limit, offset)
	if err != nil {
		c.JSON(500, gin.H{"message": "Internal error running query."})
		return
	}

	res := []string{}
	for _, r := range packages {
		res = append(res, r.([]interface{})[0].(string))
	}

	c.JSON(200, gin.H{
		"packages": res,
		"nextPage": params.Page + 1,
	})
}

func ValidateAndFormParams(c *gin.Context) (SearchParams, int, string) {
	sp := SearchParams{}

	search := c.DefaultQuery("search", "")
	page := c.DefaultQuery("page", "0")
	limit := c.DefaultQuery("limit", "20")
	suffix := c.DefaultQuery("suffix", "false")

	parsedSuffix, err := strconv.ParseBool(suffix)
	if err != nil {
		return sp, 400, "Provided suffix option was not valid. Valid options are: [true, false]"
	}
	sp.Suffix = parsedSuffix

	parsedPage, err := strconv.Atoi(page)
	if err != nil {
		return sp, 400, fmt.Sprintf("Provided page option was not valid. Could not convert %s to an int.", page)
	}
	sp.Page = parsedPage

	parsedLimit, err := strconv.Atoi(limit)
	if err != nil {
		return sp, 400, fmt.Sprintf("Provided limit option was not valid. Could not convert %s to an int.", limit)
	}

	if parsedLimit > 2000 {
		return sp, 400, fmt.Sprintf("Provided limit option was not valid. Max limit is 2000, you provided %s", limit)
	}
	sp.Limit = parsedLimit

	sp.Search = search

	return sp, 200, ""
}
