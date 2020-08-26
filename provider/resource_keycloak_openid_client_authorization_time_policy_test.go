package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakOpenidClientAuthorizationTimePolicy(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	clientId := "terraform-" + acctest.RandString(10)
	policyName := "terraform-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testResourceKeycloakOpenidClientAuthorizationTimePolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceKeycloakOpenidClientAuthorizationTimePolicy_basic(realmName, policyName, clientId),
				Check:  testResourceKeycloakOpenidClientAuthorizationTimePolicyExists("keycloak_openid_client_time_policy.test"),
			},
		},
	})
}

func getResourceKeycloakOpenidClientAuthorizationTimePolicyFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientAuthorizationTimePolicy, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	policyId := rs.Primary.ID

	policy, err := keycloakClient.GetOpenidClientAuthorizationTimePolicy(realm, resourceServerId, policyId)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client auth role policy config with alias %s: %s", resourceServerId, err)
	}

	return policy, nil
}

func testResourceKeycloakOpenidClientAuthorizationTimePolicyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_time_policy" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			policyId := rs.Primary.ID

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			policy, _ := keycloakClient.GetOpenidClientAuthorizationTimePolicy(realm, resourceServerId, policyId)
			if policy != nil {
				return fmt.Errorf("policy config with id %s still exists", policyId)
			}
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationTimePolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getResourceKeycloakOpenidClientAuthorizationTimePolicyFromState(s, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationTimePolicy_basic(realm, policyName, clientId string) string {

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

	resource keycloak_openid_client_time_policy test {
		resource_server_id = "${keycloak_openid_client.test.resource_server_id}"
		realm_id = "${keycloak_realm.test.id}"
		name = "%s"
		not_on_or_after = "2500-12-12 01:01:11"
		not_before = "2400-12-12 01:01:11"
		day_month = "1"
		day_month_end = "2"
		year = "2500"
		year_end = "2501"
		month = "1"
		month_end = "5"
		hour = "1"
		hour_end = "5"
		minute = "10"
		minute_end = "30"
		logic = "POSITIVE"
		decision_strategy = "UNANIMOUS"
	}
	`, realm, clientId, policyName)
}
