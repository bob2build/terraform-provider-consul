# Setup the Consul provisioner to use the demo cluster
provider "consul" {
  address = "127.0.0.1:8500"
  datacenter = "dc1"
  token = "af4c087e-f0b5-4074-b976-47431ddcdc1f"
}

resource "consul_policy" "sensu" {
  policy = "{}"

}

resource "consul_policy" "sensu2" {
  policy = <<EOF
{
  "key": {
    "": {
      "policy": "read"
    },
    "cloudtrust/sensu/": {
      "policy": "write"
    },

    "rabbitmq": {
      "policy": "write"
    }
  }
}
EOF
  type = "management"
  name = "test_terraform"
}
