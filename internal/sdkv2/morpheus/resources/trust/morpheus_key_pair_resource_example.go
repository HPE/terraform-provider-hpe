// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package trust

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_key_pair/resource.tf morpheus_key_pair_resource.tf.tmpl Name example-key-pair PublicKey "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC..." PrivateKey "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...\n-----END RSA PRIVATE KEY-----"
