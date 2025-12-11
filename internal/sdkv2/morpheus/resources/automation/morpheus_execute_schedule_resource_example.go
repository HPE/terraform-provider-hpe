// (C) Copyright 2025 Hewlett Packard Enterprise Development LP

package automation

//go:generate go run ../../../../../cmd/render -out examples/resources/morpheus_execute_schedule/resource.tf hpe_morpheus_execute_schedule_resource.tf.tmpl Name "Run daily at 7 AM" Description "This schedule runs daily at 7 AM Mountain Time" Enabled false TimeZone "America/Denver" Schedule "7 0 * * *"
