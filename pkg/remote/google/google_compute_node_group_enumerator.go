package google

import (
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/remote/google/repository"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/google"
)

type GoogleComputeNodeGroupEnumerator struct {
	repository repository.AssetRepository
	factory    resource.ResourceFactory
}

func NewGoogleComputeNodeGroupEnumerator(repo repository.AssetRepository, factory resource.ResourceFactory) *GoogleComputeNodeGroupEnumerator {
	return &GoogleComputeNodeGroupEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *GoogleComputeNodeGroupEnumerator) SupportedType() resource.ResourceType {
	return google.GoogleComputeNodeGroupResourceType
}

func (e *GoogleComputeNodeGroupEnumerator) Enumerate() ([]*resource.Resource, error) {
	checks, err := e.repository.SearchAllNodeGroups()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(checks))
	for _, res := range checks {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				trimResourceName(res.GetName()),
				map[string]interface{}{
					"name": res.GetName(),
				},
			),
		)
	}

	return results, err
}
