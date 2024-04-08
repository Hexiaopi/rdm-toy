package v1

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"github.com/hexiaopi/rdm-toy/internal/config"
	"github.com/hexiaopi/rdm-toy/internal/model"
	"github.com/hexiaopi/rdm-toy/internal/retcode"
	"github.com/hexiaopi/rdm-toy/internal/service"
)

type ConnController struct {
}

func (c *ConnController) List(ctx *gin.Context) (res interface{}, total int64, err error) {
	resp, err := service.ListConn(ctx)
	if err != nil {
		return nil, 0, retcode.ListConnFail
	}
	return resp, 0, nil
}

func (c *ConnController) Get(ctx *gin.Context) (res interface{}, err error) {
	response := make(map[string]interface{})
	client, err := c.Client(ctx)
	if err != nil {
		return res, err
	}
	//获取info
	result, err := client.Info(ctx).Result()
	if err != nil {
		slog.Error("redis info err", err)
		return nil, err
	}
	info := ParseInfo(result)
	response["Info"] = info
	//获取配置
	conf, err := client.ConfigGet(ctx, "*").Result()
	if err != nil {
		slog.Error("redis get config err", err)
		return nil, err
	}
	response["Config"] = conf
	//获取slowlog
	logs, err := client.SlowLogGet(ctx, 10).Result()
	if err != nil {
		slog.Error("redis slow log get err", err)
		return nil, err
	}
	response["Slowlog"] = logs
	return response, nil
}

func (c *ConnController) Delete(ctx *gin.Context) error {
	client, err := c.Client(ctx)
	if err != nil {
		return err
	}
	client.Close()
	conn := ctx.Param("conn")
	config.Conns.Delete(conn)
	return nil
}

func (c *ConnController) Test(ctx *gin.Context) error {
	var req model.Conn
	if err := ctx.BindJSON(&req); err != nil {
		return err
	}
	client, err := req.NewClient()
	if err != nil {
		return err
	}
	defer client.Close()
	return nil
}

func (c *ConnController) Create(ctx *gin.Context) error {
	var req model.Conn
	if err := ctx.BindJSON(&req); err != nil {
		return err
	}
	client, err := req.NewClient()
	if err != nil {
		return err
	}
	config.Conns.Store(req.Name, client)
	config.AppEngine.Redis[req.Name] = req
	return nil
}

func (c *ConnController) Client(ctx *gin.Context) (*redis.Client, error) {
	conn := ctx.Param("conn")
	connection, ok := config.Conns.Load(conn)
	if !ok {
		return nil, retcode.GetConnFail
	}
	client, ok := connection.(*redis.Client)
	if !ok {
		return nil, retcode.GetConnFail
	}
	return client, nil
}

func (c *ConnController) Clients(ctx *gin.Context) (interface{}, error) {
	conn, err := c.Client(ctx)
	if err != nil {
		return nil, err
	}
	clients := make([]model.Client, 0)
	result, err := conn.ClientList(ctx).Result()
	if err != nil {
		return nil, err
	}
	slice := strings.Split(result, "\n")
	for _, v := range slice {
		if v == "" {
			continue
		}
		var client model.Client
		info := strings.Split(v, " ")
		for _, i := range info {
			if i == "" {
				continue
			}
			kv := strings.Split(i, "=")
			if len(kv) == 2 {
				k := kv[0]
				v := kv[1]
				switch k {
				case "id":
					client.Id, _ = strconv.Atoi(v)
				case "addr":
					client.Addr = v
				case "db":
					client.DB, _ = strconv.Atoi(v)
				case "age":
					client.Age, _ = strconv.ParseInt(v, 10, 64)
				case "idle":
					client.Idle, _ = strconv.ParseInt(v, 10, 64)
				case "flags":
					client.Flags = v
				case "cmd":
					client.Cmd = v
				case "sub":
					client.Sub, _ = strconv.Atoi(v)
				case "psub":
					client.Psub, _ = strconv.Atoi(v)
				}
			}
		}
		clients = append(clients, client)
	}
	return clients, nil
}

type DeleteClientReq struct {
	Addr string `json:"addr"`
}

func (c *ConnController) DeleteClient(ctx *gin.Context) error {
	var req DeleteClientReq
	if err := ctx.BindJSON(&req); err != nil {
		return err
	}
	conn, err := c.Client(ctx)
	if err != nil {
		return err
	}
	_, err = conn.ClientKill(ctx, req.Addr).Result()
	if err != nil {
		return err
	}
	return nil
}

var upgrader = websocket.Upgrader{
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (c *ConnController) Cmd(ctx *gin.Context) {
	client, err := c.Client(ctx)
	if err != nil {
		return
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer ws.Close()

	for {
		mt, msg, err := ws.ReadMessage()
		if err != nil {
			break
		}
		slice := strings.Split(string(msg), " ")
		args := make([]interface{}, 0)
		for _, v := range slice {
			args = append(args, v)
		}
		res, err := client.Do(ctx, args...).Result()
		if err != nil {
			log.Println(err)
			_ = ws.WriteMessage(mt, []byte(err.Error()))
		}
		log.Println(res)
		var resp string
		switch resType := res.(type) {
		case []interface{}:
			for k, v := range resType {
				resp += fmt.Sprintf("%d) %q \r\n", k+1, InterfaceToString(v))
			}
		default:
			resp = InterfaceToString(resType)
		}
		_ = ws.WriteMessage(mt, []byte(resp))
	}
}

/**
* 解析redis服务信息
 */
func ParseInfo(infoString string) (result map[string]map[string]interface{}) {
	result = map[string]map[string]interface{}{}
	lines := strings.Split(infoString, "\n")
	section := "default"
	for i := range lines {

		if strings.TrimSpace(lines[i]) == "" {
			continue
		}

		if strings.HasPrefix(lines[i], "#") {
			sections := strings.Fields(lines[i])
			if len(sections) > 1 {
				section = sections[1]
			}
			continue
		}

		items := strings.Split(lines[i], ":")
		if len(items) < 1 {
			continue
		}
		if strings.Contains(items[1], ",") {
			valueItems := strings.Split(items[1], ",")
			subMap := make(map[string]interface{}, len(valueItems))
			for i := range valueItems {
				subValueitems := strings.Split(valueItems[i], "=")
				if len(subValueitems) < 1 {
					continue
				}
				floatValue, err := strconv.ParseFloat(strings.TrimSpace(subValueitems[1]), 64)
				if err == nil {
					subMap[subValueitems[0]] = floatValue
				} else {
					subMap[subValueitems[0]] = strings.TrimSpace(subValueitems[1])
				}
			}
			if _, ok := result[section]; ok {
				result[section][items[0]] = subMap
			} else {
				result[section] = map[string]interface{}{
					items[0]: subMap,
				}
			}

		} else {
			floatValue, err := strconv.ParseFloat(strings.TrimSpace(items[1]), 64)
			var vlaue interface{} = strings.TrimSpace(items[1])
			if err == nil {
				vlaue = floatValue
			}
			if _, ok := result[section]; ok {
				result[section][items[0]] = vlaue
			} else {
				result[section] = map[string]interface{}{
					items[0]: vlaue,
				}
			}
		}
	}
	return
}

func InterfaceToString(value interface{}) (s string) {
	var key string
	if value == nil {
		return key
	}

	switch value.(type) {

	case float64:
		ft := value.(float64)
		key = strconv.FormatFloat(ft, 'f', -1, 64)
	case float32:
		ft := value.(float32)
		key = strconv.FormatFloat(float64(ft), 'f', -1, 64)
	case int:
		it := value.(int)
		key = strconv.Itoa(it)
	case uint:
		it := value.(uint)
		key = strconv.Itoa(int(it))
	case int8:
		it := value.(int8)
		key = strconv.Itoa(int(it))
	case uint8:
		it := value.(uint8)
		key = strconv.Itoa(int(it))
	case int16:
		it := value.(int16)
		key = strconv.Itoa(int(it))
	case uint16:
		it := value.(uint16)
		key = strconv.Itoa(int(it))
	case int32:
		it := value.(int32)
		key = strconv.Itoa(int(it))
	case uint32:
		it := value.(uint32)
		key = strconv.Itoa(int(it))
	case int64:
		it := value.(int64)
		key = strconv.FormatInt(it, 10)
	case uint64:
		it := value.(uint64)
		key = strconv.FormatUint(it, 10)
	case string:
		key = value.(string)
	case bool:
		val, _ := value.(bool)
		if val {
			key = "True"
		} else {
			key = "False"
		}
	case []byte:
		key = string(value.([]byte))
	default:
		newValue, _ := json.Marshal(value)
		key = string(newValue)
	}

	return key
}
