packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "project_id" {
  type = string
  #description = "The Projetc ID"
  default = "csye6225-dev-414220"
}

variable "zone" {
  type = string
  #description = "The zone in the GCP"
  default = "us-central1-b"
}

variable "binary_path" {
  type        = string
  description = "The path to the Go binary"
}

variable "db_user" {
  type        = string
  description = "Database User"
}

variable "db_password" {
  type        = string
  description = "Database Password"
}

variable "db_name" {
  type        = string
  description = "Database Name"
}

locals {
  timestamp = "${lower(regex_replace(timestamp(), "[-:TZ]", "-"))}"
}


source "googlecompute" "centos_stream" {
  project_id              = var.project_id
  zone                    = var.zone
  source_image_family     = "centos-stream-8"
  source_image_project_id = ["centos-cloud"]
  machine_type            = "e2-medium"
  ssh_username            = "packer"
  image_name              = "webapp-golden-${local.timestamp}-image"
  image_family            = "centos-stream-custom"
}

build {
  sources = [
    "source.googlecompute.centos_stream"
  ]

  provisioner "shell" {
    inline = [
      "echo Creating group csye6225",
      "sudo groupadd csye6225",
      "echo Creating user csye6225 with no login shell",
      "sudo useradd -r -g csye6225 -s /usr/sbin/nologin csye6225",
      "echo User and group csye6225 created successfully"
    ]
  }

  #provisioner "shell" {
  #script = "./custom_image/update_system.sh"
  #}

  provisioner "shell" {
    script = "./custom_image/install_postgresql.sh"
    environment_vars = [
      "DB_USER=${var.db_user}",
      "DB_PASSWORD=${var.db_password}",
      "DB_NAME=${var.db_name}"
    ]
  }

  provisioner "shell" {
    script = "./custom_image/install_golang.sh"
  }

  provisioner "file" {
    source      = var.binary_path
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    inline = [
      "echo Installing policycoreutils-python-utils for semanage",
      "sudo dnf install -y policycoreutils-python-utils",
      "echo policycoreutils-python-utils installed successfully"
    ]
  }

  provisioner "shell" {
    inline = [
      "echo Moving /tmp/webapp to /usr/local/bin",
      "sudo mv /tmp/webapp /usr/local/bin/webapp",
      "sudo chmod +x /usr/local/bin/webapp",
      "sudo chown csye6225:csye6225 /usr/local/bin/webapp",
      "echo Setting SELinux context for /usr/local/bin/webapp",
      "sudo semanage fcontext -a -t bin_t '/usr/local/bin/webapp'",
      "sudo restorecon -v '/usr/local/bin/webapp'",
      "echo Listing contents of /usr/local/bin",
      "ls -la /usr/local/bin/"
    ]
  }

  provisioner "file" {
    source      = "./custom_image/webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/webapp.service /etc/systemd/system/webapp.service",
      "sudo chown csye6225:csye6225 /etc/systemd/system/webapp.service",
      "sudo systemctl daemon-reload",
    ]
  }

  provisioner "shell" {
    script = "./custom_image/replace_envs.sh"
    environment_vars = [
      "DB_USER=${var.db_user}",
      "DB_HOST=localhost",
      "DB_PASSWORD=${var.db_password}",
      "DB_NAME=${var.db_name}"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo systemctl enable webapp.service",
      "sudo systemctl start webapp.service"
    ]
  }

  provisioner "file" {
    source      = "./custom_image/restart_webapp.sh"
    destination = "/tmp/restart_webapp.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/restart_webapp.sh /usr/local/bin/restart_webapp.sh",
      "sudo chmod +x /usr/local/bin/restart_webapp.sh",
      "sudo chown csye6225:csye6225 /usr/local/bin/restart_webapp.sh",
      "sudo semanage fcontext -a -t bin_t '/usr/local/bin/restart_webapp.sh'",
      "sudo restorecon -v '/usr/local/bin/restart_webapp.sh'"
    ]
  }

  provisioner "file" {
    source      = "./custom_image/restart_webapp.service"
    destination = "/tmp/restart_webapp.service"
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/restart_webapp.service /etc/systemd/system/restart_webapp.service",
      "sudo chown csye6225:csye6225 /etc/systemd/system/restart_webapp.service",
      "sudo systemctl enable restart_webapp.service"
    ]
  }
}
