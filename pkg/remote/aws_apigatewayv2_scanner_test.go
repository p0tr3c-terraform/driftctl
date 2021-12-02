package remote

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/cloudskiff/driftctl/mocks"
	"github.com/cloudskiff/driftctl/pkg/filter"
	"github.com/cloudskiff/driftctl/pkg/remote/alerts"
	"github.com/cloudskiff/driftctl/pkg/remote/aws"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	remoteerr "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	testresource "github.com/cloudskiff/driftctl/test/resource"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiGatewayV2Api(t *testing.T) {
	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository.MockApiGatewayV2Repository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no api gateway v2 api",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return([]*apigatewayv2.Api{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "single api gateway v2 api",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return([]*apigatewayv2.Api{
					{ApiId: awssdk.String("f5vdrg12tk")},
				}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 1)

				assert.Equal(t, got[0].ResourceId(), "f5vdrg12tk")
				assert.Equal(t, got[0].ResourceType(), resourceaws.AwsApiGatewayV2ApiResourceType)
			},
		},
		{
			test: "cannot list api gateway v2 apis",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return(nil, dummyError)
				alerter.On("SendAlert", resourceaws.AwsApiGatewayV2ApiResourceType, alerts.NewRemoteAccessDeniedAlert(common.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(dummyError, resourceaws.AwsApiGatewayV2ApiResourceType, resourceaws.AwsApiGatewayV2ApiResourceType), alerts.EnumerationPhase)).Return()
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceaws.AwsApiGatewayV2ApiResourceType),
		},
	}

	providerVersion := "3.19.0"
	schemaRepository := testresource.InitFakeSchemaRepository("aws", providerVersion)
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository.MockApiGatewayV2Repository{}
			c.mocks(fakeRepo, alerter)

			var repo repository.ApiGatewayV2Repository = fakeRepo

			remoteLibrary.AddEnumerator(aws.NewApiGatewayV2ApiEnumerator(repo, factory))

			testFilter := &filter.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.wantErr)
			if err != nil {
				return
			}

			c.assertExpected(tt, got)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
		})
	}
}

func TestApiGatewayV2Route(t *testing.T) {
	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository.MockApiGatewayV2Repository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no api gateway v2 api",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return([]*apigatewayv2.Api{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "single api gateway v2 api with a single route",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return([]*apigatewayv2.Api{
					{ApiId: awssdk.String("a-gateway")},
				}, nil)
				repository.On("ListAllApiRoutes", awssdk.String("a-gateway")).
					Return([]*apigatewayv2.Route{{RouteId: awssdk.String("a-route")}}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 1)

				assert.Equal(t, "a-route", got[0].ResourceId())
				assert.Equal(t, resourceaws.AwsApiGatewayV2RouteResourceType, got[0].ResourceType())
			},
		},
		{
			test: "cannot list api gateway v2 apis",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllApis").Return(nil, dummyError)
				alerter.On("SendAlert", resourceaws.AwsApiGatewayV2RouteResourceType, alerts.NewRemoteAccessDeniedAlert(common.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(dummyError, resourceaws.AwsApiGatewayV2ApiResourceType, resourceaws.AwsApiGatewayV2RouteResourceType), alerts.EnumerationPhase)).Return()
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceaws.AwsApiGatewayV2RouteResourceType),
		},
	}

	providerVersion := "3.19.0"
	schemaRepository := testresource.InitFakeSchemaRepository("aws", providerVersion)
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository.MockApiGatewayV2Repository{}
			c.mocks(fakeRepo, alerter)

			var repo repository.ApiGatewayV2Repository = fakeRepo

			remoteLibrary.AddEnumerator(aws.NewApiGatewayV2RouteEnumerator(repo, factory))

			testFilter := &filter.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.wantErr)
			if err != nil {
				return
			}

			c.assertExpected(tt, got)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
		})
	}
}
