package models

import "time"

// ---------- Topic ----------

type Topic struct {
	ID         string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name       string    `bson:"name" json:"name"`
	Partitions int       `bson:"partitions" json:"partitions"`
	Owner      string    `bson:"owner" json:"owner"`
	Status     string    `bson:"status" json:"status"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}

// ---------- Policy ----------

type Policy struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Principal string    `bson:"principal" json:"principal"`
	Resource  string    `bson:"resource" json:"resource"`
	Action    string    `bson:"action" json:"action"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
}

// ---------- API Requests ----------

type CreateTopicRequest struct {
	Name       string `json:"name"`
	Partitions int    `json:"partitions"`
	Owner      string `json:"owner"`
}

type CreatePolicyRequest struct {
	Principal string `json:"principal"`
	Resource  string `json:"resource"`
	Action    string `json:"action"`
}
