# Proxy

Register an existing HTTP or HTTPS proxy with HTML/CSS to Image. Use the resource's `id` as `proxy_id` when creating images or configuring templates.

**Resource token:** `html-css-to-image:index:Proxy`

This resource manages the proxy's HCTI configuration. It does not provision or operate the proxy server. Updating its settings keeps the same ID; destroying the resource removes the configuration from HCTI.

The examples use Pulumi YAML, which works with the current local provider without a language SDK. See [provider setup](../../README.md#development) for installation and API credentials. The API key needs `proxies:read`, `proxies:create_update`, and `proxies:delete` to manage the full lifecycle.

## Basic example

Add this resource to your Pulumi YAML program:

```yaml
resources:
  rendering:
    type: html-css-to-image:index:Proxy
    properties:
      name: Rendering proxy
      url: https://proxy.example.com

outputs:
  proxyId: ${rendering.id}
```

Omitting `port` lets the API choose the URL scheme's default: 80 for HTTP or 443 for HTTPS. The input remains unset; `effectivePort` reports the port returned by the API. Omitting `authentication` configures a proxy without credentials.

## Proxy with authentication

```yaml
config:
  proxyPassword:
    type: string
    secret: true

resources:
  rendering:
    type: html-css-to-image:index:Proxy
    properties:
      name: Rendering proxy
      url: https://proxy.example.com
      port: 8080
      authentication:
        username: rendering
        password: ${proxyPassword}
      bypassHosts:
        - static.example.com
        - 192.0.2.10

outputs:
  proxyId: ${rendering.id}
```

Set the secret interactively, then preview the change:

```sh
pulumi config set --secret proxyPassword
pulumi preview
```

Put credentials and port in their separate fields rather than including them in `url`. The password is marked secret in the resource schema and is stored as a secret in Pulumi state.

## Inputs

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `name` | String | Yes | Display name, 3–500 characters after leading and trailing whitespace is removed. |
| `url` | String | Yes | Absolute HTTP or HTTPS proxy URL, up to 512 characters. No embedded credentials, port, query, fragment, or path other than `/`. |
| `port` | Integer | No | From 1 to 65535. Omitted or null lets the API select the scheme's default port. |
| `disabled` | Boolean | No | Disable the proxy without deleting its configuration. Defaults to `false`. |
| `bypassHosts` | Array of strings | No | Up to 100 hostnames, IP addresses, or absolute URLs that should connect directly. Omitted, null, or empty clears the list. |
| `authentication` | Object | No | Credentials described below. Omitted or null removes authentication. |

### Authentication

| Input | Type | Required | Description |
| --- | --- | --- | --- |
| `username` | String | Yes, within `authentication` | Up to 512 characters. An empty string is valid. Whitespace is preserved. |
| `password` | Secret string | On creation or username changes | Up to 484 UTF-8 bytes. An empty string sets an empty password. Whitespace is preserved. Omitted or null on an existing proxy retains its password when the username is unchanged. |

Password retention is handled by the provider. There is no `retainPassword` input on the Pulumi resource.

## Outputs

The resource also exposes its configured inputs and these read-only outputs:

| Output | Type | Description |
| --- | --- | --- |
| `id` | String | HCTI proxy ID. Use this for `proxy_id` in image or template settings. |
| `effectivePort` | Integer | Port returned by the API; may be absent if the response does not specify one. |
| `createdAt` | String | Creation timestamp in RFC 3339 format. |
| `updatedAt` | String | Last update timestamp in RFC 3339 format. |

## Updating credentials

For an existing authenticated proxy, this configuration retains the stored password while changing its name:

```yaml
resources:
  rendering:
    type: html-css-to-image:index:Proxy
    properties:
      name: Renamed rendering proxy
      url: https://proxy.example.com
      port: 8080
      authentication:
        username: rendering
      bypassHosts:
        - static.example.com
        - 192.0.2.10
```

Keep the rest of your desired configuration in place: updates replace settings, so omitted ordinary options are cleared or reset.

| Intended change | Configuration |
| --- | --- |
| Replace the password | Supply the new `authentication.password`. |
| Set an empty password | Set `password: ""` inside `authentication`. |
| Keep the stored password | Keep the same username and omit `password`, or set it to null. |
| Change the username | Supply the new username and a password together. |
| Remove authentication | Remove the entire `authentication` object, or set it to null. |

A retained password must already exist. The same rules apply when the proxy is disabled. If the username changes outside Pulumi, refresh exposes the new username and a subsequent change back requires a password.

## Bypass hosts

Bypass entries are matched by hostname. For example, `STATIC.EXAMPLE.COM` and `https://static.example.com/assets` both normalize to `static.example.com`. URL paths are not bypass rules. Hosts are lowercased and deduplicated.

The provider preserves equivalent configured values when refreshing, so server normalization alone does not cause another update. Omitting `bypassHosts` or setting it to `[]` clears the complete list.

## Import

With the provider available locally, import an existing HCTI proxy by its management ID:

```sh
pulumi import html-css-to-image:index:Proxy rendering PROXY_ID
```

Add the imported resource to your program using its current name, URL, port, disabled state, bypass hosts, and username. The API never returns the password, so import leaves it unset. Keep `authentication` with the existing username and omit `password` to retain it during later updates. Removing the authentication object would remove credentials on the next apply.

Import records the port returned by the API. Keep that value in your program for an unchanged preview, or intentionally omit it to reset to the scheme default on an update.

Run `pulumi preview` after adding the resource to your program to review any differences before applying.

## Preview, refresh, and deletion

Preview performs no writes. Inputs that are not yet known remain unknown until apply. Refresh reads metadata without rendering an image or testing connectivity to the proxy. An unchanged deployment performs no writes.

If HCTI reports the proxy missing with HTTP 404, refresh removes it from state. Authentication, authorization, rate-limit, and other errors are reported without treating the proxy as deleted.

Destroy deletes the HCTI proxy configuration. The external proxy server remains running. Deleting a configuration that is already missing succeeds. Review image and template references before removing a proxy they use.
