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
)

// SubmitDataHandler handler that accepts data for training
// @Description Start training process
// @Summary Start training process
// @Tags versions
// @Accept multipart/form-data
// @Produce json
// @Param class query string true "Class name"
// @Param id query string true "Payment ID"
// @Param data formData []file true "Data"
// @Success 201 {object} SubmitDataResponse
// @Router /v1/versions/submit [post]
func SubmitDataHandler(ctx *fiber.Ctx) error {
	paymentID := ctx.Query("id")
	if len(paymentID) <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Provide payment id",
		})
	}

	class := ctx.Query("class")
	if len(class) <= 0 {
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

	form, err := ctx.MultipartForm()
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": err.Error(),
		})
	}

	zip, err := convertImagesToZip(form)
	if err != nil {
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

	inputData := replicate.ConstructDreamBoothInputs(class, zipURL)

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

	return ctx.Status(fiber.StatusCreated).JSON(SubmitDataResponse{
		Code:      fiber.StatusCreated,
		Message:   "Successfully!",
		VersionID: version.ID,
	})
}

// IsReadyHandler handler that sends status of version
// @Description Get status of version
// @Summary Get status of version
// @Tags versions
// @Produce json
// @Param version_id query string true "Version ID"
// @Success 200 {object} IsReadyResponse
// @Router /v1/versions/is-ready [get]
func IsReadyHandler(ctx *fiber.Ctx) error {
	id := ctx.Query("version_id"); if len(id) <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    fiber.StatusBadRequest,
			"message": "Provide version id",
		})
	}

	version, err := models.GetVersion(id)
	if err != nil {
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

	isReady := true; if version.PushedAt == nil {
		isReady = false
	}

	return ctx.Status(fiber.StatusOK).JSON(IsReadyResponse{
		Code:    fiber.StatusOK,
		IsReady: isReady,
	})
}
