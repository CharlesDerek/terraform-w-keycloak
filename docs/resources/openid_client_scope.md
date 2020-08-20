---
page_title: "keycloak_openid_client_scope Resource"
---

# keycloak\_openid\_client\_scope Resource

Allows for creating and managing Keycloak client scopes that can be attached to clients that use the OpenID Connect protocol.

Client Scopes can be used to share common protocol and role mappings between multiple clients within a realm. They can also
be used by clients to conditionally request claims or roles for a user based on the OAuth 2.0 `scope` parameter.

## Example Usage

```hcl
resource "keycloak_realm" "realm" {
  realm   = "my-realm"
  enabled = true
}

resource "keycloak_openid_client_scope" "openid_client_scope" {
  realm_id               = keycloak_realm.realm.id
  name                   = "groups"
  description            = "When requested, this scope will map a user's group memberships to a claim"
  include_in_token_scope = true
  gui_order              = 1
}
```

## Argument Reference

- `realm_id` - (Required) The realm this client scope belongs to.
- `name` - (Required) The display name of this client scope in the GUI.
- `description` - (Optional) The description of this client scope in the GUI.
- `consent_screen_text` - (Optional) When set, a consent screen will be displayed to users authenticating to clients with this scope attached. The consent screen will display the string value of this attribute.
- `include_in_token_scope` - (Optional) When `true`, the name of this client scope will be added to the access token property 'scope' as well as to the Token Introspection Endpoint response.
- `gui_order` - (Optional) Specify order of the client scope in GUI (such as in Consent page) as integer.

## Import

Client scopes can be imported using the format `{{realm_id}}/{{client_scope_id}}`, where `client_scope_id` is the unique ID that Keycloak
assigns to the client scope upon creation. This value can be found in the URI when editing this client scope in the GUI, and is typically a GUID.

Example:

```bash
$ terraform import keycloak_openid_client_scope.openid_client_scope my-realm/8e8f7fe1-df9b-40ed-bed3-4597aa0dac52
```
