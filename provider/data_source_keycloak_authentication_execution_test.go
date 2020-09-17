package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakDataSourceAuthenticationExecution_basic(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	parentFlowAlias := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config: testDataSourceKeycloakAuthenticationExecution_basic(realm, parentFlowAlias),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakAuthenticationExecutionExists("keycloak_authentication_execution.execution"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "id", "data.keycloak_authentication_execution.execution", "id"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "realm_id", "data.keycloak_authentication_execution.execution", "realm_id"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "parent_flow_alias", "data.keycloak_authentication_execution.execution", "parent_flow_alias"),
					resource.TestCheckResourceAttrPair("keycloak_authentication_execution.execution", "authenticator", "data.keycloak_authentication_execution.execution", "provider_id"),
					testAccCheckDataKeycloakAuthenticationExecution("data.keycloak_authentication_execution.execution"),
				),
			},
		},
	})
}

func TestAccKeycloakDataSourceAuthenticationExecution_errorNoExecutions(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	parentFlowAlias := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceKeycloakAuthenticationExecution_errorNoExecutions(realm, parentFlowAlias),
				ExpectError: regexp.MustCompile("no authentication executions found for parent flow alias .*"),
			},
		},
	})
}

func TestAccKeycloakDataSourceAuthenticationExecution_errorWrongProviderId(t *testing.T) {
	realm := "terraform-" + acctest.RandString(10)
	parentFlowAlias := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakAuthenticationExecutionConfigDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testDataSourceKeycloakAuthenticationExecution_errorWrongProviderId(realm, parentFlowAlias, acctest.RandString(10)),
				ExpectError: regexp.MustCompile("no authentication execution under parent flow alias .* with provider id .* found"),
			},
		},
	})
}

func testAccCheckDataKeycloakAuthenticationExecution(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}

		keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

		id := rs.Primary.ID
		realmID := rs.Primary.Attributes["realm_id"]
		parentFlowAlias := rs.Primary.Attributes["parent_flow_alias"]
		providerID := rs.Primary.Attributes["provider_id"]

		authenticationExecutionInfo, err := keycloakClient.GetAuthenticationExecutionInfoFromProviderId(realmID, parentFlowAlias, providerID)
		if err != nil {
			return err
		}

		if authenticationExecutionInfo.Id != id {
			return fmt.Errorf("expected authenticationExecutionInfo with ID %s but got %s", id, authenticationExecutionInfo.Id)
		}

		return nil
	}
}

func testDataSourceKeycloakAuthenticationExecution_basic(realm, parentFlowAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias
	authenticator     = "identity-provider-redirector"
	requirement       = "REQUIRED"
}

data "keycloak_authentication_execution" "execution" {
	realm_id 			= keycloak_realm.realm.id
	parent_flow_alias   = keycloak_authentication_flow.flow.alias
	provider_id     	= "identity-provider-redirector"
}
	`, realm, parentFlowAlias)
}

func testDataSourceKeycloakAuthenticationExecution_errorNoExecutions(realm, parentFlowAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = keycloak_realm.realm.id
	alias    = "%s"
}

data "keycloak_authentication_execution" "execution" {
	realm_id 			= keycloak_realm.realm.id
	parent_flow_alias   = keycloak_authentication_flow.flow.alias
	provider_id     	= "foo"
}
	`, realm, parentFlowAlias)
}

func testDataSourceKeycloakAuthenticationExecution_errorWrongProviderId(realm, parentFlowAlias, providerId string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm   = "%s"
	enabled = true
}

resource "keycloak_authentication_flow" "flow" {
	realm_id = keycloak_realm.realm.id
	alias    = "%s"
}

resource "keycloak_authentication_execution" "execution" {
	realm_id          = keycloak_realm.realm.id
	parent_flow_alias = keycloak_authentication_flow.flow.alias
	authenticator     = "identity-provider-redirector"
	requirement       = "REQUIRED"
}

data "keycloak_authentication_execution" "execution" {
	realm_id 			= keycloak_realm.realm.id
	parent_flow_alias   = keycloak_authentication_flow.flow.alias
	provider_id     	= "%s"
}
	`, realm, parentFlowAlias, providerId)
}
