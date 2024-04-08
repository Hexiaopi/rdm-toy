package v1

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/hexiaopi/rdm-toy/internal/model"
	"github.com/hexiaopi/rdm-toy/internal/retcode"
)

type KeyController struct {
	Conn DBController
}

func (c *KeyController) List(ctx *gin.Context) (res interface{}, total int64, err error) {
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	filter := ctx.Query("filter")
	if filter == "" {
		filter = "*"
	}
	keys, _, err := conn.Scan(ctx, 0, filter, 100).Result()
	if err != nil {
		return nil, 0, retcode.ListKeyFail
	}
	resp := make([]model.Key, 0)
	for _, key := range keys {
		resp = append(resp, model.Key{Name: key})
	}
	return resp, int64(len(keys)), nil
}

func (c *KeyController) DeleteKeys(ctx *gin.Context) error {
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err = conn.FlushDB(ctx).Result()
	if err != nil {
		return retcode.DeleteKeysFail
	}
	return nil
}

func (c *KeyController) Get(ctx *gin.Context) (res interface{}, err error) {
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return
	}
	defer conn.Close()

	key := ctx.Param("key")
	keyType, err := conn.Type(ctx, key).Result()
	if err != nil {
		return nil, retcode.GetKeyFail
	}

	ttl, err := conn.TTL(ctx, key).Result()
	if err != nil {
		return nil, retcode.GetKeyFail
	}
	resp := model.Key{
		Name: key,
		Type: keyType,
		TTL:  int64(ttl.Seconds()),
	}

	switch keyType {
	case "string":
		resp.Value, err = conn.Get(ctx, key).Result()
	case "hash":
		resp.Value, err = conn.HGetAll(ctx, key).Result()
	case "list":
		llen, _ := conn.LLen(ctx, key).Result()
		resp.Value, err = conn.LRange(ctx, key, 0, llen).Result()
	case "set":
		slen, _ := conn.SCard(ctx, key).Result()
		resp.Value, _, err = conn.SScan(ctx, key, 0, "", slen).Result()
	case "zset":
		zlen, _ := conn.ZCard(ctx, key).Result()
		var result []string
		result, _, err = conn.ZScan(ctx, key, 0, "", zlen).Result()
		value := make([]model.ZSetValue, 0)
		for i := 0; i < len(result); i += 2 {
			score, _ := strconv.ParseFloat(result[i+1], 64)
			value = append(value, model.ZSetValue{Member: result[i], Score: score})
		}
		resp.Value = value
	case "stream":
		resp.Value, err = conn.XRange(ctx, key, "-", "+").Result()
	}
	if err != nil {
		return nil, retcode.GetKeyFail
	}
	return resp, nil
}

type KeyCreateReq struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	TTL   int64       `json:"ttl"`
	Value interface{} `json:"value"`
}

func (c *KeyController) Create(ctx *gin.Context) error {
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	var req KeyCreateReq
	if err = ctx.ShouldBind(&req); err != nil {
		return retcode.CreateKeyFail
	}

	exist, err := conn.Exists(ctx, req.Name).Result()
	if err != nil {
		return retcode.CreateKeyFail
	}
	if exist == 1 {
		return retcode.ExistKeyAlready
	}

	ttl := time.Duration(req.TTL * int64(time.Second))
	switch req.Type {
	case "string":
		_, err = conn.Set(ctx, req.Name, req.Value, ttl).Result()
	case "hash":
		value := req.Value.([]interface{})
		fields := make([]interface{}, 0, len(value)*2)
		for _, v := range value {
			f := v.(map[string]interface{})
			fields = append(fields, f["field"])
			fields = append(fields, f["value"])
		}
		_, err = conn.HSet(ctx, req.Name, fields).Result()
	case "list":
		_, err = conn.LPush(ctx, req.Name, req.Value).Result()
	case "set":
		_, err = conn.SAdd(ctx, req.Name, req.Value).Result()
	case "zset":
		value := req.Value.([]interface{})
		members := make([]redis.Z, 0, len(value))
		for _, v := range value {
			f := v.(map[string]interface{})
			score := f["score"]
			sc, _ := strconv.ParseFloat(score.(string), 64)
			member := f["member"]
			members = append(members, redis.Z{
				Score:  sc,
				Member: member,
			})
		}
		_, err = conn.ZAdd(ctx, req.Name, members...).Result()
	case "stream":
		value := req.Value.(map[string]interface{})
		id := value["id"].(string)
		slice := value["value"].([]interface{})
		fields := make([]interface{}, 0, len(slice)*2)
		for _, v := range slice {
			field := v.(map[string]interface{})
			fields = append(fields, field["field"])
			fields = append(fields, field["value"])
		}
		_, err = conn.XAdd(ctx, &redis.XAddArgs{Stream: req.Name, ID: id, Values: fields}).Result()
	}
	if err != nil {
		return retcode.CreateKeyFail
	}
	if ttl > 0 {
		_, err = conn.Expire(ctx, req.Name, ttl).Result()
		if err != nil {
			return retcode.CreateKeyFail
		}
	}
	return nil
}

func (c *KeyController) Delete(ctx *gin.Context) error {
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := ctx.Param("key")
	if _, err = conn.Del(ctx, key).Result(); err != nil {
		return retcode.DeleteKeyFail
	}
	return nil
}

type CreateFieldReq struct {
	Way   string      `json:"way,omitempty"`
	Field string      `json:"field,omitempty"`
	Id    string      `json:"id"`
	Value interface{} `json:"value"`
}

func (c *KeyController) CreateField(ctx *gin.Context) error {
	var req CreateFieldReq
	if err := ctx.ShouldBind(&req); err != nil {
		return retcode.CreateFieldFail
	}
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := ctx.Param("key")
	keyType, err := conn.Type(ctx, key).Result()
	if err != nil {
		return retcode.GetKeyFail
	}
	switch keyType {
	case "hash":
		fields := make([]interface{}, 0, 2)
		fields = append(fields, req.Field)
		fields = append(fields, req.Value)
		_, err = conn.HSet(ctx, key, fields...).Result()
	case "list":
		switch req.Way {
		case "left":
			_, err = conn.LPushX(ctx, key, req.Value).Result()
		case "right":
			_, err = conn.RPushX(ctx, key, req.Value).Result()
		default:
			_, err = conn.RPushX(ctx, key, req.Value).Result()
		}
	case "set":
		_, err = conn.SAdd(ctx, key, req.Value).Result()
	case "zset":
		score := req.Value
		sc, _ := strconv.ParseFloat(score.(string), 64)
		member := req.Field
		_, err = conn.ZAdd(ctx, key, redis.Z{Score: sc, Member: member}).Result()
	case "stream":
		id := req.Id
		values := make(map[string]interface{}, 0)
		values[req.Field] = req.Value
		_, err = conn.XAdd(ctx, &redis.XAddArgs{Stream: key, ID: id, Values: values}).Result()
	default:
		return retcode.NotSupportKeyType
	}
	if err != nil {
		return retcode.CreateFieldFail
	}
	return nil
}

type PatchFieldReq struct {
	Field string      `json:"field"`
	Index int64       `json:"index"`
	Value interface{} `json:"value"`
}

func (c *KeyController) PatchField(ctx *gin.Context) error {
	var req PatchFieldReq
	if err := ctx.ShouldBind(&req); err != nil {
		return retcode.PatchFieldFail
	}
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := ctx.Param("key")
	keyType, err := conn.Type(ctx, key).Result()
	if err != nil {
		return retcode.GetKeyFail
	}
	switch keyType {
	case "hash":
		fields := make([]interface{}, 0, 2)
		fields = append(fields, req.Field)
		fields = append(fields, req.Value)
		_, err = conn.HSet(ctx, key, fields...).Result()
	case "list":
		_, err = conn.LSet(ctx, key, req.Index, req.Value).Result()
	case "set":
		_, err = conn.SRem(ctx, key, req.Field).Result()
		if err != nil {
			return retcode.PatchFieldFail
		}
		_, err = conn.SAdd(ctx, key, req.Value).Result()
	case "zset":
		_, err = conn.ZRem(ctx, key, req.Field).Result()
		if err != nil {
			return retcode.PatchFieldFail
		}
		score := req.Value
		sc, _ := strconv.ParseFloat(score.(string), 64)
		member := req.Field
		_, err = conn.ZAdd(ctx, key, redis.Z{Score: sc, Member: member}).Result()
	default:
		return retcode.NotSupportKeyType
	}
	if err != nil {
		return retcode.PatchFieldFail
	}
	return nil
}

type DeleteFieldReq struct {
	Fields []string `json:"fields"`
}

func (c *KeyController) DeleteField(ctx *gin.Context) error {
	var req DeleteFieldReq
	if err := ctx.ShouldBind(&req); err != nil {
		return retcode.DeleteFieldFail
	}
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	key := ctx.Param("key")
	keyType, err := conn.Type(ctx, key).Result()
	if err != nil {
		return retcode.GetKeyFail
	}
	switch keyType {
	case "hash":
		_, err = conn.HDel(ctx, key, req.Fields...).Result()
	case "list":
		for _, field := range req.Fields {
			_, err = conn.LRem(ctx, key, 0, field).Result()
			if err != nil {
				return retcode.DeleteFieldFail
			}
		}
	case "set":
		fields := make([]interface{}, 0, len(req.Fields))
		for _, field := range req.Fields {
			fields = append(fields, field)
		}
		_, err = conn.SRem(ctx, key, fields...).Result()
	case "zset":
		fields := make([]interface{}, 0, len(req.Fields))
		for _, field := range req.Fields {
			fields = append(fields, field)
		}
		_, err = conn.ZRem(ctx, key, fields...).Result()
	case "stream":
		_, err = conn.XDel(ctx, key, req.Fields...).Result()
	default:
		return retcode.NotSupportKeyType
	}
	if err != nil {
		return retcode.DeleteFieldFail
	}
	return nil
}

type PatchNameReq struct {
	Name string `json:"name"`
}

func (c *KeyController) PatchName(ctx *gin.Context) error {
	key := ctx.Param("key")
	var req PatchNameReq
	if err := ctx.BindJSON(&req); err != nil {
		return retcode.UpdateKeyFail
	}
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	result, err := conn.Exists(ctx, req.Name).Result()
	if err != nil {
		return retcode.UpdateKeyFail
	}
	if result != 0 {
		return retcode.ExistKeyAlready
	}
	if _, err = conn.Rename(ctx, key, req.Name).Result(); err != nil {
		return retcode.UpdateKeyFail
	}
	return nil
}

type PatchTTLReq struct {
	TTL int64 `json:"ttl"`
}

func (c *KeyController) PatchTTL(ctx *gin.Context) error {
	key := ctx.Param("key")
	var req PatchTTLReq
	if err := ctx.BindJSON(&req); err != nil {
		return retcode.UpdateKeyFail
	}
	conn, err := c.Conn.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if req.TTL < 0 {
		_, err = conn.Persist(ctx, key).Result()
		if err != nil {
			return err
		}
		return nil
	}
	_, err = conn.Expire(ctx, key, time.Duration(req.TTL*int64(time.Second))).Result()
	return err
}
