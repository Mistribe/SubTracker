package entity

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

func CalculateETag(etagEntities ...ETagEntity) string {
	hasher := sha256.New()
	var etagFields []interface{}
	if len(etagEntities) > 0 {
		for _, etagEntity := range etagEntities {
			etagFields = append(etagFields, etagEntity.ETagFields()...)
		}
	}

	for _, value := range etagFields {
		switch v := value.(type) {
		case string:
			hasher.Write([]byte(v))
		case int, int8, int16, int32, int64:
			hasher.Write([]byte(strconv.FormatInt(v.(int64), 10)))
		case uint, uint8, uint16, uint32, uint64:
			hasher.Write([]byte(strconv.FormatUint(v.(uint64), 10)))
		case float32, float64:
			hasher.Write([]byte(strconv.FormatFloat(v.(float64), 'f', -1, 64)))
		case bool:
			if v {
				hasher.Write([]byte("1"))
			} else {
				hasher.Write([]byte("0"))
			}
		case time.Time:
			hasher.Write([]byte(v.String()))
		case nil:
			// Skip nil values
		default:
			hasher.Write([]byte(fmt.Sprintf("%v", v)))
		}
	}
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
