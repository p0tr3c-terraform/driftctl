package aws

import "github.com/cloudskiff/driftctl/pkg/resource"

const AwsLaunchTemplateResourceType = "aws_launch_template"

func initAwsLaunchTemplateMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetFlags(AwsLaunchTemplateResourceType, resource.FlagDeepMode)
}
