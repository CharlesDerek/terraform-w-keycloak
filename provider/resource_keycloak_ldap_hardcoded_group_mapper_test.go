package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
)

func TestAccKeycloakLdapHardcodedGroupMapper_basic(t *testing.T) {
	t.Parallel()
	groupName := acctest.RandomWithPrefix("tf-acc")
	groupMapperName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapHardcodedGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapHardcodedGroupMapper(groupName, groupMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedGroupMapperExists("keycloak_ldap_hardcoded_group_mapper.hardcoded_group_mapper"),
			},
			{
				ResourceName:      "keycloak_ldap_hardcoded_group_mapper.hardcoded_group_mapper",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: getLdapGenericMapperImportId("keycloak_ldap_hardcoded_group_mapper.hardcoded_group_mapper"),
			},
		},
	})
}

func TestAccKeycloakLdapHardcodedGroupMapper_createAfterManualDestroy(t *testing.T) {
	t.Parallel()
	var mapper = &keycloak.LdapHardcodedGroupMapper{}

	groupName := acctest.RandomWithPrefix("tf-acc")
	groupMapperName := acctest.RandomWithPrefix("tf-acc")

	resource.Test(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      testAccCheckKeycloakLdapHardcodedGroupMapperDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakLdapHardcodedGroupMapper(groupName, groupMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedGroupMapperFetch("keycloak_ldap_hardcoded_group_mapper.hardcoded_group_mapper", mapper),
			},
			{
				PreConfig: func() {
					err := keycloakClient.DeleteLdapHardcodedGroupMapper(testCtx, mapper.RealmId, mapper.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakLdapHardcodedGroupMapper(groupName, groupMapperName),
				Check:  testAccCheckKeycloakLdapHardcodedGroupMapperExists("keycloak_ldap_hardcoded_group_mapper.hardcoded_group_mapper"),
			},
		},
	})
}

func testAccCheckKeycloakLdapHardcodedGroupMapperExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getLdapHardcodedGroupMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakLdapHardcodedGroupMapperFetch(resourceName string, mapper *keycloak.LdapHardcodedGroupMapper) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedMapper, err := getLdapHardcodedGroupMapperFromState(s, resourceName)
		if err != nil {
			return err
		}

		mapper.Id = fetchedMapper.Id
		mapper.RealmId = fetchedMapper.RealmId

		return nil
	}
}

func testAccCheckKeycloakLdapHardcodedGroupMapperDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_ldap_hardcoded_group_mapper" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			ldapMapper, _ := keycloakClient.GetLdapHardcodedGroupMapper(testCtx, realm, id)
			if ldapMapper != nil {
				return fmt.Errorf("ldap hardcoded group mapper with id %s still exists", id)
			}
		}

		return nil
	}
}

func getLdapHardcodedGroupMapperFromState(s *terraform.State, resourceName string) (*keycloak.LdapHardcodedGroupMapper, error) {
	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	ldapMapper, err := keycloakClient.GetLdapHardcodedGroupMapper(testCtx, realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting ldap group mapper with id %s: %s", id, err)
	}

	return ldapMapper, nil
}

func testKeycloakLdapHardcodedGroupMapper(groupName, groupMapperName string) string {
	return fmt.Sprintf(`
data "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_ldap_user_federation" "openldap" {
	name                    = "openldap"
	realm_id                = data.keycloak_realm.realm.id

	enabled                 = true

	username_ldap_attribute = "cn"
	rdn_ldap_attribute      = "cn"
	uuid_ldap_attribute     = "entryDN"
	user_object_classes     = [
		"simpleSecurityObject",
		"organizationalGroup"
	]
	connection_url          = "ldap://openldap"
	users_dn                = "dc=example,dc=org"
	bind_dn                 = "cn=admin,dc=example,dc=org"
	bind_credential         = "admin"
}

resource "keycloak_group" "hardcoded_group_mapper_test" {
    realm_id    = data.keycloak_realm.realm.id
    name        = "%s"
}

resource "keycloak_ldap_hardcoded_group_mapper" "hardcoded_group_mapper" {
	name                    = "%s"
	realm_id                = data.keycloak_realm.realm.id
	ldap_user_federation_id = keycloak_ldap_user_federation.openldap.id
	group                   = keycloak_group.hardcoded_group_mapper_test.name
}
	`, testAccRealmUserFederation.Realm, groupName, groupMapperName)
}
