package provider_test

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakRealm_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakRealmDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRealm_basic(realmName),
				Check:  testAccCheckKeycloakRealmExists("keycloak_realm.realm"),
			},
			{
				Config: testKeycloakRealm_notEnabled(realmName),
				Check:  testAccCheckKeycloakRealmEnabled("keycloak_realm.realm", false),
			},
		},
	})
}

func testAccCheckKeycloakRealmExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realmName := rs.Primary.Attributes["realm"]

		_, err := keycloakClient.GetRealm(realmName)
		if err != nil {
			return fmt.Errorf("error getting realm %s: %s", realmName, err)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmEnabled(resourceName string, enabled bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realmName := rs.Primary.Attributes["realm"]

		realm, err := keycloakClient.GetRealm(realmName)
		if err != nil {
			return fmt.Errorf("error getting realm %s: %s", realmName, err)
		}

		if realm.Enabled != enabled {
			return fmt.Errorf("expected realm %s to have enabled set to %t, but was %t", realm.Realm, enabled, realm.Enabled)
		}

		return nil
	}
}

func testAccCheckKeycloakRealmDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_realm" {
				continue
			}

			realmName := rs.Primary.ID
			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			realm, _ := keycloakClient.GetRealm(realmName)
			if realm != nil {
				return fmt.Errorf("realm %s still exists", realmName)
			}
		}

		return nil
	}
}

func testKeycloakRealm_basic(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}
	`, realm)
}

func testKeycloakRealm_notEnabled(realm string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = false
}
	`, realm)
}
