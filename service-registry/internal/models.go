package serviceRegistry

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type InstanceStatus string

const (
	StatusUp          InstanceStatus = "UP"
	StatusDown        InstanceStatus = "DOWN"
	StatusMaintenance InstanceStatus = "MAINTENANCE"
)

type Service struct {
	ID        bson.ObjectID `bson:"_id,omitempty"  json:"id"`
	Name      string        `bson:"name"           json:"name"`
	CreatedAt time.Time     `bson:"created_at"     json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"     json:"updated_at"`
}

type ServiceInstance struct {
	ID            bson.ObjectID  `bson:"_id,omitempty"  json:"id"`
	ServiceID     bson.ObjectID  `bson:"service_id"     json:"service_id"`
	Host          string         `bson:"host"           json:"host"`
	Port          int            `bson:"port"           json:"port"`
	Hash          uint32         `bson:"hash"           json:"hash"`
	Status        InstanceStatus `bson:"status"         json:"status"`
	LastHeartbeat time.Time      `bson:"last_heartbeat" json:"last_heartbeat"`
}
