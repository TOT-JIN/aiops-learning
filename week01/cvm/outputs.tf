output "public_ip" {
  description = "vm public ip address"
  value       = tencentcloud_instance.docker_vm.public_ip
}
