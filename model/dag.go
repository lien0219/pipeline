package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"gorm.io/gorm"
	"time"
)

// DAGNode 表示DAG中的节点
type DAGNode struct {
	ID           string   `json:"id" gorm:"primaryKey"`
	Name         string   `json:"name"`
	Type         string   `json:"type"` // task, condition, parallel, etc.
	Config       JSONMap  `json:"config"`
	Dependencies []string `json:"dependencies"` // 依赖的节点ID列表
	Position     JSONMap  `json:"position"`     // 节点在UI中的位置
}

// JSONMap 是一个可以存储在数据库中的JSON对象
type JSONMap map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// JSONArray 是一个可以存储在数据库中的JSON数组
type JSONArray []interface{}

// Value 实现driver.Valuer接口
func (j JSONArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现sql.Scanner接口
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, j)
}

// DAG 表示有向无环图
type DAG struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	Version     int            `gorm:"default:1" json:"version"` // 版本号
	PipelineID  uint           `json:"pipeline_id"`              // 关联的流水线ID
	Pipeline    Pipeline       `gorm:"foreignKey:PipelineID" json:"pipeline"`
	Nodes       []byte         `gorm:"type:json" json:"-"` // 存储节点的JSON数据
	NodesData   []DAGNode      `gorm:"-" json:"nodes"`     // 用于JSON序列化和反序列化
	CreatorID   uint           `json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator"`
	IsActive    bool           `gorm:"default:true" json:"is_active"` // 是否为活动版本
}

// TableName 设置表名
func (DAG) TableName() string {
	return "dags"
}

// BeforeSave 在保存前将NodesData序列化为Nodes
func (d *DAG) BeforeSave(tx *gorm.DB) error {
	if d.NodesData != nil {
		data, err := json.Marshal(d.NodesData)
		if err != nil {
			return err
		}
		d.Nodes = data
	}
	return nil
}

// AfterFind 在查询后将Nodes反序列化为NodesData
func (d *DAG) AfterFind(tx *gorm.DB) error {
	if d.Nodes != nil {
		return json.Unmarshal(d.Nodes, &d.NodesData)
	}
	return nil
}
