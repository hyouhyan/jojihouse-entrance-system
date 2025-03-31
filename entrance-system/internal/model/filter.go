package model

import "time"

// LogFilter はアクセスログのフィルタ条件を表す DTO です。
type LogFilter struct {
	UserID     int       // フィルタ対象のユーザーID
	DayBefore  time.Time // 指定日以前のログを取得
	DayAfter   time.Time // 指定日以降のログを取得
	AccessType string    // "entry" または "exit"
	Limit      int64     // 取得件数の制限
}
