
~> **Note** Some network types (for example, Amazon) update the network `name`
and some other attributes a few minutes after creation. This requires using a
`lifecycle` block as shown below. If this `lifecycle` block is missing, then
a subsequent `terraform apply` may attempt to delete the network.

If required, a `lifecycle` block may be added as follows:

```hcl
resource "hpe_morpheus_network" "a" {
  name = "network A"
  display_name = "network A"
  description = "First network"

  .
  .
  .

  lifecycle {
    ignore_changes = [ name, display_name, description ]
  }
}
```
