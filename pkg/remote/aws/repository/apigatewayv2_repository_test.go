package repository

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	awstest "github.com/cloudskiff/driftctl/test/aws"
	"github.com/pkg/errors"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

func Test_apigatewayv2Repository_ListAllApis(t *testing.T) {
	apis := []*apigatewayv2.Api{
		{ApiId: aws.String("api1")},
		{ApiId: aws.String("api2")},
		{ApiId: aws.String("api3")},
		{ApiId: aws.String("api4")},
		{ApiId: aws.String("api5")},
		{ApiId: aws.String("api6")},
	}

	remoteError := errors.New("remote error")

	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache)
		want    []*apigatewayv2.Api
		wantErr error
	}{
		{
			name: "list multiple apis",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetApis",
					&apigatewayv2.GetApisInput{}).Return(&apigatewayv2.GetApisOutput{Items: apis}, nil).Once()

				store.On("Get", "apigatewayv2ListAllApis").Return(nil).Times(1)
				store.On("Put", "apigatewayv2ListAllApis", apis).Return(false).Times(1)
			},
			want: apis,
		},
		{
			name: "should hit cache",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				store.On("Get", "apigatewayv2ListAllApis").Return(apis).Times(1)
			},
			want: apis,
		},
		{
			name: "should return remote error",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetApis",
					&apigatewayv2.GetApisInput{}).Return(nil, remoteError).Once()

				store.On("Get", "apigatewayv2ListAllApis").Return(nil).Times(1)
			},
			wantErr: remoteError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &cache.MockCache{}
			client := &awstest.MockFakeApiGatewayV2{}
			tt.mocks(client, store)
			r := &apigatewayv2Repository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllApis()
			assert.Equal(t, tt.wantErr, err)

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %s -> %s", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
			store.AssertExpectations(t)
			client.AssertExpectations(t)
		})
	}
}

func Test_apigatewayv2Repository_ListAllApiRoutes(t *testing.T) {
	routes := []*apigatewayv2.Route{
		{RouteId: aws.String("route1")},
		{RouteId: aws.String("route2")},
		{RouteId: aws.String("route3")},
	}

	remoteError := errors.New("remote error")

	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache)
		want    []*apigatewayv2.Route
		wantErr error
	}{
		{
			name: "list multiple routes",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetRoutes",
					&apigatewayv2.GetRoutesInput{ApiId: aws.String("an-id")}).
					Return(&apigatewayv2.GetRoutesOutput{Items: routes}, nil).Once()

				store.On("Get", "apigatewayv2ListAllApiRoutes_api_an-id").Return(nil).Times(1)
				store.On("Put", "apigatewayv2ListAllApiRoutes_api_an-id", routes).Return(false).Times(1)
			},
			want: routes,
		},
		{
			name: "should hit cache",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				store.On("Get", "apigatewayv2ListAllApiRoutes_api_an-id").Return(routes).Times(1)
			},
			want: routes,
		},
		{
			name: "should return remote error",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetRoutes",
					&apigatewayv2.GetRoutesInput{ApiId: aws.String("an-id")}).Return(nil, remoteError).Once()

				store.On("Get", "apigatewayv2ListAllApiRoutes_api_an-id").Return(nil).Times(1)
			},
			wantErr: remoteError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &cache.MockCache{}
			client := &awstest.MockFakeApiGatewayV2{}
			tt.mocks(client, store)
			r := &apigatewayv2Repository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllApiRoutes(aws.String("an-id"))
			assert.Equal(t, tt.wantErr, err)

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %s -> %s", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
			store.AssertExpectations(t)
			client.AssertExpectations(t)
		})
	}
}
