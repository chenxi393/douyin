package database

import (
	"douyin/model"
	"douyin/package/constant"
	"douyin/storage/cache"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CreateVideo 新增视频，返回的 videoID 是为了将 videoID 放入布隆过滤器
// 这里简单的先写到数据库 后序使用redis + 布隆过滤器
func CreateVideo(video *model.Video) (uint64, error) {
	err := constant.DB.Transaction(func(tx *gorm.DB) error {
		// If value doesn't contain a matching primary key, value is inserted.
		err := tx.Create(video).Error
		if err != nil {
			return err
		}
		// TODO 这里 操作别的表了 可能得去调rpc 接口 而不是 直接操作
		// 暂时可以这样 
		cnt, err := SelectWorkCount(video.AuthorID)
		if err != nil {
			return err
		}
		// TODO 目前数据库层面仍然是耦合的 这里是否是关注分布式事务 ？？
		err = tx.Model(&model.User{}).Where("id = ?", video.AuthorID).Update("work_count", cnt+1).Error
		if err != nil {
			return err
		}
		return cache.PublishVideo(video.AuthorID, video.ID)
	})
	if err != nil {
		return 0, err
	}
	return video.ID, nil
}

func SelectVideosByUserID(userID uint64) ([]*model.Video, error) {
	videos := make([]*model.Video, 0)
	err := constant.DB.Model(&model.Video{}).Where("author_id = ? ", userID).Order("publish_time desc").Find(&videos).Error
	if err != nil {
		return nil, err
	}
	return videos, nil
}

// 根据视频ID集合查询视频信息 按id排序
// 但是这里只用在 favorite_id视频上
func SelectVideosByVideoID(videoIDList []uint64) ([]*model.Video, error) {
	res := make([]*model.Video, 0, len(videoIDList))
	err := constant.DB.Where("id IN (?)", videoIDList).Clauses(clause.OrderBy{
		Expression: clause.Expr{SQL: "FIELD(id,?)", Vars: []interface{}{videoIDList}, WithoutParentheses: true},
	}).Find(&res).Error
	// 如果想用指定顺序 可以用filed ？？ 是用field还是 业务中再去排序呢？
	if err != nil {
		return nil, err
	}
	return res, err
}

func UpdateVideoURL(playURL, coverURL string, videoID uint64) error {
	//  Don’t use Save with Model, it’s an Undefined Behavior.
	return constant.DB.Model(&model.Video{ID: videoID}).Updates(&model.Video{PlayURL: playURL, CoverURL: coverURL}).Error
}

func SelectFeedVideoList(numberVideos int, lastTime int64) ([]*model.Video, error) {
	if lastTime == 0 {
		lastTime = time.Now().UnixMilli()
	}
	res := make([]*model.Video, 0, 30)
	err := constant.DB.Model(&model.Video{}).Where("video.publish_time < ? ",
		time.UnixMilli(lastTime)).Order("video.publish_time desc").Limit(numberVideos).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SelectFeedVideoByTopic(numberVideos int, lastTime int64, topic string) ([]*model.Video, error) {
	if lastTime == 0 {
		lastTime = time.Now().UnixMilli()
	}
	res := make([]*model.Video, 0, 30)
	err := constant.DB.Model(&model.Video{}).Where("video.publish_time < ? and topic like ?",
		time.UnixMilli(lastTime), topic+"%").Order("video.publish_time desc").Limit(numberVideos).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func SearchVideoByKeyword(keyword string) ([]*model.Video, error) {
	var videos []*model.Video
	err := constant.DB.Raw("select * from video where match(title,topic) against(?) order by publish_time desc", keyword).Scan(&videos).Error
	return videos, err
}

// 已弃用 分布式或者分库分表 不宜使用join
// Scan支持的数据类型仅为struct及struct slice以及它们的指针类型
// Scan要不结构体加tag  gorm:"column:col_name" 指定列名 要不改造结构体
// func SelectFeedVideoByTopicWithJoin(numberVideos int, lastTime int64, topic string) ([]*response.VideoData, error) {
// 	if lastTime == 0 {
// 		lastTime = time.Now().UnixMilli()
// 	}
// 	res := make([]*response.VideoData, 0, 30)
// 	// 这里使用外连接 双表联查 可以考虑改多次单表 联查太麻烦
// 	err := constant.DB.Model(&model.User{}).Select(`user.*,
//     video.id as vid,
//     video.play_url,
//     video.cover_url,
//     video.favorite_count as vfavorite_count,
//     video.comment_count,
//     video.title,
// 	video.publish_time,
// 	video.topic`).Joins(
// 		"right join video on video.author_id = user.id").Where("video.publish_time < ? and video.topic like ?",
// 		time.UnixMilli(lastTime), topic+"%").Order("video.publish_time desc").Limit(numberVideos).Scan(&res).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return res, nil
// }
