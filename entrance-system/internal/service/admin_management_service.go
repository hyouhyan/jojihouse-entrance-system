package service

import (
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/api/model/response"
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
	paymentLogRepository          *repository.PaymentLogRepository
}

func NewAdminManagementService(userRepository *repository.UserRepository, roleRepository *repository.RoleRepository, accessLogRepository *repository.AccessLogRepository, remainingEntriesLogRepository *repository.RemainingEntriesLogRepository, paymentLogRepository *repository.PaymentLogRepository) *AdminManagementService {
	return &AdminManagementService{userRepository: userRepository, roleRepository: roleRepository, accessLogRepository: accessLogRepository, remainingEntriesLogRepository: remainingEntriesLogRepository, paymentLogRepository: paymentLogRepository}
}

func (s *AdminManagementService) CreateUser(req *request.CreateUser) (*response.User, error) {
	// パース的な、model.userに合わせて再構築
	user := &model.User{
		Name:              req.Name,
		Description:       req.Description,
		Barcode:           req.Barcode,
		Contact:           req.Contact,
		Remaining_entries: req.Remaining_entries,
		Allergy:           req.Allergy,
		Number:            req.Number,
	}

	// ユーザーを作成
	user, err := s.userRepository.CreateUser(user)
	if err != nil {
		return nil, err
	}

	res := &response.User{
		ID:                user.ID,
		Name:              user.Name,
		Description:       user.Description,
		Barcode:           user.Barcode,
		Contact:           user.Contact,
		Remaining_entries: user.Remaining_entries,
		Registered_at:     user.Registered_at,
		Total_entries:     user.Total_entries,
		Allergy:           user.Allergy,
		Number:            user.Number,
	}

	return res, nil
}

func (s *AdminManagementService) UpdateUser(userID int, user *request.UpdateUser) error {
	userModel, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.Name != nil {
		userModel.Name = user.Name
	}
	if user.Description != nil {
		userModel.Description = user.Description
	}
	if user.Barcode != nil {
		userModel.Barcode = user.Barcode
	}
	if user.Contact != nil {
		userModel.Contact = user.Contact
	}
	if user.Remaining_entries != nil {
		userModel.Remaining_entries = user.Remaining_entries
	}
	if user.Allergy != nil {
		userModel.Allergy = user.Allergy
	}
	if user.Number != nil {
		userModel.Number = user.Number
	}

	return s.userRepository.UpdateUser(userModel)
}

func (s *AdminManagementService) DeleteUser(userID int) error {
	// ユーザーが実在するかの確認
	_, err := s.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}

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

func (s *AdminManagementService) CreatePaymentLog(log *model.PaymentLog) error {
	return s.paymentLogRepository.CreatePaymentLog(log)
}

func (s *AdminManagementService) GetAllPaymentLogs(lastID string, limit int64) ([]response.PaymentLog, error) {
	var objectID primitive.ObjectID
	var err error

	// lastIDを変換
	if lastID == "" {
		objectID = primitive.NilObjectID
	} else {
		objectID, err = primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
	}

	logs, err := s.paymentLogRepository.GetAllPaymentLogs(objectID, limit)
	if err != nil {
		return nil, err
	}

	// UserIDの一覧を作成
	userIDs := make([]int, len(logs))
	for i, log := range logs {
		userIDs[i] = log.UserID
	}

	// PostgreSQL から UserID に対応する UserName を取得
	users, err := s.userRepository.GetUsersByIDs(userIDs)
	if err != nil {
		return nil, err
	}

	// UserID -> UserName のマッピング
	userMap := make(map[int]string)
	for _, user := range users {
		userMap[*user.ID] = *user.Name
	}

	// レスポンスデータの作成
	var responseLogs []response.PaymentLog
	for _, log := range logs {
		responseLogs = append(responseLogs, response.PaymentLog{
			ID:          log.ID.Hex(),
			UserID:      log.UserID,
			UserName:    userMap[log.UserID], // UserIDからUserNameを取得
			Time:        log.Time,
			Description: log.Description,
			Amount:      log.Amount,
			Payway:      log.Payway,
		})
	}

	return responseLogs, nil
}

func (s *AdminManagementService) GetMonthlyPaymentLogs(year int, month int) (response.MonthlyPaymentLog, error) {
	logs, err := s.paymentLogRepository.GetMonthlyPaymentLogs(year, month)
	if err != nil {
		return response.MonthlyPaymentLog{}, err
	}

	// UserIDの一覧を作成
	userIDs := make([]int, len(logs))
	for i, log := range logs {
		userIDs[i] = log.UserID
	}

	// PostgreSQL から UserID に対応する UserName を取得
	users, err := s.userRepository.GetUsersByIDs(userIDs)
	if err != nil {
		return response.MonthlyPaymentLog{}, err
	}

	// UserID -> UserName のマッピング
	userMap := make(map[int]string)
	for _, user := range users {
		userMap[*user.ID] = *user.Name
	}

	// レスポンスデータの作成
	var responseLogs []response.PaymentLog
	totalAmount := 0
	for _, log := range logs {
		responseLogs = append(responseLogs, response.PaymentLog{
			ID:          log.ID.Hex(),
			UserID:      log.UserID,
			UserName:    userMap[log.UserID], // UserIDからUserNameを取得
			Time:        log.Time,
			Description: log.Description,
			Amount:      log.Amount,
			Payway:      log.Payway,
		})
		totalAmount += log.Amount
	}

	return response.MonthlyPaymentLog{
		Year:  year,
		Month: month,
		Total: totalAmount,
		Logs:  responseLogs,
	}, nil
}
