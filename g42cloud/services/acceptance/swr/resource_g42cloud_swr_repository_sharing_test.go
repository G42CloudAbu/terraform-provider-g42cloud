package swr

import (
	"fmt"
	"testing"

	"github.com/g42cloud-terraform/terraform-provider-g42cloud/g42cloud/services/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"

	"github.com/chnsz/golangsdk/openstack/swr/v2/domains"
)

func getResourceRepositorySharing(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	swrClient, err := conf.SwrV2Client(acceptance.G42_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("Error creating G42Cloud SWR client: %s", err)
	}

	return domains.Get(swrClient, state.Primary.Attributes["organization"],
		state.Primary.Attributes["repository"], state.Primary.ID).Extract()
}

func TestAccSWRRepositorySharing_basic(t *testing.T) {
	var domain domains.AccessDomain
	rName := acceptance.RandomAccResourceName()
	resourceName := "g42cloud_swr_repository_sharing.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&domain,
		getResourceRepositorySharing,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckSWRDomian(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccSWRRepositorySharing_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${g42cloud_swr_organization.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "repository",
						"${g42cloud_swr_repository.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "sharing_account", acceptance.G42_SWR_SHARING_ACCOUNT),
					resource.TestCheckResourceAttr(resourceName, "deadline", "forever"),
					resource.TestCheckResourceAttr(resourceName, "permission", "pull"),
				),
			},
			{
				Config: testAccSWRRepositorySharing_update(rName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "organization",
						"${g42cloud_swr_organization.test.name}"),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "repository",
						"${g42cloud_swr_repository.test.name}"),
					resource.TestCheckResourceAttr(resourceName, "sharing_account", acceptance.G42_SWR_SHARING_ACCOUNT),
					resource.TestCheckResourceAttr(resourceName, "deadline", "2099-12-31"),
					resource.TestCheckResourceAttr(resourceName, "permission", "pull"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSWRRepositorySharingImportStateIdFunc(),
			},
		},
	})
}

func testAccSWRRepositorySharingImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		var organization string
		var repositoryID string
		var sharingAccount string
		for _, rs := range s.RootModule().Resources {
			if rs.Type == "g42cloud_swr_organization" {
				organization = rs.Primary.Attributes["name"]
			} else if rs.Type == "g42cloud_swr_repository" {
				repositoryID = rs.Primary.ID
			} else if rs.Type == "g42cloud_swr_repository_sharing" {
				sharingAccount = rs.Primary.ID
			}
		}
		if organization == "" || repositoryID == "" || sharingAccount == "" {
			return "", fmt.Errorf("resource not found: %s/%s/%s", organization, repositoryID, sharingAccount)
		}
		return fmt.Sprintf("%s/%s/%s", organization, repositoryID, sharingAccount), nil
	}
}

func testAccSWRRepositorySharing_basic(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_swr_repository_sharing" "test" {
  organization    = g42cloud_swr_organization.test.name
  repository      = g42cloud_swr_repository.test.name
  sharing_account = "%s"
  permission      = "pull"
  deadline        = "forever"
}
`, testAccSWRRepository_basic(rName), acceptance.G42_SWR_SHARING_ACCOUNT)
}

func testAccSWRRepositorySharing_update(rName string) string {
	return fmt.Sprintf(`
%s

resource "g42cloud_swr_repository_sharing" "test" {
  organization    = g42cloud_swr_organization.test.name
  repository      = g42cloud_swr_repository.test.name
  sharing_account = "%s"
  permission      = "pull"
  deadline        = "2099-12-31"
}
`, testAccSWRRepository_basic(rName), acceptance.G42_SWR_SHARING_ACCOUNT)
}
