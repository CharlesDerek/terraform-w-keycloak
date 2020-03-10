package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestResourceKeycloakOpenidClientAuthorizationClientPolicy(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	roleName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testResourceKeycloakOpenidClientAuthorizationClientPolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceKeycloakOpenidClientAuthorizationClientPolicy_basic(realmName, roleName, clientId),
				Check:  testResourceKeycloakOpenidClientAuthorizationClientPolicyExists("keycloak_openid_client_client_policy.test"),
			},
		},
	})
}

func getResourceKeycloakOpenidClientAuthorizationClientPolicyFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientAuthorizationClientPolicy, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	policyId := rs.Primary.ID

	policy, err := keycloakClient.GetOpenidClientAuthorizationClientPolicy(realm, resourceServerId, policyId)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client auth role policy config with alias %s: %s", resourceServerId, err)
	}

	return policy, nil
}

func testResourceKeycloakOpenidClientAuthorizationClientPolicyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_client_policy" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			policyId := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			policy, _ := keycloakClient.GetOpenidClientAuthorizationClientPolicy(realm, resourceServerId, policyId)
			if policy != nil {
				return fmt.Errorf("policy config with id %s still exists", policyId)
			}
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationClientPolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getResourceKeycloakOpenidClientAuthorizationClientPolicyFromState(s, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationClientPolicy_basic(realm, roleName, clientId string) string {

	return fmt.Sprintf(`
	resource keycloak_realm test {
		realm = "%s"
	}
	
	resource keycloak_openid_client test {
		client_id                = "%s"
		realm_id                 = "${keycloak_realm.test.id}"
		access_type              = "CONFIDENTIAL"
		service_accounts_enabled = true
		authorization {
			policy_enforcement_mode = "ENFORCING"
		}
	}
	
	resource keycloak_openid_client_client_policy test {
		resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
		realm_id = "${keycloak_realm.test.id}"
		name = "keycloak_openid_client_client_policy"
		decision_strategy = "AFFIRMATIVE"
		logic = "POSITIVE"
		clients = ["${keycloak_openid_client.test.resource_server_id}"]
	}
	`, realm, clientId)
}
