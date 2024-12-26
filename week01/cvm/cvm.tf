provider "tencentcloud" {
  region     = var.tencentcloud_region
}

resource "tencentcloud_security_group" "docker_sg" {
  name        = "docker-security-group"
  description = "Security group for Docker"
}

resource "tencentcloud_security_group_lite_rule" "default" {
  security_group_id = tencentcloud_security_group.docker_sg.id
  ingress = [
    "ACCEPT#0.0.0.0/0#22#TCP",
    "ACCEPT#0.0.0.0/0#6443#TCP",
  ]

  egress = [
    "ACCEPT#0.0.0.0/0#ALL#ALL"
  ]
}
resource "tencentcloud_key_pair" "default" {
  key_name   = "skey_8psgl2e3"
  public_key = file("~/.ssh/id_rsa.pub")
}

data "tencentcloud_images" "default" {
  image_type = ["PUBLIC_IMAGE"]
  os_name    = "ubuntu"
}

resource "tencentcloud_instance" "docker_vm" {
  instance_name       = "terraform-docker-vm"
  availability_zone   = var.tencentcloud_availability_zone
  instance_type       = "SA5.MEDIUM2"
  image_id            = data.tencentcloud_images.default.images.0.image_id 
  system_disk_type    = "CLOUD_BSSD"
  system_disk_size    = 50
  allocate_public_ip         = true
  internet_max_bandwidth_out = 100

  key_ids                 = [tencentcloud_key_pair.default.id]
  orderly_security_groups = [tencentcloud_security_group.docker_sg.id]

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      host        = self.public_ip
      user        = "ubuntu"
      private_key = file("~/.ssh/id_rsa")
    }

    inline = [
      "sudo apt-get update -y",
      "sudo apt-get install -qqy --no-install-recommends docker.io",
      "sudo systemctl enable docker",
      "sudo systemctl start docker"
    ]
  }
}
