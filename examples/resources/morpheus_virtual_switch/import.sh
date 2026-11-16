# Virtual switches are imported by their numeric Morpheus ID. The cluster_id and
# the user-chosen segment map keys are populated by the read that follows the
# import (keys are synthesised from each segment's traffic type and VLAN and may
# be renamed afterwards).
terraform import hpe_morpheus_virtual_switch.example 42
