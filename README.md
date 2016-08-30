[![Buildstatus](https://travis-ci.org/universityofderby/terraform-provider-infoblox.svg)](https://travis-ci.org/universityofderby/terraform-provider-infoblox)

# [Terraform](https://github.com/hashicorp/terraform) Infoblox Provider

The Terraform Infoblox provider is used to interact with the
resources supported by Infoblox. The provider needs to be configured
with the proper credentials before it can be used.

##  Download
Download builds for Darwin, Linux and Windows from the [releases page](https://github.com/universityofderby/terraform-provider-infoblox/releases/).

## Example Usage

```
# Configure the Infoblox provider
provider "infoblox" {
    username = "${var.infoblox_username}"
    password = "${var.infoblox_password}"
    host  = "${var.infoblox_host}"
    sslverify = "${var.infoblox_sslverify}"
    usecookies = "${var.infoblox_usecookies}"
}

# Create a record
resource "infoblox_record" "www" {
    ...
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) The Infoblox username. It must be provided, but it can also be sourced from the `INFOBLOX_USERNAME` environment variable.
* `password` - (Required) The password associated with the username. It must be provided, but it can also be sourced from the `INFOBLOX_PASSWORD` environment variable.
* `host` - (Required) The base url for the Infoblox REST API, but it can also be sourced from the `INFOBLOX_HOST` environment variable.
* `sslverify` - (Required) Enable ssl for the REST api, but it can also be sourced from the `INFOBLOX_SSLVERIFY` environment variable.
* `usecookies` - (Optional) Use cookies to connect to the REST API, but it can also be sourced from the `INFOBLOX_USECOOKIES` environment variable

# infoblox\_record

Provides a Infoblox record resource.

## Example Usage

```
# Create A record
resource "infoblox_record" "foobar" {
	name = "terraform"
	domain = "mydomain.com"
	type = "A"
	value = "192.168.0.10"
	ttl = 3600
}
```

```
# Create A record using nextavailableip
resource "infoblox_record" "foobar" {
	name = "terraform"
	domain = "mydomain.com"
	nextavailableip = true
	value = "192.168.0.0/24"
	type = "A"
}
```

```
# Create host record using nextavailableip
resource "infoblox_record" "foobar" {
	name = "terraform"
	domain = "mydomain.com"
	nextavailableip = true
	value = "192.168.0.0/24"
	type = "Host"
	comment = "Terraform test record"
}
```

## Argument Reference

See [related part of Infoblox Docs](https://godoc.org/github.com/fanatic/go-infoblox) for details about valid values.

The following arguments are supported:

* `comment` - (Optional) The comment of the record
* `domain` - (Required) The domain to add the record to
* `name` - (Required) The name of the record
* `nextavailableip` - (Boolean, Optional) Get next available IP address using CIDR or other search
* `ttl` - (Integer, Optional) The TTL of the record
* `type` - (Required) The type of the record
* `value` - (Required) The value of the record; its usage will depend on the `type` (see below)
* `view` - (Optional) The view of the record

## DNS Record Types

The type of record being created affects the interpretation of the `value` argument.

#### A Record

* `value` is the IPv4 address, or the search value (e.g. CIDR) if nextavailableip = true

#### AAAA Record

* `value` is the IPv6 address, or the search value (e.g. CIDR) if nextavailableip = true

#### CNAME Record

* `value` is the canonical name

#### Host Record

* `value` is the IPv4 address, or the search value (e.g. CIDR) if nextavailableip = true

## Attributes Reference

The following attributes are exported:

* `comment` - The comment of the record [string]
* `domain` - The domain of the record [string]
* `fqdn` - The FQDN (computed) of the record [string]
* `ipv4addr` - The IPv4 address (computed) of the record [string]
* `ipv6addr` - The IPv6 address (computed) of the record [string]
* `name` - The name of the record [string]
* `nextavailableip` - Get next available IP address using CIDR or other search [bool]
* `type` - The type of the record [string]
* `ttl` - The TTL of the record [int]
* `value` - The value of the record [string]
* `view` - The view of the record [string]

# infoblox\_ip

Queries the next available IP address from a network and returns it in a computed variable
that can be used by the infoblox_record resource.

## Example Usage

```
# Acquire the next available IP from a network CIDR
#it will create a variable called "ipaddress"
resource "infoblox_ip" "theIPAddress" {
	cidr = "10.0.0.0/24"
}


# Add a record to the domain
resource "infoblox_record" "foobar" {
	value = "${infoblox_ip.theIPAddress.ipaddress}"
	name = "terraform"
	domain = "mydomain.com"
	type = "A"
	ttl = 3600
}
```

## Argument Reference

The following arguments are supported:

* `cidr` - (Required) The network to search for - example 10.0.0.0/24
