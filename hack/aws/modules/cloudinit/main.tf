data "template_cloudinit_config" "this" {
  gzip          = !var.debug
  base64_encode = !var.debug

  part {
    filename     = "init.cfg"
    content_type = "text/cloud-config"
    content      = <<-EOT
    repo_update: false
    repo_upgrade: false

    output:
      all: '| tee -a /var/log/cloud-init-output.log'
    EOT
  }

  part {
    filename     = "part_handler.py"
    content_type = "text/part-handler"
    content      = <<-EOT
    #part-handler

    import os

    def list_types():
        return ["text/plain", "text/shellscript"]

    def handle_part(data, ctype, filename, payload):
        if ctype in ["__begin__", "__end__"]:
            return
        if not os.path.exists(os.path.dirname(filename)):
            os.makedirs(os.path.dirname(filename))
        with open(filename, "w") as f:
            f.write(payload)
        if ctype == "text/shellscript":
            os.chmod(filename, 0o755)
    EOT
  }

  dynamic "part" {
    for_each = var.files != null ? fileset(abspath(var.files.src_dir), "**") : []
    content {
      filename     = "${var.files.root_dir}${part.value}"
      content_type = "text/plain"
      content      = templatefile("${abspath(var.files.src_dir)}/${part.value}", var.vars)
    }
  }

  dynamic "part" {
    for_each = var.scripts != null ? fileset(abspath(var.scripts.src_dir), "**") : []
    content {
      filename     = "${var.scripts.dest_dir}/${part.value}"
      content_type = "text/shellscript"
      content      = templatefile("${abspath(var.scripts.src_dir)}/${part.value}", var.vars)
    }
  }

  dynamic "part" {
    for_each = var.scripts != null ? [var.scripts.dest_dir]: []
    content {
      filename     = "run-parts.sh"
      content_type = "text/x-shellscript"
      content      = <<-EOT
      #!/bin/bash

      run-parts --exit-on-error --regex '\.sh$' ${part.value}
      EOT
    }
  }
}
