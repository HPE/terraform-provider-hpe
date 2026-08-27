variable "password" {
  description = "Password for the admin user"
  type        = string
  sensitive   = true
}

variable "email" {
  description = "Email for the admin user"
  type        = string
}