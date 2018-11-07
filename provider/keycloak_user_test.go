package provider

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/charlesderek/terraform-w-keycloak/keycloak"
	"testing"
)

func TestAccKeycloakUser_basic(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_basic(realmName, username),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
			{
				ResourceName:        resourceName,
				ImportState:         true,
				ImportStateVerify:   true,
				ImportStateIdPrefix: realmName + "/",
			},
		},
	})
}

func TestAccKeycloakUser_createAfterManualDestroy(t *testing.T) {
	var user = &keycloak.User{}

	realmName := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_basic(realmName, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					testAccCheckKeycloakUserFetch(resourceName, user),
				),
			},
			{
				PreConfig: func() {
					keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

					err := keycloakClient.DeleteUser(user.RealmId, user.Id)
					if err != nil {
						t.Fatal(err)
					}
				},
				Config: testKeycloakUser_basic(realmName, username),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakUser_updateRealm(t *testing.T) {
	realmOne := "terraform-" + acctest.RandString(10)
	realmTwo := "terraform-" + acctest.RandString(10)
	username := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_updateRealmBefore(realmOne, realmTwo, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "realm_id", realmOne),
				),
			},
			{
				Config: testKeycloakUser_updateRealmAfter(realmOne, realmTwo, username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "realm_id", realmTwo),
				),
			},
		},
	})
}

func TestAccKeycloakUser_updateUsername(t *testing.T) {
	realmName := "terraform-" + acctest.RandString(10)
	usernameOne := "terraform-user-" + acctest.RandString(10)
	usernameTwo := "terraform-user-" + acctest.RandString(10)

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_basic(realmName, usernameOne),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "username", usernameOne),
				),
			},
			{
				Config: testKeycloakUser_basic(realmName, usernameTwo),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "username", usernameTwo),
				),
			},
		},
	})
}

func TestAccKeycloakUser_updateInPlace(t *testing.T) {
	userOne := &keycloak.User{
		RealmId:   "terraform-" + acctest.RandString(10),
		Username:  "terraform-user-" + acctest.RandString(10),
		Email:     fmt.Sprintf("%s@gmail.com", acctest.RandString(10)),
		FirstName: acctest.RandString(10),
		LastName:  acctest.RandString(10),
		Enabled:   randomBool(),
	}

	userTwo := &keycloak.User{
		RealmId:   userOne.RealmId,
		Username:  userOne.Username,
		Email:     fmt.Sprintf("%s@gmail.com", acctest.RandString(10)),
		FirstName: acctest.RandString(10),
		LastName:  acctest.RandString(10),
		Enabled:   randomBool(),
	}

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_fromInterface(userOne),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
			{
				Config: testKeycloakUser_fromInterface(userTwo),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
		},
	})
}

func TestAccKeycloakUser_unsetOptionalAttributes(t *testing.T) {
	userWithOptionalAttributes := &keycloak.User{
		RealmId:   "terraform-" + acctest.RandString(10),
		Username:  "terraform-user-" + acctest.RandString(10),
		Email:     fmt.Sprintf("%s@gmail.com", acctest.RandString(10)),
		FirstName: acctest.RandString(10),
		LastName:  acctest.RandString(10),
		Enabled:   randomBool(),
	}

	resourceName := "keycloak_user.user"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckKeycloakUserDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testKeycloakUser_fromInterface(userWithOptionalAttributes),
				Check:  testAccCheckKeycloakUserExists(resourceName),
			},
			{
				Config: testKeycloakUser_basic(userWithOptionalAttributes.RealmId, userWithOptionalAttributes.Username),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKeycloakUserExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "email", ""),
					resource.TestCheckResourceAttr(resourceName, "first_name", ""),
					resource.TestCheckResourceAttr(resourceName, "last_name", ""),
				),
			},
		},
	})
}

func testAccCheckKeycloakUserExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, err := getUserFromState(s, resourceName)
		if err != nil {
			return err
		}

		return nil
	}
}

func testAccCheckKeycloakUserFetch(resourceName string, user *keycloak.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		fetchedUser, err := getUserFromState(s, resourceName)
		if err != nil {
			return err
		}

		user.Id = fetchedUser.Id
		user.RealmId = fetchedUser.RealmId

		return nil
	}
}

func testAccCheckKeycloakUserDestroy() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "keycloak_user" {
				continue
			}

			id := rs.Primary.ID
			realm := rs.Primary.Attributes["realm_id"]

			keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

			user, _ := keycloakClient.GetUser(realm, id)
			if user != nil {
				return fmt.Errorf("user with id %s still exists", id)
			}
		}

		return nil
	}
}

func getUserFromState(s *terraform.State, resourceName string) (*keycloak.User, error) {
	keycloakClient := testAccProvider.Meta().(*keycloak.KeycloakClient)

	rs, ok := s.RootModule().Resources[resourceName]
	if !ok {
		return nil, fmt.Errorf("resource not found: %s", resourceName)
	}

	id := rs.Primary.ID
	realm := rs.Primary.Attributes["realm_id"]

	user, err := keycloakClient.GetUser(realm, id)
	if err != nil {
		return nil, fmt.Errorf("error getting user with id %s: %s", id, err)
	}

	return user, nil
}

func testKeycloakUser_basic(realm, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id = "${keycloak_realm.realm.id}"
	username = "%s"
}
	`, realm, username)
}

func testKeycloakUser_updateRealmBefore(realmOne, realmTwo, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id  = "${keycloak_realm.realm_1.id}"
	username  = "%s"
}
	`, realmOne, realmTwo, username)
}

func testKeycloakUser_updateRealmAfter(realmOne, realmTwo, username string) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm_1" {
	realm = "%s"
}

resource "keycloak_realm" "realm_2" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id  = "${keycloak_realm.realm_2.id}"
	username  = "%s"
}
	`, realmOne, realmTwo, username)
}

func testKeycloakUser_fromInterface(user *keycloak.User) string {
	return fmt.Sprintf(`
resource "keycloak_realm" "realm" {
	realm = "%s"
}

resource "keycloak_user" "user" {
	realm_id   = "${keycloak_realm.realm.id}"
	username   = "%s"

	email      = "%s"
	first_name = "%s"
	last_name  = "%s"
	enabled    = %t
}
	`, user.RealmId, user.Username, user.Email, user.FirstName, user.LastName, user.Enabled)
}
