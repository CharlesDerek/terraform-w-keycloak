# keycloak_custom_user_federation

Allows for creating and managing custom user federation providers within Keycloak.

A custom user federation provider is an implementation of Keycloak's
[User Storage SPI](https://www.keycloak.org/docs/4.2/server_development/index.html#_user-storage-spi).
An example of this implementation can be found [here](https://github.com/charlesderek/terraform-w-keycloak/tree/master/custom-user-federation-example).

### Example Usage

```hcl
resource "keycloak_realm" "realm" {
    realm   = "test"
    enabled = true
}

resource "keycloak_custom_user_federation" "custom_user_federation" {
    name        = "custom"
    realm_id    = "${keycloak_realm.realm.id}"
    provider_id = "custom"

    enabled     = true
}
```

### Argument Reference

The following arguments are supported:

- `realm_id` - (Required) The realm that this provider will provide user federation for.
- `name` - (Required) Display name of the provider when displayed in the console.
- `provider_id` - (Required) The unique ID of the custom provider, specified in the `getId` implementation for the `UserStorageProviderFactory` interface.
- `enabled` - (Optional) When `false`, this provider will not be used when performing queries for users. Defaults to `true`.
- `priority` - (Optional) Priority of this provider when looking up users. Lower values are first. Defaults to `0`.
- `cache_policy` - (Optional) Can be one of `DEFAULT`, `EVICT_DAILY`, `EVICT_WEEKLY`, `MAX_LIFESPAN`, or `NO_CACHE`. Defaults to `DEFAULT`.
- `parent_id` - (Optional) Must be set to the realms' `internal_id`  when it differs from the realm. This can happen when existing resources are imported into the state.

### Import

Custom user federation providers can be imported using the format `{{realm_id}}/{{custom_user_federation_id}}`.
The ID of the custom user federation provider can be found within the Keycloak GUI and is typically a GUID:

```bash
$ terraform import keycloak_custom_user_federation.custom_user_federation my-realm/af2a6ca3-e4d7-49c3-b08b-1b3c70b4b860
```
