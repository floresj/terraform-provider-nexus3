package nexus

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestComponentRaw_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest: true,
		Providers:  testProviders,
		Steps: []resource.TestStep{
			{
				Config: `
				locals {
					repository = "automation"
				}

				resource "nexus_component_raw" "files" {
					repository = local.repository
					src = "https://dl.google.com/go/go1.13.linux-amd64.tar.gz"
					dest = "/go/1.13"
				}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("nexus_component_raw.files", "repository", "automation"),
				),
			},
		},
	})
}
