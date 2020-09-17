package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakAuthenticationExecution_basic(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	parentAuthFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationExecution_basic(realmName, parentAuthFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
			},
			{
				ResourceName:      "keycloak_authentication_execution.execution",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getExecutionImportId("keycloak_authentication_execution.execution"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationExecution_createAfterManualDestroy(t *testing.T) {
	var authenticationExecution = &keycloak.AuthenticationExecution{}

	realmName := "terraform-" + acctest.RandString(10)
	authParentFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationExecution_basic(realmName, authParentFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					testAccCheckKeycloakAuthenticationExecutionFetch("keycloak_authentication_execution.execution", authenticationExecution),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteAuthenticationExecution(authenticationExecution.RealmId, authenticationExecution.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakAuthenticationExecution_basic(realmName, authParentFlowAlias),
				Check:  testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
			},
		},
	})
}

func TestAccKeycloakAuthenticationExecution_updateAuthenticationExecutionRequirement(t *testing.T) {
	realmName := "terraform-r-" + acctest.RandString(10)
	authParentFlowAlias := "terraform-parent-flow-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationSubFlowDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakAuthenticationExecution_basic(realmName, authParentFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution.execution", "requirement", "DISABLED"),
				),
			},
			{
				Config: testKeycloakAuthenticationExecution_basicWithRequirement(realmName, authParentFlowAlias, "REQUIRED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution.execution", "requirement", "REQUIRED"),
				),
			},
			{
				Config: testKeycloakAuthenticationExecution_basicWithRequirement(realmName, authParentFlowAlias, "DISABLED"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					resource.TestCheckResourceAttr("keycloak_authentication_execution.execution", "requirement", "DISABLED"),
				),
			},
		},
	})
}

func testAccCheckKeycloakAuthenticationExecutionExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getAuthenticationExecutionFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakAuthenticationExecutionFetch(resourceName string, authenticationExecution *keycloak.AuthenticationExecution) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedAuthenticationExecution, err := getAuthenticationExecutionFromState(s, resourceName)
		if err != nil {
			return err
		}

		authenticationExecution.Id = fetchedAuthenticationExecution.Id
		authenticationExecution.ParentFlowAlias = fetchedAuthenticationExecution.ParentFlowAlias
		authenticationExecution.RealmId = fetchedAuthenticationExecution.RealmId

		return nil
	}
}

func testAccCheckKeycloakAuthenticationExecutionDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_authentication_execution" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]
			parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			authenticationExecution, _ := keycloakClient.GetAuthenticationExecution(realm, parentFlowAlias, id)
			if authenticationExecution != nil {
				return fmt.Errorf("authentication flow with id %s still exists", id)
			}
		}

		return nil
	}
}

func getAuthenticationExecutionFromState(s *terraform.State, resourceName string) (*keycloak.AuthenticationExecution, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]
	parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]

	authenticationExecution, err := keycloakClient.GetAuthenticationExecution(realm, parentFlowAlias, id)

	if err != nil {
		return nil, fmt.Errorf("error getting authentication execution with id %s: %s", id, err)
	}

	return authenticationExecution, nil
}

func getExecutionImportId(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("resource not found: %s", resourceName)
		}

		id := rs.Primary.ID
		parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]
		realmId := rs.Primary.Attributes["realm_id"]

		return fmt.Sprintf("%s/%s/%s", realmId, parentFlowAlias, id), nil
	}
}

func testKeycloakAuthenticationExecution_basic(realm, parentAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "%s"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator = "auth-cookie"
}
	`, realm, parentAlias)
}

func testKeycloakAuthenticationExecution_basicWithRequirement(realm, parentAlias, requirement string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = "${keycloak_realm.realm.id}"
	alias    = "%s"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id = "${keycloak_realm.realm.id}"
	parent_flow_alias = "${keycloak_authentication_flow.flow.alias}"
	authenticator = "auth-cookie"
	requirement = "%s"
}
	`, realm, parentAlias, requirement)
}
