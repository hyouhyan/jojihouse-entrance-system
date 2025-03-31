package service

import (
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminManagementService struct {
	userRepository                *repository.UserRepository
	roleRepository                *repository.RoleRepository
	accessLogRepository           *repository.AccessLogRepository
	remainingEntriesLogRepository *repository.RemainingEntriesLogRepository
}

func NewAdminManagementService(userRepository *repository.UserRepository, roleRepository *repository.RoleRepository, accessLogRepository *repository.AccessLogRepository, remainingEntriesLogRepository *repository.RemainingEntriesLogRepository) *AdminManagementService {
	return &AdminManagementService{userRepository: userRepository, roleRepository: roleRepository, accessLogRepository: accessLogRepository, remainingEntriesLogRepository: remainingEntriesLogRepository}
}

func (s *AdminManagementService) CreateUser(user *model.User) (*model.User, error) {
	return s.userRepository.CreateUser(user)
}

func (s *AdminManagementService) UpdateUser(user *model.User) error {
	return s.userRepository.UpdateUser(user)
}

func (s *AdminManagementService) DeleteUser(userID int) error {
	return s.userRepository.DeleteUser(userID)
}

// 入場可能回数の追加
func (s *AdminManagementService) IncreaseRemainingEntries(userID int, count int, reason string, updatedBy string) error {
	// 入場可能回数 追加
	beforeCount, afterCount, err := s.userRepository.IncreaseRemainingEntries(userID, count)
	if err != nil {
		return err
	}

	// ログに保存
	log := &model.RemainingEntriesLog{
		UserID:          userID,
		PreviousEntries: beforeCount,
		NewEntries:      afterCount,
		Reason:          reason,
		UpdatedBy:       updatedBy,
		UpdatedAt:       time.Now(),
	}

	return s.remainingEntriesLogRepository.CreateRemainingEntriesLog(log)
}

func (s *AdminManagementService) AddRoleToUser(userID, roleID int) error {
	return s.roleRepository.AddRoleToUser(userID, roleID)
}

func (s *AdminManagementService) RemoveRoleFromUser(userID, roleID int) error {
	return s.roleRepository.RemoveRoleFromUser(userID, roleID)
}

func (s *AdminManagementService) GetAllAccessLogs(lastID primitive.ObjectID) ([]model.AccessLog, error) {
	// フィルターをあえて指定しない
	options := model.AccessLogFilter{}

	return s.accessLogRepository.GetAccessLogsByAnyFilter(lastID, options)
}

func (s *AdminManagementService) GetRemainingEntriesLogs(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogs(lastID, 50)
}

func (s *AdminManagementService) GetRemainingEntriesLogsOnlyIncrease(lastID primitive.ObjectID) ([]model.RemainingEntriesLog, error) {
	return s.remainingEntriesLogRepository.GetRemainingEntriesLogsOnlyIncrease(lastID, 50)
}
