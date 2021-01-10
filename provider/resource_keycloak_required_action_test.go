package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"regexp"
	"testing"
)

func TestAccKeycloakRequiredAction_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "CONFIGURE_TOTP"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRequiredAction_basic(realmName, requiredActionAlias, 37),
				Check:  testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realmName, requiredActionAlias, 37),
			},
		},
	})
}

func TestAccKeycloakRequiredAction_unregisteredAction(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "webauthn-register"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRequiredAction_basic(realmName, requiredActionAlias, 37),
				Check:  testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realmName, requiredActionAlias, 37),
			},
		},
	})
}

func TestAccKeycloakRequiredAction_invalidAlias(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	randomReqActionAlias := "randomRequiredAction-" + acctest.RandString(10)

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRequiredAction_basic(realmName, randomReqActionAlias, 37),
				ExpectError: regexp.MustCompile("validation error: required action .+ does not exist on the server, installed providers: .+"),
			},
		},
	})
}

func TestAccKeycloakRequiredAction_import(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "terms_and_conditions"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRequiredAction_import(realmName, requiredActionAlias),
				Check:  testAccCheckKeycloakRequiresActionExists(realmName, requiredActionAlias),
			},
			{
				ResourceName:      "keycloak_required_action.required_action",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     realmName + "/" + requiredActionAlias,
			},
		},
	})
}

func TestAccKeycloakRequiredAction_disabledDefault(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "terms_and_conditions"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config:      testKeycloakRequiredAction_disabledDefault(realmName, requiredActionAlias),
				ExpectError: regexp.MustCompile("validation error: a 'default' required action should be enabled, set 'defaultAction' to 'false' or set 'enabled' to 'true'"),
			},
		},
	})
}
func TestAccKeycloakRequiredAction_computedPriority(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	requiredActionAlias := "terms_and_conditions"

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testKeycloakRequiredAction_computedPriority(realmName, requiredActionAlias, 37, 14),
				Check:  testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realmName, requiredActionAlias, 51),
			},
		},
	})
}

func testKeycloakRequiredAction_basic(realm, requiredActionAlias string, priority int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_required_action" "required_action" {
	realm_id		= "${keycloak_realm.realm.realm}"
	alias			= "%s"
	default_action 	= true
	enabled			= true
	name			= "My required Action"
	priority		= %d
}
	`, realm, requiredActionAlias, priority)
}

func testKeycloakRequiredAction_import(realm, requiredActionAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_required_action" "required_action" {
	realm_id		= "${keycloak_realm.realm.realm}"
	alias			= "%s"
}
	`, realm, requiredActionAlias)
}

func testKeycloakRequiredAction_disabledDefault(realm, requiredActionAlias string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_required_action" "required_action" {
	realm_id		= "${keycloak_realm.realm.realm}"
	alias			= "%s"
	default_action 	= true
	enabled			= false
	name			= "My required Action"
	priority		= 56
}
	`, realm, requiredActionAlias)
}

func testKeycloakRequiredAction_computedPriority(realm, requiredActionAlias string, priority1, priorityPlus int) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_required_action" "required_action" {
	realm_id		= "${keycloak_realm.realm.realm}"
	alias			= "VERIFY_EMAIL"
	name			= "My required Action"
	priority		= %d
}

resource "keycloak_required_action" "required_action2" {
	realm_id		= "${keycloak_realm.realm.realm}"
	alias			= "%s"
	name			= "My required Action 2"
	priority		= "${keycloak_required_action.required_action.priority+%d}"
}
	`, realm, priority1, requiredActionAlias, priorityPlus)
}

func testAccCheckKeycloakRequiresActionExistsWithCorrectPriority(realm, requiredActionAlias string, priority int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		action, err := keycloakClient.GetRequiredAction(realm, requiredActionAlias)
		if err != nil {
			return fmt.Errorf("required action not found: %s", requiredActionAlias)
		}

		if action.Priority != priority {
			return fmt.Errorf("expected required action to have priority %d, but got %d", priority, action.Priority)
		}

		return nil
	}
}

func testAccCheckKeycloakRequiresActionExists(realm, requiredActionAlias string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := keycloakClient.GetRequiredAction(realm, requiredActionAlias)
		if err != nil {
			return fmt.Errorf("required action not found: %s", requiredActionAlias)
		}

		return nil
	}
}
