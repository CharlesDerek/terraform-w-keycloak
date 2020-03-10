package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"testing"
)

func TestGenericRoleMapper_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	parentClientName := "client1-" + acctest.RandString(10)
	parentRoleName := "role-" + acctest.RandString(10)
	childClientName := "client2-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
			},
		},
	})
}

func TestGenericRoleMapper_createAfterManualDestroy(t *testing.T) {
	var role = &keycloak.Role{}
	var childClient = &keycloak.GenericClient{}

	realmName := "terraform-" + acctest.RandString(10)
	parentClientName := "client1-" + acctest.RandString(10)
	parentRoleName := "role-" + acctest.RandString(10)
	childClientName := "client2-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
					testAccCheckKeycloakRoleFetch("keycloak_role.parent-role", role),
					testAccCheckKeycloakGenericClientFetch("keycloak_openid_client.child-client", childClient),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteRoleScopeMapping(childClient.RealmId, childClient.Id, role)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName),
				Check:  testAccCheckKeycloakScopeMappingExists("keycloak_generic_client_role_mapper.child-client-with-parent-client-role"),
			},
		},
	})
}

func testKeycloakGenericRoleMapping_basic(realmName, parentClientName, parentRoleName, childClientName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "parent-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "PUBLIC"
}

resource "keycloak_role" "parent-role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "${keycloak_openid_client.parent-client.id}"
  name      = "%s"
}

resource "keycloak_openid_client" "child-client" {
	realm_id    = "${keycloak_realm.realm.id}"
	client_id   = "%s"
	access_type = "PUBLIC"
}

resource "keycloak_generic_client_role_mapper" "child-client-with-parent-client-role" {
  realm_id  = "${keycloak_realm.realm.id}"
  client_id = "${keycloak_openid_client.child-client.id}"
  role_id   = "${keycloak_role.parent-role.id}"
}
	`, realmName, parentClientName, parentRoleName, childClientName)
}

func testAccCheckKeycloakScopeMappingExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		return nil
	}
}

func testAccCheckKeycloakGenericClientFetch(resourceName string, client *keycloak.GenericClient) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedClient, err := getGenericClientFromState(s, resourceName)
		if err != nil {
			return err
		}

		client.Id = fetchedClient.Id
		client.ClientId = fetchedClient.ClientId
		client.RealmId = fetchedClient.RealmId

		return nil
	}
}

func getGenericClientFromState(s *terraform.State, resourceName string) (*keycloak.GenericClient, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	client, err := keycloakClient.GetGenericClient(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting generic client %s: %s", id, err)
	}

	return client, nil
}
