package aws

import "github.com/snyk/driftctl/pkg/resource"

const AwsDefaultRouteTableResourceType = "aws_default_route_table"

func initAwsDefaultRouteTableMetadata(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetResolveReadAttributesFunc(AwsDefaultRouteTableResourceType, func(res *resource.Resource) map[string]string {
		return map[string]string{
			"vpc_id": *res.Attributes().GetString("vpc_id"),
		}
	})
	resourceSchemaRepository.SetFlags(AwsDefaultRouteTableResourceType, resource.FlagDeepMode)
}
