// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"

	"gorm.io/gorm"
)

const TableNameCommentMetum = "comment_meta"

// CommentMetum mapped from table <comment_meta>
type CommentMetum struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	CommentID int64          `gorm:"column:comment_id;not null;comment:评论ID" json:"comment_id"`                  // 评论ID
	ItemID    int64          `gorm:"column:item_id;not null;comment:视频ID，文章ID等 抽象的物品item id" json:"item_id"`     // 视频ID，文章ID等 抽象的物品item id
	ParentID  int64          `gorm:"column:parent_id;not null;comment:0:根评论 非0:子评论" json:"parent_id"`            // 0:根评论 非0:子评论
	UserID    int64          `gorm:"column:user_id;not null;comment:评论的用户ID" json:"user_id"`                     // 评论的用户ID
	Status    int32          `gorm:"column:status;not null;comment:1:全都可见 2:已删除 3:置顶 4:审核不通过 。。。" json:"status"` // 1:全都可见 2:已删除 3:置顶 4:审核不通过 。。。
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;not null" json:"deleted_at"`
	CreatedAt time.Time      `gorm:"column:created_at;not null" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not null" json:"updated_at"`
}

// TableName CommentMetum's table name
func (*CommentMetum) TableName() string {
	return TableNameCommentMetum
}
