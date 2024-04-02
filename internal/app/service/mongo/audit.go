package mongo

import "go.mongodb.org/mongo-driver/bson/primitive"

type AuditLog struct {
	ObjectID  primitive.ObjectID `bson:"_id"`
	Operation string             `json:"operation" bson:"operation"`
	RequestID string             `json:"request_id" bson:"request_id"`
}

func (log AuditLog) CollectionName() string {
	return "audit_logs"
}
