package utils

import (
	"go.mongodb.org/mongo-driver/bson"
)

func BsonAToStructArray[T any](bsonA bson.A) (result []T, err error) {
	for i := 0; i < len(bsonA); i++ {
		var (
			item        T
			bsonItemRaw []byte
		)
		bsonItem := bsonA[i]
		bsonItemRaw, err = bson.Marshal(bsonItem)
		if err != nil {
			return result, err
		}
		if err = bson.Unmarshal(bsonItemRaw, &item); err != nil {
			return result, err
		}
		result = append(result, item)
	}

	return result, nil
}
