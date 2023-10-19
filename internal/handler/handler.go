package handler

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"flappy-backend/internal/domain"
)

var (
	errEmptyNickname = errors.New("empty nickname")
)

type Handler struct {
	service Service
	logger  Logger
}

func NewHandler(service Service, logger Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func returnBadRequest(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func returnInternalError(err error, c *fiber.Ctx) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": err.Error(),
	})
}

func (h *Handler) AddRecord(c *fiber.Ctx) error {
	payload := addRecordIn{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	if len(payload.Nickname) == 0 {
		h.logger.Warn("empty nickname", map[string]interface{}{})
		return returnBadRequest(errEmptyNickname, c)
	}

	if payload.Score < 0 {
		h.logger.Warn("invalid score", map[string]interface{}{
			"score": payload.Score,
		})
		return returnBadRequest(fmt.Errorf("invalid score: %d", payload.Score), c)
	}

	err := h.service.AddRecord(c.Context(), domain.Record{
		Nickname: payload.Nickname,
		Score:    payload.Score,
	})
	if err != nil {
		h.logger.Error("AddRecord: error from service", map[string]interface{}{
			"err": err.Error(),
		})

		return returnInternalError(err, c)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *Handler) GetRecord(c *fiber.Ctx) error {
	payload := getRecordIn{}
	if err := c.BodyParser(&payload); err != nil {
		return returnBadRequest(err, c)
	}

	if len(payload.Nickname) == 0 {
		h.logger.Warn("empty nickname", map[string]interface{}{})
		return returnBadRequest(errEmptyNickname, c)
	}

	score, err := h.service.GetRecord(c.Context(), payload.Nickname)
	if err != nil {
		h.logger.Error("GetRecord: error from service", map[string]interface{}{
			"err": err.Error(),
		})

		return returnInternalError(err, c)
	}

	return c.JSON(getRecordOut{
		Score: score,
	})
}

func (h *Handler) GetTop10Records(c *fiber.Ctx) error {
	records, err := h.service.GetTop10Records(c.Context())
	if err != nil {
		h.logger.Error("GetTop10Records: error from service", map[string]interface{}{
			"err": err.Error(),
		})

		return returnInternalError(err, c)
	}

	out := getTop10RecordsOut{
		Records: make([]record, 0, len(records)),
	}

	for _, r := range records {
		out.Records = append(out.Records, record{
			Nickname: r.Nickname,
			Score:    r.Score,
		})
	}

	return c.JSON(out)
}
