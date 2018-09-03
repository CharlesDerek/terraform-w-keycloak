provider "keycloak" {
  client_id = "terraform"
  client_secret = "d0e0f95-0f42-4a63-9b1f-94274655669e"
  url = "http://localhost:8080"
}

resource "keycloak_realm" "test" {
  realm = "test"
}
