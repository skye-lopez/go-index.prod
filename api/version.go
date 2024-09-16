package api

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	goquery "github.com/skye-lopez/go-query"
)

type Pkg struct {
	Pkg      string `json:"pkg"`
	Versions []struct {
		Version string `json:"version"`
		Time    string `json:"time"`
	} `json:"versions"`
}

func (p *Pkg) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("Could not type assert query response to []byte")
	}

	return json.Unmarshal(b, &p)
}

func Versions(c *gin.Context, db *goquery.GoQuery) {
	search := c.DefaultQuery("package", "")
	if search == "" {
		c.JSON(400, gin.H{"message": "Please provide a package= search parameter."})
		return
	}

	query := `
        SELECT
        JSONB_BUILD_OBJECT( 'pkg', owner,
        'versions', JSONB_AGG(JSONB_BUILD_OBJECT(
            'version', version,
            'time', time
        )))
        FROM package_versions
        WHERE owner = $1
        GROUP BY owner
    `

	resp := &Pkg{}
	err := db.Conn.QueryRow(query, search).Scan(&resp)
	if err != nil {
		fmt.Println(err)
		c.JSON(500, gin.H{"message": "Internal error performing query."})
		return
	}

	c.JSON(200, gin.H{"package": resp})
}
