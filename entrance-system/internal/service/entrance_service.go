package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"log"
	"time"
)

type EntranceService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
	currentUsersRepository        *repository.CurrentUsersRepository
}

func NewEntranceService(
	userRepository *repository.UserRepository,
	roleRepository *repository.RoleRepository,
	accessLogRepository *repository.AccessLogRepository,
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository,
	currentUsersRepository *repository.CurrentUsersRepository,
) *EntranceService {
	return &EntranceService{
		userRepository:                userRepository,
		roleRepository:                roleRepository,
		accessLogRepository:           accessLogRepository,
		remainingEntriesLogRepository: remainingEntriesLogRepository,
		currentUsersRepository:        currentUsersRepository,
	}
}

// 入場したときの処理
func (s *EntranceService) EnterUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// 入場ログ作成
	err = s.accessLogRepository.CreateEntryAccessLog(user.ID)
	if err != nil {
		return err
	}

	// 在室ユーザーに追加
	err = s.currentUsersRepository.AddUserToCurrentUsers(user.ID)
	if err != nil {
		return err
	}

	isDecreaseTarget := true
	// ハウス管理者か
	isHouseAdmin, err := s.roleRepository.IsHouseAdmin(user.ID)
	if err != nil {
		log.Fatalf("Failed to check if the user is a house admin: %v", err)
	}
	if isHouseAdmin {
		isDecreaseTarget = false
	}

	// 最後に「入場可能回数を消費した」入場を取得
	lastRemainingLog, err := s.remainingEntriesLogRepository.GetLastRemainingEntriesLogByUserID(user.ID)
	if err != nil {
		return err
	}

	// ログの日が今日なら同日再入場
	if isSameDate(lastRemainingLog.UpdatedAt, time.Now()) {
		isDecreaseTarget = false
	}

	if isDecreaseTarget {
		// 残り回数を減らす
		beforeCount, afterCount, err := s.userRepository.DecreaseRemainingEntries(user.ID, 1)
		if err != nil {
			return err
		}
		// ログ保存

		log := &model.RemainingEntriesLog{
			UserID:          user.ID,
			PreviousEntries: beforeCount,
			NewEntries:      afterCount,
			Reason:          "ハウス入場のため",
			UpdatedBy:       "システム",
		}

		// ログ作成
		s.remainingEntriesLogRepository.CreateRemainingEntriesLog(log)
	}

	return nil
}

// 退場したときの処理
func (s *EntranceService) ExitUser(barcode string) error {
	// ユーザー情報を取得(存在するかの確認)
	user, err := s.userRepository.GetUserByBarcode(barcode)
	if err != nil {
		return err
	}

	// 退場ログ作成
	err = s.accessLogRepository.CreateExitAccessLog(user.ID)
	if err != nil {
		return err
	}

	// 在室ユーザーから削除
	err = s.currentUsersRepository.DeleteUserToCurrentUsers(user.ID)
	if err != nil {
		return err
	}

	// 入場回数を増やす
	err = s.userRepository.IncreaseTotalEntries(user.ID)
	if err != nil {
		return err
	}

	return nil
}

func isSameDate(a, b time.Time) bool {
	aDate := a.Truncate(24 * time.Hour)
	bDate := b.Truncate(24 * time.Hour)
	return aDate.Equal(bDate)
}
