package entity

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type ETagEntity interface {
	ETagFields() []interface{}
	ETag() string
}

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
		case int:
			hasher.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int8:
			hasher.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int16:
			hasher.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int32:
			hasher.Write([]byte(strconv.FormatInt(int64(v), 10)))
		case int64:
			hasher.Write([]byte(strconv.FormatInt(v, 10)))
		case uint:
			hasher.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint8:
			hasher.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint16:
			hasher.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint32:
			hasher.Write([]byte(strconv.FormatUint(uint64(v), 10)))
		case uint64:
			hasher.Write([]byte(strconv.FormatUint(v, 10)))
		case float32:
			hasher.Write([]byte(strconv.FormatFloat(float64(v), 'f', -1, 64)))
		case float64:
			hasher.Write([]byte(strconv.FormatFloat(v, 'f', -1, 64)))
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
