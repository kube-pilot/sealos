/*
Copyright 2022 labring.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package aws

import (
	"context"
	"fmt"
	"os"

	"github.com/labring/sealos/pkg/utils/logger"

	"github.com/labring/sealos/controllers/infra/common"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/labring/sealos/pkg/types/v1beta1"

	v1 "github.com/labring/sealos/controllers/infra/api/v1"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2DescribeInstancesAPI defines the interface for the DescribeInstances function.
// We use this interface to test the function using a mocked service.
type EC2DescribeInstancesAPI interface {
	DescribeInstances(ctx context.Context,
		params *ec2.DescribeInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// GetInstances retrieves information about your Amazon Elastic Compute Cloud (Amazon EC2) instances.
// Inputs:
//     c is the context of the method call, which includes the AWS Region.
//     api is the interface that defines the method call.
//     input defines the input arguments to the service call.
// Output:
//     If success, a DescribeInstancesOutput object containing the result of the service call and nil.
//     Otherwise, nil and an error from the call to DescribeInstances.
func GetInstances(c context.Context, api EC2DescribeInstancesAPI, input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	return api.DescribeInstances(c, input)
}

func ReconcileInstance(infra *v1.Infra) (*v1beta1.Cluster, error) {
	fmt.Println("access key id is: ", os.Getenv("AWS_DEFAULT_REGION"), os.Getenv("AWS_ACCESS_KEY_ID"))
	if len(infra.Spec.Hosts) == 0 {
		logger.Debug("desired host len is 0")
		return nil, nil
	}

	tag := infra.GetInstancesTag()
	instances, err := GetInstancesByLabel(common.InfraInstancesLabel, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to query instances: %v", err)
	}

	// TODO create instances
	// if len(instances) == 0 {
	// }

	for _, i := range instances {
		if err := checkInstanceTags(i.Tags, common.InfraInstancesLabel, tag); err != nil {
			return nil, fmt.Errorf("instance not contains infra label: %s, %v, %v", *i.InstanceId, i.Tags, err)
		}
	}

	return nil, nil
}

func checkInstanceTags(tags []types.Tag, key, value string) error {
	// TODO check tags contain key, value, if not return an error

	return nil
}

/*
Use this to uniquely identify which cluster the virtual machine belongs to

  key=infra.sealos.io/instances/label
  value=[namespace]/[infra name]

For example:
  apiVersion: infra.sealos.io/v1
  kind: Infra
  metadata:
    name: aws-infra-demo
    namespace: default
The value should be: default/aws-infra-demo
*/
func GetInstancesByLabel(key string, value string) ([]types.Instance, error) {
	var instances []types.Instance
	fmt.Println("access key id is: ", os.Getenv("AWS_DEFAULT_REGION"), os.Getenv("AWS_ACCESS_KEY_ID"))

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("load default config failed %s", err)
	}
	client := ec2.NewFromConfig(cfg)
	tag := fmt.Sprintf("tag:%s", key)
	input := &ec2.DescribeInstancesInput{
		Filters: []types.Filter{
			{
				Name:   &tag,
				Values: []string{value},
			},
		},
	}

	result, err := GetInstances(context.TODO(), client, input)
	if err != nil {
		return nil, fmt.Errorf("got an error retrieving information about your Amazon EC2 instances: %v", err)
	}

	for _, r := range result.Reservations {
		instances = append(instances, r.Instances...)
	}

	return instances, nil
}
