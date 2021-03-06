// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package main

import (
	"testing"

	"github.com/MustWin/baremetal-sdk-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/oracle/terraform-provider-baremetal/client/mocks"
)

func TestLoadBalancerCertificatesDatasource(t *testing.T) {
	client := &mocks.BareMetalClient{}
	providers := map[string]terraform.ResourceProvider{
		"baremetal": Provider(func(d *schema.ResourceData) (interface{}, error) {
			return client, nil
		}),
	}
	resourceName := "data.baremetal_load_balancer_certificates.t"
	config := `
data "baremetal_load_balancer_certificates" "t" {
  load_balancer_id = "ocid1.loadbalancer.stub_id"
}
`
	config += testProviderConfig

	loadbalancerID := "ocid1.loadbalancer.stub_id"
	list := &baremetal.ListCertificates{
		Certificates: []baremetal.Certificate{
			{CertificateName: "stub_name1"},
			{CertificateName: "stub_name2"},
		},
	}
	client.On(
		"ListCertificates",
		loadbalancerID,
		(*baremetal.ClientRequestOptions)(nil),
	).Return(list, nil)

	resource.UnitTest(t, resource.TestCase{
		PreventPostDestroyRefresh: true,
		Providers:                 providers,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "load_balancer_id", loadbalancerID),
					resource.TestCheckResourceAttr(resourceName, "certificates.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "certificates.0.certificate_name", "stub_name1"),
					resource.TestCheckResourceAttr(resourceName, "certificates.1.certificate_name", "stub_name2"),
				),
			},
		},
	})
}
