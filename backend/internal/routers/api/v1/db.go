package v1

import (
	"fmt"
	"log/slog"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/hexiaopi/rdm-toy/internal/model"
	"github.com/hexiaopi/rdm-toy/internal/retcode"
)

type DBController struct {
	Client ConnController
}

func (c *DBController) List(ctx *gin.Context) (res interface{}, total int64, err error) {
	client, err := c.Client.Client(ctx)
	if err != nil {
		return res, 0, err
	}
	result, err := client.Info(ctx, "keyspace").Result()
	if err != nil {
		slog.Error("redis do err", err)
	}
	Slice := strings.Split(result, "\r\n")
	dbs := make(map[string]int64)
	for _, dbInfo := range Slice[1:] {
		compileNumRegex := regexp.MustCompile(`(.*?):keys=(\d+)`)
		keyNum := compileNumRegex.FindSubmatch([]byte(dbInfo))
		if len(keyNum) != 0 {
			fmt.Println(string(keyNum[0]), string(keyNum[1]), string(keyNum[2]))
			value, _ := strconv.ParseInt(string(keyNum[2]), 10, 64)
			dbs[string(keyNum[1])] = value
		}
	}
	dbNum, err := client.ConfigGet(ctx, "databases").Result()
	if err != nil {
		return nil, 0, retcode.ListDBFail
	}
	num, _ := strconv.Atoi(dbNum["databases"])
	resp := make([]model.DB, 0)
	for i := 0; i < num; i++ {
		db := model.DB{Id: fmt.Sprintf("%d", i), Name: fmt.Sprintf("db%d", i)}
		if v, ok := dbs[db.Name]; ok {
			db.Total = v
		}
		resp = append(resp, db)
	}
	return resp, 0, nil
}

func (c *DBController) Conn(ctx *gin.Context) (*redis.Conn, error) {
	client, err := c.Client.Client(ctx)
	if err != nil {
		return nil, err
	}
	db := ctx.Param("db")
	index, err := strconv.Atoi(db)
	if err != nil {
		return nil, retcode.GetDBFail
	}
	conn := client.Conn()
	_, err = conn.Select(ctx, index).Result()
	if err != nil {
		return nil, retcode.GetDBFail
	}
	return conn, nil
}
