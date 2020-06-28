package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"regexp"
	"testing"
)

func TestAccKeycloakUserRoles_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	realmRoleName := "terraform-role-" + acctest.RandString(10)
	openIdClientName := "terraform-openid-client-" + acctest.RandString(10)
	openIdRoleName := "terraform-role-" + acctest.RandString(10)
	samlClientName := "terraform-saml-client-" + acctest.RandString(10)
	samlRoleName := "terraform-role-" + acctest.RandString(10)
	userName := "terraform-user-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUserRoles_basic(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			{
				ResourceName:      "keycloak_user_roles.user_roles",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// check destroy
			{
				Config: testKeycloakUserRoles_noUserRoles(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName),
				Check:  testAccCheckKeycloakUserHasNoRoles("keycloak_user.user"),
			},
		},
	})
}

func TestAccKeycloakUserRoles_update(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)

	realmRoleOneName := "terraform-role-" + acctest.RandString(10)
	realmRoleTwoName := "terraform-role-" + acctest.RandString(10)
	openIdClientName := "terraform-openid-client-" + acctest.RandString(10)
	openIdRoleOneName := "terraform-role-" + acctest.RandString(10)
	openIdRoleTwoName := "terraform-role-" + acctest.RandString(10)
	samlClientName := "terraform-saml-client-" + acctest.RandString(10)
	samlRoleOneName := "terraform-role-" + acctest.RandString(10)
	samlRoleTwoName := "terraform-role-" + acctest.RandString(10)
	userName := "terraform-user-" + acctest.RandString(10)

	allRoleIds := []string{
		"${keycloak_role.realm_role_one.id}",
		"${keycloak_role.realm_role_two.id}",
		"${keycloak_role.openid_client_role_one.id}",
		"${keycloak_role.openid_client_role_two.id}",
		"${keycloak_role.saml_client_role_one.id}",
		"${keycloak_role.saml_client_role_two.id}",
		"${data.keycloak_role.offline_access.id}",
	}

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			// initial setup, resource is defined but no roles are specified
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, []string{}),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// add all roles
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, allRoleIds),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// remove some
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, []string{
					"${keycloak_role.realm_role_two.id}",
					"${keycloak_role.openid_client_role_one.id}",
					"${keycloak_role.openid_client_role_two.id}",
					"${data.keycloak_role.offline_access.id}",
				}),
				Check: testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// add some and remove some
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, []string{
					"${keycloak_role.saml_client_role_one.id}",
					"${keycloak_role.saml_client_role_two.id}",
					"${keycloak_role.realm_role_one.id}",
				}),
				Check: testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// add some and remove some again
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, []string{
					"${keycloak_role.saml_client_role_one.id}",
					"${keycloak_role.openid_client_role_two.id}",
					"${keycloak_role.realm_role_two.id}",
					"${data.keycloak_role.offline_access.id}",
				}),
				Check: testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// add all back
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, allRoleIds),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// random scenario 1
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, randomStringSliceSubset(allRoleIds)),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// random scenario 2
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, randomStringSliceSubset(allRoleIds)),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// random scenario 3
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, randomStringSliceSubset(allRoleIds)),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
			// remove all
			{
				Config: testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, []string{}),
				Check:  testAccCheckKeycloakUserHasRoles("keycloak_user_roles.user_roles"),
			},
		},
	})
}

func flattenUserRoles(userRoleMapping *keycloak.UserRoleMapping) ([]string, error) {
	var roles []string

	for _, realmRole := range userRoleMapping.RealmMappings {
		roles = append(roles, realmRole.Name)
	}

	for _, clientRoleMapping := range userRoleMapping.ClientMappings {
		for _, clientRole := range clientRoleMapping.Mappings {
			roles = append(roles, fmt.Sprintf("%s/%s", clientRoleMapping.Id, clientRole.Name))
		}
	}

	return roles, nil
}

func testAccCheckKeycloakUserHasRoles(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		userId := rs.Primary.Attributes["user_id"]

		var roles []*keycloak.Role
		for k, v := range rs.Primary.Attributes {
			if match, _ := regexp.MatchString("role_ids\\.[^#]+", k); !match {
				continue
			}

			role, err := keycloakClient.GetRole(realm, v)
			if err != nil {
				return err
			}

			roles = append(roles, role)
		}

		user, err := keycloakClient.GetUser(realm, userId)
		if err != nil {
			return err
		}

		userRoleMappings, err := keycloakClient.GetUserRoleMappings(realm, userId)
		if err != nil {
			return err
		}

		userRoles, err := flattenUserRoles(userRoleMappings)
		if err != nil {
			return err
		}

		if len(userRoles) != len(roles) {
			return fmt.Errorf("expected number of user roles to be %d, got %d", len(roles), len(userRoles))
		}

		for _, role := range roles {
			var expectedRoleString string
			if role.ClientRole {
				expectedRoleString = fmt.Sprintf("%s/%s", role.ClientId, role.Name)
			} else {
				expectedRoleString = role.Name
			}

			found := false

			for _, userRole := range userRoles {
				if userRole == expectedRoleString {
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("expected to find role %s assigned to user %s", expectedRoleString, user.Username)
			}
		}

		return nil
	}
}

func testAccCheckKeycloakUserHasNoRoles(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		realm := rs.Primary.Attributes["realm_id"]
		id := rs.Primary.ID

		user, err := keycloakClient.GetUser(realm, id)
		if err != nil {
			return err
		}

		userRoleMapping, err := keycloakClient.GetUserRoleMappings(realm, id)
		if err != nil {
			return err
		}

		if len(userRoleMapping.RealmMappings) != 0 || len(userRoleMapping.ClientMappings) != 0 {
			return fmt.Errorf("expected user %s to have no roles", user.Username)
		}

		return nil
	}
}

func testKeycloakUserRoles_basic(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_openid_client" "account" {
	realm_id = "${keycloak_realm.realm.id}"
	client_id = "account"
}

data "keycloak_role" "manage_account" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${data.keycloak_openid_client.account.id}"
	name 	  = "manage-account"
}

data "keycloak_role" "view_profile" {
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${data.keycloak_openid_client.account.id}"
	name 	  = "view-profile"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

data "keycloak_role" "uma_authorization" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "uma_authorization"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_user_roles" "user_roles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.user.id}"

	role_ids = [
		"${keycloak_role.realm_role.id}",
		"${keycloak_role.openid_client_role.id}",
		"${keycloak_role.saml_client_role.id}",

		# default roles
		"${data.keycloak_role.offline_access.id}",
		"${data.keycloak_role.uma_authorization.id}",
		"${data.keycloak_role.manage_account.id}",
		"${data.keycloak_role.view_profile.id}",
	]
}
	`, realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName)
}

func testKeycloakUserRoles_noUserRoles(realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}
	`, realmName, openIdClientName, samlClientName, realmRoleName, openIdRoleName, samlRoleName, userName)
}

func testKeycloakUserRoles_update(realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName string, roleIds []string) string {
	tfRoleIds := fmt.Sprintf("role_ids = %s", arrayOfStringsForTerraformResource(roleIds))

	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_openid_client" "openid_client" {
	client_id   = "%s"
	realm_id    = "${keycloak_realm.realm.id}"
	access_type = "CONFIDENTIAL"
}

resource "keycloak_saml_client" "saml_client" {
	client_id = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role_one" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "realm_role_two" {
	name     = "%s"
	realm_id = "${keycloak_realm.realm.id}"
}

resource "keycloak_role" "openid_client_role_one" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "openid_client_role_two" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_openid_client.openid_client.id}"
}

resource "keycloak_role" "saml_client_role_one" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

resource "keycloak_role" "saml_client_role_two" {
	name      = "%s"
	realm_id  = "${keycloak_realm.realm.id}"
	client_id = "${keycloak_saml_client.saml_client.id}"
}

data "keycloak_role" "offline_access" {
	realm_id  = "${keycloak_realm.realm.id}"
	name      = "offline_access"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}

resource "keycloak_user_roles" "user_roles" {
	realm_id = "${keycloak_realm.realm.id}"
	user_id = "${keycloak_user.user.id}"

	%s
}
	`, realmName, openIdClientName, samlClientName, realmRoleOneName, realmRoleTwoName, openIdRoleOneName, openIdRoleTwoName, samlRoleOneName, samlRoleTwoName, userName, tfRoleIds)
}
