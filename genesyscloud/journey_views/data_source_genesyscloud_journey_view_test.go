package journey_views

import (
	"fmt"
	"terraform-provider-genesyscloud/genesyscloud/provider"
	"terraform-provider-genesyscloud/genesyscloud/util"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceJourneyViewBasic(t *testing.T) {
	var (
		journeyResource = "test-journey"
		journeyName     = "TerraformTestJourney-" + uuid.NewString()
		duration        = "P1Y"
		elementsBlock   = ""

		journeyDataSource = "test-journey-ds"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				// Create
				Config: generateJourneyView(
					journeyResource,
					journeyName,
					duration,
					elementsBlock,
				) + generateJourneyViewDataSource(
					journeyDataSource,
					journeyName,
					"genesyscloud_journey_views."+journeyResource,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.genesyscloud_journey_views."+journeyDataSource,
						"id", "genesyscloud_journey_views."+journeyResource, "id",
					),
				),
			},
		},
	})
}

func TestAccDataSourceJourneyViewCaching(t *testing.T) {
	var (
		journey1ResourceId = "journey1"
		journeyName1       = "terraform test journey " + uuid.NewString()
		journey2ResourceId = "journey2"
		journeyName2       = "terraform test journey " + uuid.NewString()
		journey3ResourceId = "journey3"
		journeyName3       = "terraform test journey " + uuid.NewString()
		duration           = "P1Y"
		elementsBlock      = ""
		dataSource1Id      = "data-1"
		dataSource2Id      = "data-2"
		dataSource3Id      = "data-3"
	)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { util.TestAccPreCheck(t) },
		ProviderFactories: provider.GetProviderFactories(providerResources, providerDataSources),
		Steps: []resource.TestStep{
			{
				PreConfig: func() {
					time.Sleep(45 * time.Second)
				},
				Config: generateJourneyView( // journey resource
					journey1ResourceId,
					journeyName1,
					duration,
					elementsBlock,
				) + generateJourneyView( // journey resource
					journey2ResourceId,
					journeyName2,
					duration,
					elementsBlock,
				) + generateJourneyView( // journey resource
					journey3ResourceId,
					journeyName3,
					duration,
					elementsBlock,
				) + generateJourneyViewDataSource( // journey data source
					dataSource1Id,
					journeyName1,
					"genesyscloud_journey_views."+journey1ResourceId,
				) + generateJourneyViewDataSource( // journey data source
					dataSource2Id,
					journeyName2,
					"genesyscloud_journey_views."+journey2ResourceId,
				) + generateJourneyViewDataSource( // journey data source
					dataSource3Id,
					journeyName3,
					"genesyscloud_journey_views."+journey3ResourceId,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("genesyscloud_journey_views."+journey1ResourceId, "id",
						"data.genesyscloud_journey_views."+dataSource1Id, "id"),
					resource.TestCheckResourceAttrPair("genesyscloud_journey_views."+journey2ResourceId, "id",
						"data.genesyscloud_journey_views."+dataSource2Id, "id"),
					resource.TestCheckResourceAttrPair("genesyscloud_journey_views."+journey3ResourceId, "id",
						"data.genesyscloud_journey_views."+dataSource3Id, "id"),
				),
			},
		},
		CheckDestroy: testVerifyJourneyViewsDestroyed,
	})
}

func generateJourneyViewDataSource(
	resourceID string,
	name string,
	dependsOnResource string) string {
	return fmt.Sprintf(`data "genesyscloud_journey_views" "%s" {
		name = "%s"
		depends_on = [%s]
	}
	`, resourceID, name, dependsOnResource)
}
