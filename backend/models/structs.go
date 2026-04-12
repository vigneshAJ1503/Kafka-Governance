package models

import "time"

type TopicStatus string

const (
	TopicPending  TopicStatus = "PENDING"
	TopicApproved TopicStatus = "APPROVED"
)

type Topic struct {
	ID         string      `bson:"_id,omitempty" json:"id"`
	Name       string      `bson:"name" json:"name"`
	Cluster    string      `bson:"cluster" json:"cluster"`
	Partitions int         `bson:"partitions" json:"partitions"`
	Replicas   int         `bson:"replicas" json:"replicas"`
	Status     TopicStatus `bson:"status" json:"status"`
	RequestedBy string     `bson:"requestedBy" json:"requestedBy"`
	ApprovedBy  string     `bson:"approvedBy,omitempty" json:"approvedBy,omitempty"`
	CreatedAt  time.Time   `bson:"createdAt" json:"createdAt"`
	ApprovedAt *time.Time  `bson:"approvedAt,omitempty" json:"approvedAt,omitempty"`
}

type Policy struct {
	ID        string    `bson:"_id,omitempty" json:"id"`
	Principal string    `bson:"principal" json:"principal"` // User::"u_123"
	Action    string    `bson:"action" json:"action"`       // Action::"CreateTopic"
	Resource  string    `bson:"resource" json:"resource"`   // Topic::"orders.created"
	Effect    string    `bson:"effect" json:"effect"`        // permit / forbid
	CreatedAt time.Time `bson:"createdAt" json:"createdAt"`
}

