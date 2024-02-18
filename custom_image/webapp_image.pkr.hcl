packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "project_id" {
  type    = string
  default = "csye6225-dev-414220"
}

variable "zone" {
  type    = string
  default = "us-central1-b"
}

variable "binary_path" {
  type    = string
  description = "The path to the Go binary"
}

source "googlecompute" "centos_stream" {
  project_id              = var.project_id
  zone                    = var.zone
  source_image_family     = "centos-stream-8"
  source_image_project_id = ["centos-cloud"]
  machine_type            = "e2-medium"
  ssh_username            = "packer"
  image_name              = "centos-stream-postgres16-golang-image"
  image_family            = "centos-stream-custom"
}

build {
  sources = ["source.googlecompute.centos_stream"]

  #provisioner "shell" {
  #script = "./custom_image/update_system.sh"
  #}

  provisioner "shell" {
    script = "./custom_image/install_postgresql.sh"
  }

  provisioner "shell" {
    script = "./custom_image/install_golang.sh"
  }
  provisioner "file" {
    source = "${var.binary_path}"
    destination = "/tmp/myapp"
  }

  provisioner "shell" {
  inline = [
    "echo Moving /tmp/myapp to /usr/local/bin",
    "sudo mv /tmp/myapp /usr/local/bin/",
    "sudo chmod +x /usr/local/bin/myapp",
    "echo Listing contents of /usr/local/bin",
    "ls -la /usr/local/bin/"
  ]

  provisioner "file" {
  source      = "./custom_image/myapp.service"
  destination = "/tmp/myapp.service"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/myapp.service /etc/systemd/system/myapp.service",
      "sudo systemctl daemon-reload",
      "sudo systemctl enable myapp.service",
      "sudo systemctl start myapp.service"
    ]
  }

}
}
