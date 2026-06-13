package service

import (
	"context"
	"strings"
	"time"

	"github.com/anish/backend-development-task/internal/models"
	"github.com/anish/backend-development-task/internal/repository"
	"github.com/go-playground/validator/v10"
)

const dateLayout = "2006-01-02"

type UserService struct {
	repository repository.UserRepository
	validate   *validator.Validate
}

func NewUserService(repository repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
		validate:   validator.New(),
	}
}

func (s *UserService) CreateUser(ctx context.Context, req models.CreateUserRequest) (models.UserMutationResponse, error) {
	name, dob, err := s.validateInput(req.Name, req.DOB)
	if err != nil {
		return models.UserMutationResponse{}, err
	}

	user, err := s.repository.CreateUser(ctx, name, dob)
	if err != nil {
		return models.UserMutationResponse{}, err
	}

	return toMutationResponse(user), nil
}

func (s *UserService) GetUserByID(ctx context.Context, id int32) (models.UserResponse, error) {
	user, err := s.repository.GetUserByID(ctx, id)
	if err != nil {
		return models.UserResponse{}, err
	}

	return toUserResponse(user, time.Now().UTC()), nil
}

func (s *UserService) UpdateUser(ctx context.Context, id int32, req models.UpdateUserRequest) (models.UserMutationResponse, error) {
	name, dob, err := s.validateInput(req.Name, req.DOB)
	if err != nil {
		return models.UserMutationResponse{}, err
	}

	user, err := s.repository.UpdateUser(ctx, id, name, dob)
	if err != nil {
		return models.UserMutationResponse{}, err
	}

	return toMutationResponse(user), nil
}

func (s *UserService) DeleteUser(ctx context.Context, id int32) error {
	return s.repository.DeleteUser(ctx, id)
}

func (s *UserService) ListUsers(ctx context.Context) ([]models.UserResponse, error) {
	users, err := s.repository.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	today := time.Now().UTC()
	response := make([]models.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, toUserResponse(user, today))
	}

	return response, nil
}

func (s *UserService) validateInput(nameValue, dobValue string) (string, time.Time, error) {
	payload := struct {
		Name string `validate:"required"`
		DOB  string `validate:"required"`
	}{
		Name: nameValue,
		DOB:  dobValue,
	}

	if err := s.validate.Struct(payload); err != nil {
		return "", time.Time{}, validationError("name and dob are required")
	}

	name := strings.TrimSpace(nameValue)
	if name == "" {
		return "", time.Time{}, validationError("name is required")
	}

	dobText := strings.TrimSpace(dobValue)
	dob, err := time.Parse(dateLayout, dobText)
	if err != nil {
		return "", time.Time{}, validationError("dob must be a valid date in YYYY-MM-DD format")
	}

	if dob.After(startOfTodayUTC()) {
		return "", time.Time{}, validationError("dob cannot be in the future")
	}

	return name, dob, nil
}

func CalculateAge(dob time.Time, today time.Time) int {
	dob = dateOnlyUTC(dob)
	today = dateOnlyUTC(today)

	age := today.Year() - dob.Year()
	birthdayThisYear := time.Date(today.Year(), dob.Month(), dob.Day(), 0, 0, 0, 0, time.UTC)
	if today.Before(birthdayThisYear) {
		age--
	}
	if age < 0 {
		return 0
	}
	return age
}

func toMutationResponse(user repository.User) models.UserMutationResponse {
	return models.UserMutationResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  formatDate(user.Dob),
	}
}

func toUserResponse(user repository.User, today time.Time) models.UserResponse {
	return models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  formatDate(user.Dob),
		Age:  CalculateAge(user.Dob, today),
	}
}

func formatDate(value time.Time) string {
	return value.UTC().Format(dateLayout)
}

func startOfTodayUTC() time.Time {
	return dateOnlyUTC(time.Now().UTC())
}

func dateOnlyUTC(value time.Time) time.Time {
	value = value.UTC()
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, time.UTC)
}
