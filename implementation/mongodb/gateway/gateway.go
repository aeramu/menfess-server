package gateway

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func IDsFromHex(hexList []string) []primitive.ObjectID {
	var idList []primitive.ObjectID
	for _, hex := range hexList {
		id, _ := primitive.ObjectIDFromHex(hex)
		idList = append(idList, id)
	}
	return idList
}
