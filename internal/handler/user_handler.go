package handler

import (
	"errors"
	"strconv"

	"github.com/anish/backend-development-task/internal/models"
	"github.com/anish/backend-development-task/internal/repository"
	"github.com/anish/backend-development-task/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type UserHandler struct {
	service *service.UserService
	logger  *zap.Logger
}

func NewUserHandler(service *service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		service: service,
		logger:  logger,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad_request", "request body must be valid JSON")
	}

	resp, err := h.service.CreateUser(c.UserContext(), req)
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("user created", zap.Int32("user_id", resp.ID))
	return c.Status(fiber.StatusCreated).JSON(resp)
}

func (h *UserHandler) GetUserByID(c *fiber.Ctx) error {
	id, ok := parseID(c.Params("id"))
	if !ok {
		return writeError(c, fiber.StatusBadRequest, "bad_request", "id must be a positive integer")
	}

	resp, err := h.service.GetUserByID(c.UserContext(), id)
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("user fetched", zap.Int32("user_id", resp.ID))
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	id, ok := parseID(c.Params("id"))
	if !ok {
		return writeError(c, fiber.StatusBadRequest, "bad_request", "id must be a positive integer")
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return writeError(c, fiber.StatusBadRequest, "bad_request", "request body must be valid JSON")
	}

	resp, err := h.service.UpdateUser(c.UserContext(), id, req)
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("user updated", zap.Int32("user_id", resp.ID))
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	id, ok := parseID(c.Params("id"))
	if !ok {
		return writeError(c, fiber.StatusBadRequest, "bad_request", "id must be a positive integer")
	}

	if err := h.service.DeleteUser(c.UserContext(), id); err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("user deleted", zap.Int32("user_id", id))
	return c.Status(fiber.StatusNoContent).Send(nil)
}

func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	resp, err := h.service.ListUsers(c.UserContext())
	if err != nil {
		return h.handleError(c, err)
	}

	h.logger.Info("users listed", zap.Int("count", len(resp)))
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *UserHandler) handleError(c *fiber.Ctx, err error) error {
	var serviceErr *service.Error
	if errors.As(err, &serviceErr) {
		return writeError(c, fiber.StatusBadRequest, serviceErr.Code, serviceErr.Message)
	}

	if errors.Is(err, repository.ErrNotFound) {
		return writeError(c, fiber.StatusNotFound, "not_found", "user not found")
	}

	h.logger.Error("request failed", zap.Error(err))
	return writeError(c, fiber.StatusInternalServerError, "internal_error", "internal server error")
}

func parseID(value string) (int32, bool) {
	id, err := strconv.ParseInt(value, 10, 32)
	if err != nil || id <= 0 {
		return 0, false
	}
	return int32(id), true
}

func writeError(c *fiber.Ctx, status int, code string, message string) error {
	return c.Status(status).JSON(models.ErrorResponse{
		Error:   code,
		Message: message,
	})
}
