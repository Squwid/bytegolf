resource "google_compute_network" "bytegolf_local_network" {
  name = "bytegolf-vpc"
  auto_create_subnetworks = false
  routing_mode = "GLOBAL"
  enable_ula_internal_ipv6 = false
}