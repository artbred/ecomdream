package versions

import (
	"context"
	"ecomdream/src/domain/models"
	"ecomdream/src/domain/replicate"
	"ecomdream/src/pkg/storages/bucket"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"time"
)

type handler struct {}

// TrainVersionHandler handler that accepts data for training
// @Description Start training process
// @Summary Start training process
// @Tags versions
// @Accept multipart/form-data
// @Produce json
// @Param class query string true "Class name"
// @Param id path string true "Payment ID"
// @Param data formData []file true "Data"
// @Success 201 {object} TrainVersionResponse
// @Router /v1/versions/train/{id} [post]
func (h *handler) TrainVersionHandler(ctx *fiber.Ctx) error {
	paymentID := ctx.Params("id"); if len(paymentID) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Provide payment id",
		})
	}

	class := ctx.Query("class"); if len(class) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Provide class",
		})
	}

	payment, err := models.GetPayment(paymentID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if payment == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid payment id",
		})
	}

	if payment.PaidAt == nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "You must pay to train model",
		})
	}

	if !payment.Plan.IsInit {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "You can only train model if you purchased initial plan",
		})
	}

	if payment.VersionID != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"code":    fiber.StatusForbidden,
			"message": "You already have model, please wait until training is finished",
		})
	}

	form, err := ctx.MultipartForm();if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	zip, err := processImagesToZip(form); if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	zipURL, err := bucket.Upload(fmt.Sprintf("%s.zip", payment.ID), zip, int64(zip.Len()), "application/zip")
	if err != nil {
		logrus.WithField("payment_id", payment.ID).Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	inputData := replicate.ConstructDreamBoothInputs(class, zipURL, replicate.TrainerVersion)
	replicateRes, err := replicate.StartDreamBoothTraining(context.Background(), inputData)
	if err != nil {
		logrus.WithField("payment_id", payment.ID).Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "We have a problem and we are already working on fixing it",
		})
	}

	version := &models.Version{
		ID:             uuid.NewV4().String(),
		PredictionID:   replicateRes.ID,
		Identifier:     replicate.UniqueIdentifier,
		Class:          class,
		InstancePrompt: inputData.Input.InstancePrompt,
		ClassPrompt:    inputData.Input.ClassPrompt,
		InstanceData:   inputData.Input.InstanceData,
		MaxTrainStep:   replicate.MaxTrainingSteps,
		Model:          &replicate.ModelName,
		TrainerVersion: replicate.TrainerVersion,
	}

	err = version.Create(payment)
	if err != nil {
		logrus.WithField("payment_id", payment.ID).Error(err)
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "We have a problem and we are already working on fixing it",
		})
	}

	logrus.Infof("Strarted training proccess for payment %s", payment.ID)

	return ctx.Status(fiber.StatusCreated).JSON(TrainVersionResponse{
		Code:      fiber.StatusCreated,
		Message:   "Successfully!",
		VersionID: version.ID,
	})
}

// VersionInfoHandler handler that sends info about version
// @Description Get info about version
// @Summary Get info about version
// @Tags versions
// @Produce json
// @Param id path string true "Version ID"
// @Success 200 {object} VersionInfoResponse
// @Router /v1/versions/info/{id} [get]
func (h *handler) VersionInfoHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id"); if len(id) == 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Provide version id",
		})
	}

	version, err := models.GetVersion(id); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	if version == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Invalid version id",
		})
	}

	if version.PushedAt == nil {
		return ctx.Status(fiber.StatusOK).JSON(VersionInfoResponse{
			Code:         fiber.StatusOK,
			IsReady:      false,
			TimeTraining: time.Now().UTC().Sub(version.CreatedAt).String(),
			Info:         nil,
		})
	}

	info, err := version.LoadExtendedInfo(); if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"code":    fiber.StatusInternalServerError,
			"message": "Please try again later",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(VersionInfoResponse{
		Code: fiber.StatusOK,
		IsReady: true,
		TimeTraining: version.PushedAt.Sub(version.CreatedAt).String(),
		Info: info,
	})
}

func createHandler() *handler {
	return &handler{}
}
