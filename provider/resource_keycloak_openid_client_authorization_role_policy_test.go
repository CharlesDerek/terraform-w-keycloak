package provider

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakOpenidClientAuthorizationRolePolicy_basic(t *testing.T) {
	t.Parallel()

	clientId := acctest.RandomWithPrefix("tf-acc")
	roleName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testResourceKeycloakOpenidClientAuthorizationRolePolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceKeycloakOpenidClientAuthorizationRolePolicy_basic(roleName, clientId),
				Check:  testResourceKeycloakOpenidClientAuthorizationRolePolicyExists("keycloak_openid_client_role_policy.test"),
			},
		},
	})
}

func TestAccKeycloakOpenidClientAuthorizationRolePolicy_multiple(t *testing.T) {
	t.Parallel()

	clientId := acctest.RandomWithPrefix("tf-acc")
	var roleNames []string
	for i := 0; i < acctest.RandIntRange(7, 12); i++ {
		roleNames = append(roleNames, acctest.RandomWithPrefix("tf-acc"))
	}

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testResourceKeycloakOpenidClientAuthorizationRolePolicyDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testResourceKeycloakOpenidClientAuthorizationRolePolicy_multipleRoles(roleNames, clientId),
				Check:  testResourceKeycloakOpenidClientAuthorizationRolePolicyExists("keycloak_openid_client_role_policy.test"),
			},
		},
	})
}

func getResourceKeycloakOpenidClientAuthorizationRolePolicyFromState(s *terraform.State, resourceName string) (*keycloak.OpenidClientAuthorizationRolePolicy, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	realm := rs.Primary.Attributes["realm_id"]
	resourceServerId := rs.Primary.Attributes["resource_server_id"]
	policyId := rs.Primary.ID

	policy, err := keycloakClient.GetOpenidClientAuthorizationRolePolicy(realm, resourceServerId, policyId)
	if err != nil {
		return nil, fmt.Errorf("error getting openid client auth role policy config with alias %s: %s", resourceServerId, err)
	}

	return policy, nil
}

func testResourceKeycloakOpenidClientAuthorizationRolePolicyDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_openid_client_role_policy" {
				continue
			}

			realm := rs.Primary.Attributes["realm_id"]
			resourceServerId := rs.Primary.Attributes["resource_server_id"]
			policyId := rs.Primary.ID

			policy, _ := keycloakClient.GetOpenidClientAuthorizationRolePolicy(realm, resourceServerId, policyId)
			if policy != nil {
				return fmt.Errorf("policy config with id %s still exists", policyId)
			}
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationRolePolicyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getResourceKeycloakOpenidClientAuthorizationRolePolicyFromState(s, resourceName)

		if err != nil {
			return err
		}

		return nil
	}
}

func testResourceKeycloakOpenidClientAuthorizationRolePolicy_basic(roleName, clientId string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource keycloak_openid_client test {
	client_id                = "%s"
	realm_id                 = data.keycloak_realm.realm.id
	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
	authorization {
		policy_enforcement_mode = "ENFORCING"
	}
}

resource "keycloak_role" "test" {
	realm_id    = data.keycloak_realm.realm.id
	name        = "%s"
}

resource keycloak_openid_client_role_policy test {
	resource_server_id = keycloak_openid_client.test.resource_server_id
	realm_id = data.keycloak_realm.realm.id
	name = "keycloak_openid_client_role_policy"
	decision_strategy = "AFFIRMATIVE"
	logic = "POSITIVE"
	type = "role"
	role  {
		id = keycloak_role.test.id
		required = false
	}
}
	`, testAccRealm.Realm, roleName, clientId)
}

func testResourceKeycloakOpenidClientAuthorizationRolePolicy_multipleRoles(roleNames []string, clientId string) string {
	var (
		roles        strings.Builder
		rolePolicies strings.Builder
	)
	for i, roleName := range roleNames {
		roles.WriteString(fmt.Sprintf(`
resource "keycloak_role" "role_%d" {
	realm_id    = data.keycloak_realm.realm.id
	name        = "%s"
}
`, i, roleName))
		rolePolicies.WriteString(fmt.Sprintf(`
	role  {
		id = keycloak_role.role_%d.id
		required = false
	}
`, i))
	}

	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource keycloak_openid_client test {
	client_id                = "%s"
	realm_id                 = data.keycloak_realm.realm.id
	access_type              = "CONFIDENTIAL"
	service_accounts_enabled = true
	authorization {
		policy_enforcement_mode = "ENFORCING"
	}
}

%s

resource keycloak_openid_client_role_policy test {
	resource_server_id = keycloak_openid_client.test.resource_server_id
	realm_id = data.keycloak_realm.realm.id
	name = "keycloak_openid_client_role_policy"
	decision_strategy = "AFFIRMATIVE"
	logic = "POSITIVE"
	type = "role"

%s

}
	`, testAccRealm.Realm, clientId, roles.String(), rolePolicies.String())
}
