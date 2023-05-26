package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/openinfradev/tks-api/internal/middleware/auth/request"

	"github.com/google/uuid"
	"github.com/openinfradev/tks-api/internal/repository"
	"github.com/openinfradev/tks-api/pkg/domain"
	"github.com/openinfradev/tks-api/pkg/httpErrors"
	"github.com/openinfradev/tks-api/pkg/log"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IAlertUsecase interface {
	Get(alertId uuid.UUID) (domain.Alert, error)
	GetByName(organizationId string, name string) (domain.Alert, error)
	Fetch(organizationId string) ([]domain.Alert, error)
	Create(ctx context.Context, dto domain.CreateAlertRequest) (err error)
	Update(ctx context.Context, dto domain.Alert) error
	Delete(ctx context.Context, dto domain.Alert) error

	CreateAlertAction(ctx context.Context, dto domain.AlertAction) (alertActionId uuid.UUID, err error)
}

type AlertUsecase struct {
	repo         repository.IAlertRepository
	clusterRepo  repository.IClusterRepository
	appGroupRepo repository.IAppGroupRepository
}

func NewAlertUsecase(r repository.Repository) IAlertUsecase {
	return &AlertUsecase{
		repo:         r.Alert,
		clusterRepo:  r.Cluster,
		appGroupRepo: r.AppGroup,
	}
}

func (u *AlertUsecase) Create(ctx context.Context, input domain.CreateAlertRequest) (err error) {
	if input.Alerts == nil || len(input.Alerts) == 0 {
		return fmt.Errorf("No data found")
	}

	allClusters, err := u.clusterRepo.Fetch()
	if err != nil {
		return fmt.Errorf("No clusters")
	}

	for _, alert := range input.Alerts {
		clusterId := alert.Labels.TacoCluster
		organizationId, err := u.getOrganizationFromCluster(&allClusters, clusterId)
		if err != nil {
			log.Error(err)
			continue
		}

		rawData, err := json.Marshal(alert)
		if err != nil {
			rawData = []byte{}
		}

		/*
			target := ""
			// discriminative 에 target 에 대한 정보가 있다.
			// discriminative: $labels.taco_cluster, $labels.instance
			discriminative := alert.Annotations.Discriminative
			if discriminative != "" {
				trimed := strings.TrimLeft(discriminative, " ")
				trimed = strings.TrimLeft(trimed, "$")
				arr := strings.Split(trimed, ",")

				for _, refer := range arr {

				}
			}
		*/

		grafanaUrl := ""
		appGroups, err := u.appGroupRepo.Fetch(domain.ClusterId(clusterId))
		if err == nil {
			for _, appGroup := range appGroups {
				if appGroup.AppGroupType == domain.AppGroupType_LMA {
					applications, err := u.appGroupRepo.GetApplications(appGroup.ID, domain.ApplicationType_GRAFANA)
					if err != nil {
						continue
					}
					if len(applications) > 0 {
						grafanaUrl = "http://" + applications[0].Endpoint
					}
				}
			}
		}

		node := ""
		if strings.Contains(alert.Labels.AlertName, "node") {
			node = alert.Labels.Instance
		}

		dto := domain.Alert{
			OrganizationId: organizationId,
			Name:           alert.Labels.AlertName,
			Code:           alert.Labels.AlertName,
			Grade:          alert.Labels.Severity,
			Node:           node,
			Message:        alert.Annotations.Message,
			Description:    alert.Annotations.Description,
			CheckPoint:     alert.Annotations.Checkpoint,
			Summary:        alert.Annotations.Summary,
			ClusterId:      domain.ClusterId(clusterId),
			GrafanaUrl:     grafanaUrl,
			RawData:        rawData,
		}

		_, err = u.repo.Create(dto)
		if err != nil {
			continue
		}
	}

	return nil
}

func (u *AlertUsecase) Update(ctx context.Context, dto domain.Alert) error {
	return nil
}

func (u *AlertUsecase) Get(alertId uuid.UUID) (alert domain.Alert, err error) {
	alert, err = u.repo.Get(alertId)
	if err != nil {
		return alert, err
	}
	u.makeAdditionalInfo(&alert)

	return
}

func (u *AlertUsecase) GetByName(organizationId string, name string) (out domain.Alert, err error) {
	out, err = u.repo.GetByName(organizationId, name)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return out, httpErrors.NewNotFoundError(err, "", "")
		}
		return out, err
	}
	return
}

func (u *AlertUsecase) Fetch(organizationId string) (alerts []domain.Alert, err error) {
	alerts, err = u.repo.Fetch(organizationId)
	if err != nil {
		return nil, err
	}

	for i := range alerts {
		u.makeAdditionalInfo(&alerts[i])
	}

	return alerts, nil
}

func (u *AlertUsecase) Delete(ctx context.Context, dto domain.Alert) (err error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		return httpErrors.NewUnauthorizedError(fmt.Errorf("Invalid token"), "A_INVALID_TOKEN", "")
	}

	_, err = u.Get(dto.ID)
	if err != nil {
		return httpErrors.NewNotFoundError(err, "AL_NOT_FOUND_ALERT", "")
	}

	*dto.UpdatorId = user.GetUserId()

	err = u.repo.Delete(dto)
	if err != nil {
		return err
	}

	return nil
}

func (u *AlertUsecase) CreateAlertAction(ctx context.Context, dto domain.AlertAction) (alertActionId uuid.UUID, err error) {
	user, ok := request.UserFrom(ctx)
	if !ok {
		return uuid.Nil, httpErrors.NewUnauthorizedError(fmt.Errorf("Invalid token"), "A_INVALID_TOKEN", "")
	}

	_, err = u.repo.Get(dto.AlertId)
	if err != nil {
		return uuid.Nil, httpErrors.NewBadRequestError(fmt.Errorf("Not found alert"), "AL_NOT_FOUND_ALERT", "")
	}

	userId := user.GetUserId()
	dto.TakerId = &userId
	dto.CreatedAt = time.Now()

	alertActionId, err = u.repo.CreateAlertAction(dto)
	if err != nil {
		return uuid.Nil, err
	}
	log.Info("newly created alertActionId:", alertActionId)

	return
}

func (u *AlertUsecase) getOrganizationFromCluster(clusters *[]domain.Cluster, strId string) (organizationId string, err error) {
	clusterId := domain.ClusterId(strId)
	if !clusterId.Validate() {
		return "", fmt.Errorf("Invalid clusterId %s", strId)
	}

	for _, cluster := range *clusters {
		if cluster.ID == clusterId {
			return cluster.OrganizationId, nil
		}
	}

	return "", fmt.Errorf("No martched organization %s", strId)
}

func (u *AlertUsecase) makeAdditionalInfo(alert *domain.Alert) {
	alert.FiredAt = &alert.CreatedAt
	alert.Status = domain.AlertActionStatus_CREATED

	if len(alert.AlertActions) > 0 {
		alert.TakedAt = &alert.AlertActions[0].CreatedAt
		for _, action := range alert.AlertActions {
			if action.Status == domain.AlertActionStatus_CLOSED {
				alert.ClosedAt = &action.CreatedAt
				alert.ProcessingSec = int((action.CreatedAt).Sub(alert.CreatedAt).Seconds())
			}
		}

		alert.LastTaker = alert.AlertActions[len(alert.AlertActions)-1].Taker
		alert.TakedSec = int((alert.AlertActions[0].CreatedAt).Sub(alert.CreatedAt).Seconds())
		alert.Status = alert.AlertActions[len(alert.AlertActions)-1].Status
	}
}
