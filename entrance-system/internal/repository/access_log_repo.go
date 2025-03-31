package repository

import (
	"context"
	"time"

	"jojihouse-entrance-system/internal/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AccessLogRepository struct {
	db *mongo.Database
}

func NewLogRepository(db *mongo.Database) *AccessLogRepository {
	return &AccessLogRepository{db: db}
}

func (r *AccessLogRepository) CreateAccessLog(log *model.AccessLog) error {
	_, err := r.db.Collection("access_log").InsertOne(context.Background(), log)
	if err != nil {
		return err
	}
	return nil
}

func (r *AccessLogRepository) CreateEntryAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "entry",
		}

	return r.CreateAccessLog(log)
}

func (r *AccessLogRepository) CreateExitAccessLog(userid int) error {
	log :=
		&model.AccessLog{
			UserID:     userid,
			Time:       time.Now(),
			AccessType: "exit",
		}

	return r.CreateAccessLog(log)
}

func (r *AccessLogRepository) GetAccessLogsByAnyFilter(lastID primitive.ObjectID, options model.AccessLogFilter) ([]model.AccessLog, error) {
	filter := bson.D{}

	// UserID のフィルタ
	if options.UserID > 0 {
		filter = append(filter, bson.E{Key: "user_id", Value: options.UserID})
	}

	// 日時フィルタ
	timeFilter := bson.D{}
	if !options.DayAfter.IsZero() {
		timeFilter = append(timeFilter, bson.E{Key: "$gte", Value: options.DayAfter})
	}
	if !options.DayBefore.IsZero() {
		timeFilter = append(timeFilter, bson.E{Key: "$lte", Value: options.DayBefore})
	}
	if len(timeFilter) > 0 {
		filter = append(filter, bson.E{Key: "time", Value: timeFilter})
	}

	// AccessType のフィルタ
	if options.AccessType != "" {
		filter = append(filter, bson.E{Key: "access_type", Value: options.AccessType})
	}

	// リミットの設定（デフォルト50）
	limit := int64(50)
	if options.Limit > 0 {
		limit = int64(options.Limit)
	}

	return r._findAccessLogs(filter, lastID, limit)
}

// 共通の検索処理
func (r *AccessLogRepository) _findAccessLogs(filter bson.D, lastID primitive.ObjectID, limit int64) ([]model.AccessLog, error) {
	var logs []model.AccessLog

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.D{{Key: "time", Value: 1}}) // `time` で昇順ソート

	// lastID によるページネーション
	if !lastID.IsZero() {
		filter = append(filter, bson.E{Key: "_id", Value: bson.D{{Key: "$gt", Value: lastID}}})
	}

	cursor, err := r.db.Collection("access_log").Find(context.Background(), filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var log model.AccessLog
		if err := cursor.Decode(&log); err != nil {
			return nil, err
		}

		// タイムゾーンの変換
		log.Time = log.Time.In(time.Local)

		logs = append(logs, log)
	}
	return logs, nil
}
