package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccKeycloakDataSourceRealm_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)

	resourceName := "keycloak_realm.my_realm"
	dataSourceName := "data.keycloak_realm.realm"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakRealm_basic(realm),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataSourceName, "realm", resourceName, "realm"),
					resource.TestCheckResourceAttrPair(dataSourceName, "enabled", resourceName, "enabled"),
					resource.TestCheckResourceAttrPair(dataSourceName, "display_name", resourceName, "display_name"),
					resource.TestCheckResourceAttrPair(dataSourceName, "display_name_html", resourceName, "display_name_html"),
				),
			},
		},
	})
}

func testDataSourceKeycloakRealm_basic(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "my_realm" {
	realm             = "%s"
	enabled           = true
	display_name      = "foo"
	display_name_html = "<b>foo</b>"
}

data "keycloak_realm" "realm" {
	realm = keycloak_realm.my_realm.id
}`, realm)
}
