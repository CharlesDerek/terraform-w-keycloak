package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"strings"
	"testing"
)

func TestAccKeycloakDefaultRoles_basic(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakDefaultRolesDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakDefaultRoles_basic(realmName),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
		},
	})
}

func TestAccKeycloakDefaultRoles_updateDefaultRoles(t *testing.T) {
	realmName := acctest.RandomWithPrefix("tf-acc")

	groupDefaultRolesOne := &keycloak.DefaultRoles{
		RealmId:      testAccRealmUserFederation.Realm,
		DefaultRoles: []string{"\"uma_authorization\""},
	}

	groupDefaultRolesTwo := &keycloak.DefaultRoles{
		RealmId:      testAccRealmUserFederation.Realm,
		DefaultRoles: []string{"\"uma_authorization\",", "\"offline_access\""},
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakDefaultRolesDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakDefaultRoles_basicFromInterface(realmName, groupDefaultRolesOne),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
			{
				Config: testKeycloakDefaultRoles_basicFromInterface(realmName, groupDefaultRolesTwo),
				Check:  testAccCheckDefaultRolesExists("keycloak_default_roles.default_roles"),
			},
		},
	})
}

func testAccCheckDefaultRolesExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getKeycloakDefaultRolesFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakDefaultRolesDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_default_roles" || strings.HasPrefix(name, "data") {
				continue
			}

			id := rs.Primary.ID
			realmId := rs.Primary.Attributes["realm_id"]
			realm, _ := keycloakClient.GetRealm(realmId)
			// Since we started using the realm as a resource, only this destroy check will be triggered.
			if realm == nil {
				return nil
			}

			composites, err := keycloakClient.GetDefaultRoles(realmId, id)
			if err != nil {
				return fmt.Errorf("error getting defaultRoles with id %s: %s", id, err)
			}

			defaultRoles, err := getDefaultRoleNames(composites)
			if err != nil {
				return err
			}
			if len(defaultRoles) != 0 {
				return fmt.Errorf("%s with id %s still exists", name, id)
			}
		}
		return nil
	}
}

func getKeycloakDefaultRolesFromState(s *terraform.State, resourceName string) (*keycloak.DefaultRoles, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	composites, err := keycloakClient.GetDefaultRoles(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting defaultRoles with id %s: %s", id, err)
	}

	defaultRoleNamesList, _ := getDefaultRoleNames(composites)

	defaultRoles := &keycloak.DefaultRoles{
		Id:           id,
		RealmId:      realm,
		DefaultRoles: defaultRoleNamesList,
	}

	return defaultRoles, nil
}

func testKeycloakDefaultRoles_basic(realmName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_default_roles" "default_roles" {
	realm_id  = keycloak_realm.realm.id
    default_roles = ["uma_authorization"]
}
	`, realmName)
}

func testKeycloakDefaultRoles_basicFromInterface(realmName string, defaultRoles *keycloak.DefaultRoles) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_default_roles" "default_roles" {
	realm_id  = keycloak_realm.realm.id
    default_roles = %s
}
	`, realmName, defaultRoles.DefaultRoles)
}
