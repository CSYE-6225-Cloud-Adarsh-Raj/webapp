packer {
  required_plugins {
    googlecompute = {
      source  = "github.com/hashicorp/googlecompute"
      version = "~> 1"
    }
  }
}

variable "project_id" {
  type        = string
  description = "The Projetc ID"
  default     = "csye6225-dev-414220"
}

variable "zone" {
  type        = string
  description = "The zone in the GCP"
  default     = "us-central1-b"
}

variable "binary_path" {
  type        = string
  description = "The path to the Go binary"
  default     = "../webapp"
}

#variable "db_user" {
#  type        = string
#  description = "Database User"
#}

#variable "db_password" {
#  type        = string
#  description = "Database Password"
#}

#variable "db_name" {
#  type        = string
#  description = "Database Name"
#}

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
    script = "./create_user.sh"
  }

  provisioner "shell" {
    script = "./update_system.sh"
  }

  #provisioner "shell" {
  #  script = "./install_postgresql.sh"
  #  environment_vars = [
  #    "DB_USER=${var.db_user}",
  #    "DB_PASSWORD=${var.db_password}",
  #    "DB_NAME=${var.db_name}"
  #  ]
  #}

  provisioner "shell" {
    script = "./install_opsagent.sh"
  }

  provisioner "file" {
    source      = "config.yaml"
    destination = "/tmp/config.yaml"
  }

  provisioner "shell" {
    script = "./setup_opsagent.sh"
  }

  provisioner "shell" {
    script = "./restart_opsagent.sh"
  }

  provisioner "shell" {
    script = "./install_golang.sh"
  }

  provisioner "file" {
    source      = var.binary_path
    destination = "/tmp/webapp"
  }

  provisioner "shell" {
    script = "./semanage.sh"
  }

  provisioner "shell" {
    script = "./wb_system_p.sh"
  }

  provisioner "shell" {
    inline = [
      "sudo mkdir -p /etc",
      "echo '# Placeholder for environment variables' | sudo tee /etc/webapp.env"
    ]
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  provisioner "shell" {
    script = "./wb_user.sh"
  }

  #provisioner "shell" {
  #  script = "./replace_envs.sh"
  #  environment_vars = [
  #    "DB_USER=${var.db_user}",
  #    "DB_HOST=localhost",
  #    "DB_PASSWORD=${var.db_password}",
  #    "DB_NAME=${var.db_name}"
  #  ]
  #}

  provisioner "shell" {
    script = "./wb_start.sh"
  }

  provisioner "file" {
    source      = "./restart_webapp.sh"
    destination = "/tmp/restart_webapp.sh"
  }

  provisioner "shell" {
    script = "./restart_wb_p.sh"
  }

  provisioner "file" {
    source      = "./restart_webapp.service"
    destination = "/tmp/restart_webapp.service"
  }

  provisioner "shell" {
    script = "./restart_wb.sh"
  }
}
