package aws

import (
	"strings"

	"github.com/snyk/driftctl/pkg/remote/aws/repository"

	remoteerror "github.com/snyk/driftctl/pkg/remote/error"

	"github.com/snyk/driftctl/pkg/resource"
	resourceaws "github.com/snyk/driftctl/pkg/resource/aws"
)

type Route53RecordEnumerator struct {
	client  repository.Route53Repository
	factory resource.ResourceFactory
}

func NewRoute53RecordEnumerator(repo repository.Route53Repository, factory resource.ResourceFactory) *Route53RecordEnumerator {
	return &Route53RecordEnumerator{
		repo,
		factory,
	}
}

func (e *Route53RecordEnumerator) SupportedType() resource.ResourceType {
	return resourceaws.AwsRoute53RecordResourceType
}

func (e *Route53RecordEnumerator) Enumerate() ([]*resource.Resource, error) {

	zones, err := e.client.ListAllZones()
	if err != nil {
		return nil, remoteerror.NewResourceListingErrorWithType(err, string(e.SupportedType()), resourceaws.AwsRoute53ZoneResourceType)
	}

	results := make([]*resource.Resource, 0, len(zones))

	for _, hostedZone := range zones {
		records, err := e.listRecordsForZone(strings.TrimPrefix(*hostedZone.Id, "/hostedzone/"))
		if err != nil {
			return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
		}

		results = append(results, records...)
	}

	return results, err
}

func (e *Route53RecordEnumerator) listRecordsForZone(zoneId string) ([]*resource.Resource, error) {

	records, err := e.client.ListRecordsForZone(zoneId)
	if err != nil {
		return nil, err
	}

	results := make([]*resource.Resource, 0, len(records))

	for _, raw := range records {
		rawType := *raw.Type
		rawName := *raw.Name
		rawSetIdentifier := raw.SetIdentifier

		vars := []string{
			zoneId,
			strings.ToLower(strings.TrimSuffix(rawName, ".")),
			rawType,
		}
		if rawSetIdentifier != nil {
			vars = append(vars, *rawSetIdentifier)
		}

		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				strings.Join(vars, "_"),
				map[string]interface{}{
					"type": rawType,
				},
			),
		)
	}

	return results, nil
}
