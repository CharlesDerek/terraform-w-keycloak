---
page_title: "keycloak_openid_user_property_protocol_mapper Resource"
---

# keycloak\_openid\_user\_property\_protocol\_mapper Resource

Allows for creating and managing user property protocol mappers within Keycloak.

User property protocol mappers allow you to map built in properties defined on the Keycloak user interface to a claim in
a token.

Protocol mappers can be defined for a single client, or they can be defined for a client scope which can be shared between
multiple different clients.

## Example Usage (Client)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client" "openid_client" {
  realm_id  = keycloak_realm.realm.id
  client_id = "client"

  name    = "client"
  enabled = true

  access_type         = "CONFIDENTIAL"
  valid_redirect_uris = [
    "http://localhost:8080/openid-callback"
  ]
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper" {
  realm_id  = keycloak_realm.realm.id
  client_id = keycloak_openid_client.openid_client.id
  name      = "user-property-mapper"

  user_property = "email"
  claim_name    = "email"
}
```

## Example Usage (Client Scope)

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client_scope" "client_scope" {
  realm_id = keycloak_realm.realm.id
  name     = "client-scope"
}

resource "keycloak_openid_user_property_protocol_mapper" "user_property_mapper" {
  realm_id        = keycloak_realm.realm.id
  client_scope_id = keycloak_openid_client_scope.client_scope.id
  name            = "test-mapper"

  user_property = "email"
  claim_name    = "email"
}
```

## Argument Reference

- `realm_id` - (Required) The realm this protocol mapper exists within.
- `name` - (Required) The display name of this protocol mapper in the GUI.
- `user_property` - (Required) The built in user property (such as email) to map a claim for.
- `claim_name` - (Required) The name of the claim to insert into a token.
- `client_id` - (Optional) The client this protocol mapper should be attached to. Conflicts with `client_scope_id`. One of `client_id` or `client_scope_id` must be specified.
- `client_scope_id` - (Optional) The client scope this protocol mapper should be attached to. Conflicts with `client_id`. One of `client_id` or `client_scope_id` must be specified. `client_scope_id` - (Required if `client_id` is not specified) The client scope this protocol mapper is attached to.
- `claim_value_type` - (Optional) The claim type used when serializing JSON tokens. Can be one of `String`, `JSON`, `long`, `int`, or `boolean`. Defaults to `String`.
- `add_to_id_token` - (Optional) Indicates if the property should be added as a claim to the id token. Defaults to `true`.
- `add_to_access_token` - (Optional) Indicates if the property should be added as a claim to the access token. Defaults to `true`.
- `add_to_userinfo` - (Optional) Indicates if the property should be added as a claim to the UserInfo response body. Defaults to `true`.

## Import

Protocol mappers can be imported using one of the following formats:
- Client: `{{realm_id}}/client/{{client_keycloak_id}}/{{protocol_mapper_id}}`
- Client Scope: `{{realm_id}}/client-scope/{{client_scope_keycloak_id}}/{{protocol_mapper_id}}`

Example:

```bash
$ terraform import keycloak_openid_user_property_protocol_mapper.user_property_mapper my-realm/client/a7202154-8793-4656-b655-1dd18c181e14/71602afa-f7d1-4788-8c49-ef8fd00af0f4
$ terraform import keycloak_openid_user_property_protocol_mapper.user_property_mapper my-realm/client-scope/b799ea7e-73ee-4a73-990a-1eafebe8e20a/71602afa-f7d1-4788-8c49-ef8fd00af0f4
```
