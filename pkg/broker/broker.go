package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"code.cloudfoundry.org/lager"
	"github.com/pivotal-cf/brokerapi/domain"
	"github.com/pivotal-cf/brokerapi/domain/apiresponses"
)

type BrokerImpl struct {
	Logger    lager.Logger
	Config    Config
	Instances map[string]domain.GetInstanceDetailsSpec
	Bindings  map[string]domain.GetBindingSpec
}

type Config struct {
	ServiceName    string
	ServicePlan    string
	BaseGUID       string
	Credentials    interface{}
	Tags           string
	ImageURL       string
	SysLogDrainURL string
	Free           bool

	FakeAsync    bool
	FakeStateful bool
}

func NewBrokerImpl(logger lager.Logger) (bkr *BrokerImpl) {
	var credentials interface{}
	json.Unmarshal([]byte(getEnvWithDefault("CREDENTIALS", "{\"port\": \"4000\"}")), &credentials)
	fmt.Printf("Credentials: %v\n", credentials)

	return &BrokerImpl{
		Logger:    logger,
		Instances: map[string]domain.GetInstanceDetailsSpec{},
		Bindings:  map[string]domain.GetBindingSpec{},
		Config: Config{
			BaseGUID:    getEnvWithDefault("BASE_GUID", "29140B3F-0E69-4C7E-8A35"),
			ServiceName: getEnvWithDefault("SERVICE_NAME", "some-service-name"),
			ServicePlan: getEnvWithDefault("SERVICE_PLAN_NAME", "shared"),
			Credentials: credentials,
			Tags:        getEnvWithDefault("TAGS", "shared,worlds-simplest-service-broker"),
			ImageURL:    os.Getenv("IMAGE_URL"),
			Free:        true,

			FakeAsync:    os.Getenv("FAKE_ASYNC") == "true",
			FakeStateful: os.Getenv("FAKE_STATEFUL") == "true",
		},
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	if os.Getenv(key) == "" {
		return defaultValue
	}
	return os.Getenv(key)
}

func (bkr *BrokerImpl) Services(ctx context.Context) ([]domain.Service, error) {
	return []domain.Service{
		{
			ID:                   bkr.Config.BaseGUID + "-service-" + bkr.Config.ServiceName,
			Name:                 bkr.Config.ServiceName,
			Description:          "Shared service for " + bkr.Config.ServiceName,
			Bindable:             true,
			InstancesRetrievable: bkr.Config.FakeStateful,
			BindingsRetrievable:  bkr.Config.FakeStateful,
			Metadata: &domain.ServiceMetadata{
				DisplayName: bkr.Config.ServiceName,
				ImageUrl:    bkr.Config.ImageURL,
			},
			Plans: []domain.ServicePlan{
				{
					ID:          bkr.Config.BaseGUID + "-plan-" + bkr.Config.ServicePlan,
					Name:        bkr.Config.ServicePlan,
					Description: "Shared service for " + bkr.Config.ServiceName,
					Free:        &bkr.Config.Free,
				},
			},
		},
	}, nil
}

func (bkr *BrokerImpl) Provision(ctx context.Context, instanceID string, details domain.ProvisionDetails, asyncAllowed bool) (domain.ProvisionedServiceSpec, error) {
	var parameters interface{}
	json.Unmarshal(details.GetRawParameters(), &parameters)
	bkr.Instances[instanceID] = domain.GetInstanceDetailsSpec{
		ServiceID:  details.ServiceID,
		PlanID:     details.PlanID,
		Parameters: parameters,
	}
	return domain.ProvisionedServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) Deprovision(ctx context.Context, instanceID string, details domain.DeprovisionDetails, asyncAllowed bool) (domain.DeprovisionServiceSpec, error) {
	return domain.DeprovisionServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) GetInstance(ctx context.Context, instanceID string) (spec domain.GetInstanceDetailsSpec, err error) {
	if val, ok := bkr.Instances[instanceID]; ok {
		return val, nil
	}
	err = apiresponses.NewFailureResponse(fmt.Errorf("Unknown instance ID %s", instanceID), 404, "get-instance")
	return
}

func (bkr *BrokerImpl) Bind(ctx context.Context, instanceID string, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	var parameters interface{}
	json.Unmarshal(details.GetRawParameters(), &parameters)
	bkr.Bindings[bindingID] = domain.GetBindingSpec{
		Credentials: bkr.Config.Credentials,
		Parameters:  parameters,
	}
	return domain.Binding{
		Credentials: bkr.Config.Credentials,
	}, nil
}

func (bkr *BrokerImpl) Unbind(ctx context.Context, instanceID string, bindingID string, details domain.UnbindDetails, asyncAllowed bool) (domain.UnbindSpec, error) {
	return domain.UnbindSpec{}, nil
}

func (bkr *BrokerImpl) GetBinding(ctx context.Context, instanceID string, bindingID string) (spec domain.GetBindingSpec, err error) {
	if val, ok := bkr.Bindings[bindingID]; ok {
		return val, nil
	}
	err = apiresponses.NewFailureResponse(fmt.Errorf("Unknown binding ID %s", bindingID), 404, "get-binding")
	return
}

func (bkr *BrokerImpl) Update(ctx context.Context, instanceID string, details domain.UpdateDetails, asyncAllowed bool) (domain.UpdateServiceSpec, error) {
	return domain.UpdateServiceSpec{
		IsAsync: bkr.Config.FakeAsync,
	}, nil
}

func (bkr *BrokerImpl) LastOperation(ctx context.Context, instanceID string, details domain.PollDetails) (domain.LastOperation, error) {
	return domain.LastOperation{
		State: domain.Succeeded,
	}, nil
}

func (bkr *BrokerImpl) LastBindingOperation(ctx context.Context, instanceID string, bindingID string, details domain.PollDetails) (domain.LastOperation, error) {
	panic("not implemented")
}
