// (C) Copyright 2022 Hewlett Packard Enterprise Development LP

package client

import (
	"context"
	"encoding/json"
	"fmt"

	consts "github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/common"
	"github.com/HPE/terraform-provider-hpe/greenlake/cloud/sdk/vmaascmp/models"
)

type LoadBalancerAPIService struct {
	Client APIClientHandler
	Cfg    Configuration
}

func (lb *LoadBalancerAPIService) CreateLoadBalancer(
	ctx context.Context,
	request models.CreateLoadBalancerRequest,
) (models.CreateNetworkLoadBalancerResp, error) {
	LoadBalancerResp := models.CreateNetworkLoadBalancerResp{}
	loadBalancerAPI := &api{
		method: "POST",
		path:   consts.LoadBalancerPath,
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, request, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) UpdateLoadBalancer(
	ctx context.Context,
	lbID int,
	request models.CreateLoadBalancerRequest,
) (models.CreateNetworkLoadBalancerResp, error) {
	LoadBalancerResp := models.CreateNetworkLoadBalancerResp{}
	loadBalancerAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d",
			consts.LoadBalancerPath, lbID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, request, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) DeleteLoadBalancer(
	ctx context.Context,
	lbID int,
) (models.SuccessOrErrorMessage, error) {
	LoadBalancerResp := models.SuccessOrErrorMessage{}
	loadBalancerAPI := &api{
		method: "DELETE",
		path:   fmt.Sprintf("%s/%d", consts.LoadBalancerPath, lbID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, nil, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) GetLoadBalancers(
	ctx context.Context,
) (models.GetNetworkLoadBalancers, error) {
	LoadBalancerResp := models.GetNetworkLoadBalancers{}
	loadBalancerAPI := &api{
		method: "GET",
		path:   consts.LoadBalancerPath,
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, nil, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) GetLoadBalancerTypes(
	ctx context.Context,
	queryParams map[string]string,
) (models.GetLoadBalancerTypes, error) {
	LoadBalancerResp := models.GetLoadBalancerTypes{}
	loadBalancerAPI := &api{
		method: "GET",
		path:   consts.LoadBalancerTypePath,
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, nil, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) GetSpecificLoadBalancers(
	ctx context.Context,
	lbID int,
) (models.GetSpecificNetworkLoadBalancer, error) {
	LoadBalancerResp := models.GetSpecificNetworkLoadBalancer{}
	loadBalancerAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d",
			consts.LoadBalancerPath, lbID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LoadBalancerResp)
		},
	}
	err := loadBalancerAPI.do(ctx, nil, nil)

	return LoadBalancerResp, err
}

func (lb *LoadBalancerAPIService) CreateLBMonitor(
	ctx context.Context,
	request models.CreateLBMonitor,
	lbID int,
) (models.CreateLBMonitorResp, error) {
	LBMonitorResp := models.CreateLBMonitorResp{}
	LBMonitorAPI := &api{
		method: "POST",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerMonitorPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMonitorResp)
		},
	}
	err := LBMonitorAPI.do(ctx, request, nil)

	return LBMonitorResp, err
}

func (lb *LoadBalancerAPIService) UpdateLBMonitor(
	ctx context.Context,
	request models.CreateLBMonitor,
	lbID int,
	lbmonitorID int,
) (models.CreateLBMonitorResp, error) {
	LBMonitorResp := models.CreateLBMonitorResp{}
	LBMonitorAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerMonitorPath, lbmonitorID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMonitorResp)
		},
	}
	err := LBMonitorAPI.do(ctx, request, nil)

	return LBMonitorResp, err
}

func (lb *LoadBalancerAPIService) DeleteLBMonitor(
	ctx context.Context,
	lbID int,
	lbMonitorID int,
) (models.SuccessOrErrorMessage, error) {
	LBMonitorResp := models.SuccessOrErrorMessage{}
	LBMonitorAPI := &api{
		method: "DELETE",
		path: fmt.Sprintf("%s/%d/%s/%d", consts.LoadBalancerPath, lbID,
			consts.LoadBalancerMonitorPath, lbMonitorID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMonitorResp)
		},
	}
	err := LBMonitorAPI.do(ctx, nil, nil)

	return LBMonitorResp, err
}

func (lb *LoadBalancerAPIService) GetLBMonitors(
	ctx context.Context,
	lbID int,
) (models.GetLBMonitors, error) {
	LBMonitorResp := models.GetLBMonitors{}
	LBMonitorAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerMonitorPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMonitorResp)
		},
	}
	err := LBMonitorAPI.do(ctx, nil, nil)

	return LBMonitorResp, err
}

func (lb *LoadBalancerAPIService) GetSpecificLBMonitor(
	ctx context.Context,
	lbID int,
	lbmonitorID int,
) (models.GetSpecificLBMonitor, error) {
	LBMonitorResp := models.GetSpecificLBMonitor{}
	LBMonitorAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerMonitorPath, lbmonitorID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMonitorResp)
		},
	}
	err := LBMonitorAPI.do(ctx, nil, nil)

	return LBMonitorResp, err
}

func (lb *LoadBalancerAPIService) CreateLBProfile(
	ctx context.Context,
	request models.CreateLBProfile,
	lbID int,
) (models.CreateLBProfileResp, error) {
	LBProfileResp := models.CreateLBProfileResp{}
	LBProfileAPI := &api{
		method: "POST",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerProfilePath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBProfileResp)
		},
	}
	err := LBProfileAPI.do(ctx, request, nil)

	return LBProfileResp, err
}

func (lb *LoadBalancerAPIService) UpdateLBProfile(
	ctx context.Context,
	request models.CreateLBProfile,
	lbID int,
	lbProfileID int,
) (models.CreateLBProfileResp, error) {
	LBProfileResp := models.CreateLBProfileResp{}
	LBProfileAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerProfilePath, lbProfileID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBProfileResp)
		},
	}
	err := LBProfileAPI.do(ctx, request, nil)

	return LBProfileResp, err
}

func (lb *LoadBalancerAPIService) DeleteLBProfile(
	ctx context.Context,
	lbID int,
	lbProfileID int,
) (models.SuccessOrErrorMessage, error) {
	LBProfileResp := models.SuccessOrErrorMessage{}
	LBProfileAPI := &api{
		method: "DELETE",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID,
			consts.LoadBalancerProfilePath, lbProfileID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBProfileResp)
		},
	}
	err := LBProfileAPI.do(ctx, nil, nil)

	return LBProfileResp, err
}

func (lb *LoadBalancerAPIService) GetLBProfiles(
	ctx context.Context,
	lbID int,
) (models.GetLBProfile, error) {
	LBProfileResp := models.GetLBProfile{}
	LBProfileAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerProfilePath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBProfileResp)
		},
	}
	err := LBProfileAPI.do(ctx, nil, nil)

	return LBProfileResp, err
}

func (lb *LoadBalancerAPIService) GetSpecificLBProfile(
	ctx context.Context,
	lbID int,
	lbProfileID int,
) (models.GetLBSpecificProfile, error) {
	LBProfileResp := models.GetLBSpecificProfile{}
	LBProfileAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerProfilePath, lbProfileID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBProfileResp)
		},
	}
	err := LBProfileAPI.do(ctx, nil, nil)

	return LBProfileResp, err
}

func (lb *LoadBalancerAPIService) CreateLBPool(
	ctx context.Context,
	request models.CreateLBPool,
	lbID int,
) (models.CreateLBPoolResp, error) {
	LBPoolResp := models.CreateLBPoolResp{}
	LBPoolAPI := &api{
		method: "POST",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerPoolPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBPoolResp)
		},
	}
	err := LBPoolAPI.do(ctx, request, nil)

	return LBPoolResp, err
}

func (lb *LoadBalancerAPIService) UpdateLBPool(
	ctx context.Context,
	request models.CreateLBPool,
	lbID int,
	lbPoolID int,
) (models.CreateLBPoolResp, error) {
	LBPoolResp := models.CreateLBPoolResp{}
	LBPoolAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerPoolPath, lbPoolID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBPoolResp)
		},
	}
	err := LBPoolAPI.do(ctx, request, nil)

	return LBPoolResp, err
}

func (lb *LoadBalancerAPIService) DeleteLBPool(
	ctx context.Context,
	lbID int,
	lbPoolID int,
) (models.SuccessOrErrorMessage, error) {
	LBPoolResp := models.SuccessOrErrorMessage{}
	LBPoolAPI := &api{
		method: "DELETE",
		path: fmt.Sprintf("%s/%d/%s/%d", consts.LoadBalancerPath, lbID,
			consts.LoadBalancerPoolPath, lbPoolID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBPoolResp)
		},
	}
	err := LBPoolAPI.do(ctx, nil, nil)

	return LBPoolResp, err
}

func (lb *LoadBalancerAPIService) GetLBPools(
	ctx context.Context,
	lbID int,
) (models.GetLBPools, error) {
	LBPoolResp := models.GetLBPools{}
	LBPoolAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerPoolPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBPoolResp)
		},
	}
	err := LBPoolAPI.do(ctx, nil, nil)

	return LBPoolResp, err
}

func (lb *LoadBalancerAPIService) GetSpecificLBPool(
	ctx context.Context,
	lbID int,
	lbPoolID int,
) (models.GetSpecificLBPool, error) {
	LBPoolResp := models.GetSpecificLBPool{}
	LBPoolAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerPoolPath, lbPoolID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBPoolResp)
		},
	}
	err := LBPoolAPI.do(ctx, nil, nil)

	return LBPoolResp, err
}

func (lb *LoadBalancerAPIService) CreateLBVirtualServers(
	ctx context.Context,
	request models.CreateLBVirtualServers,
	lbID int,
) (models.LBVirtualServersResp, error) {
	LBVSResp := models.LBVirtualServersResp{}
	LBVSAPI := &api{
		method: "POST",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerVirtualServersPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBVSResp)
		},
	}
	err := LBVSAPI.do(ctx, request, nil)

	return LBVSResp, err
}

func (lb *LoadBalancerAPIService) UpdateLBVirtualServers(
	ctx context.Context,
	request models.CreateLBVirtualServers,
	lbID int,
	lbVirtualServerID int,
) (models.LBVirtualServersResp, error) {
	LBVSResp := models.LBVirtualServersResp{}
	LBVSAPI := &api{
		method: "PUT",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerVirtualServersPath, lbVirtualServerID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBVSResp)
		},
	}
	err := LBVSAPI.do(ctx, request, nil)

	return LBVSResp, err
}

func (lb *LoadBalancerAPIService) DeleteLBVirtualServers(
	ctx context.Context,
	lbID int,
	lbVirtualServerID int,
) (models.SuccessOrErrorMessage, error) {
	LBVSResp := models.SuccessOrErrorMessage{}
	LBVSAPI := &api{
		method: "DELETE",
		path: fmt.Sprintf("%s/%d/%s/%d", consts.LoadBalancerPath, lbID,
			consts.LoadBalancerVirtualServersPath, lbVirtualServerID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBVSResp)
		},
	}
	err := LBVSAPI.do(ctx, nil, nil)

	return LBVSResp, err
}

func (lb *LoadBalancerAPIService) GetLBVirtualServers(
	ctx context.Context,
	lbID int,
) (models.GetLBVirtualServers, error) {
	LBVSResp := models.GetLBVirtualServers{}
	LBVSAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerVirtualServersPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBVSResp)
		},
	}
	err := LBVSAPI.do(ctx, nil, nil)

	return LBVSResp, err
}

func (lb *LoadBalancerAPIService) GetSpecificLBVirtualServer(
	ctx context.Context,
	lbID int,
	lbVSID int,
) (models.GetSpecificLBVirtualServers, error) {
	LBVSResp := models.GetSpecificLBVirtualServers{}
	LBVSAPI := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%d/%s/%d",
			consts.LoadBalancerPath, lbID, consts.LoadBalancerVirtualServersPath, lbVSID),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBVSResp)
		},
	}
	err := LBVSAPI.do(ctx, nil, nil)

	return LBVSResp, err
}

func (lb *LoadBalancerAPIService) GetLBPoolMemberGroup(
	ctx context.Context,
	serviceID int,
) (models.GetMemeberGroupForPool, error) {
	LBMemberGroupResp := models.GetMemeberGroupForPool{}
	LBMemberGroupInput := &api{
		method: "GET",
		path: fmt.Sprintf("%s/%s/%d/%s",
			consts.NetworksPath, consts.ServerPath, serviceID, consts.GroupsPath),
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBMemberGroupResp)
		},
	}
	err := LBMemberGroupInput.do(ctx, nil, nil)

	return LBMemberGroupResp, err
}

func (lb *LoadBalancerAPIService) GetLBVirtualServerSSLCerts(
	ctx context.Context,
) (models.GetSSLCertificates, error) {
	LBCertResp := models.GetSSLCertificates{}
	LBCertInput := &api{
		method: "GET",
		path:   consts.LBSSLCertificatesPath,
		client: lb.Client,
		jsonParser: func(body []byte) error {
			return json.Unmarshal(body, &LBCertResp)
		},
	}
	err := LBCertInput.do(ctx, nil, nil)

	return LBCertResp, err
}
