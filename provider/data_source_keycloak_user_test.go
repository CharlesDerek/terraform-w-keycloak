package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccKeycloakDataSourceUser(t *testing.T) {
	t.Parallel()
	username := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakUser(username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists("keycloak_user.user"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "id", "data.keycloak_user.user", "id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "realm_id", "data.keycloak_user.user", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_user.user", "username", "data.keycloak_user.user", "username"),
					testAccCheckDataKeycloakUser("data.keycloak_user.user"),
				),
			},
		},
	})
}

func testAccCheckDataKeycloakUser(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		realmID := rs.Primary.Attributes["realm_id"]
		username := rs.Primary.Attributes["username"]

		user, err := keycloakClient.GetUser(realmID, id)
		if err != nil {
			return err
		}

		if user.Username != username {
			return fmt.Errorf("expected user with ID %s to have username %s, but got %s", id, username, user.Username)
		}

		return nil
	}
}

func testDataSourceKeycloakUser(username string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	username    = "%s"
	realm_id 	= data.keycloak_realm.realm.id
	enabled    	= true

    email      	= "bob@domain.com"
    first_name 	= "Bob"
    last_name  	= "Bobson"
}

data "keycloak_user" "user" {
	realm_id 	= data.keycloak_realm.realm.id
	username    = keycloak_user.user.username

	depends_on = [
		keycloak_user.user
	]
}
	`, testAccRealm.Realm, username)
}
