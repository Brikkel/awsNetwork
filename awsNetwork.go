package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type AwsNetwork struct {
	pulumi.ResourceState
	SubnetId pulumi.IDOutput
}

func NewAwsNetwork(ctx *pulumi.Context, opts ...pulumi.ResourceOption) (*AwsNetwork, error) {
	awsNetwork := &AwsNetwork{}

	err := ctx.RegisterComponentResource("sue:aws:network", "awsNetworkName", awsNetwork, opts...)
	if err != nil {
		return nil, err
	}

	basicVpc, err := ec2.NewVpc(ctx, "basicVpc", &ec2.VpcArgs{
		CidrBlock: pulumi.String("172.16.0.0/16"),
	}, pulumi.Parent(awsNetwork))

	if err != nil {
		return nil, err
	}

	basicSubnet, err := ec2.NewSubnet(ctx, "basicSubnet", &ec2.SubnetArgs{
		VpcId:            basicVpc.ID(),
		CidrBlock:        pulumi.String("172.16.10.0/24"),
		AvailabilityZone: pulumi.String("eu-west-1a"),
	}, pulumi.Parent(awsNetwork))

	if err != nil {
		return nil, err
	}

	basicIgw, err := ec2.NewInternetGateway(ctx, "basicIgw", &ec2.InternetGatewayArgs{
		VpcId: basicVpc.ID(),
	}, pulumi.Parent(awsNetwork))

	if err != nil {
		return nil, err
	}

	basicRt, err := ec2.NewRouteTable(ctx, "basicRt", &ec2.RouteTableArgs{
		VpcId: basicVpc.ID(),
		Routes: ec2.RouteTableRouteArray{
			&ec2.RouteTableRouteArgs{
				CidrBlock: pulumi.String("0.0.0.0/0"),
				GatewayId: basicIgw.ID(),
			},
		},
	}, pulumi.Parent(awsNetwork))

	if err != nil {
		return nil, err
	}

	_, err = ec2.NewRouteTableAssociation(ctx, "basicRta", &ec2.RouteTableAssociationArgs{
		SubnetId:     basicSubnet.ID(),
		RouteTableId: basicRt.ID(),
	}, pulumi.Parent(awsNetwork))

	if err != nil {
		return nil, err
	}

	awsNetwork.SubnetId = basicSubnet.ID()

	return awsNetwork, nil
}
