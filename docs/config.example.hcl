demo = false

logging {
  verbose = false
  level = "info"
}

services {
  http {
    bind = ":5080"
    mount = "/"

    tls {
      cert = "/run/secrets/tls.crt"
      key  = "/run/secrets/tls.key"
    }
  }

  management {
    bind = ":9113"

    tls {
      cert = "/run/secrets/tls.management.crt"
      key  = "/run/secrets/tls.management.key"
    }
  }

  syslog {
    bind = ":5514"
    protocol = "tcp"
    initial_allowed_origins = []

    tls {
      cert = "/run/secrets/tls.syslog.crt"
      key  = "/run/secrets/tls.syslog.key"
    }
  }
}

storage {
  backend "badgerdb" {
    auth_dir = "./data/auth"
    log_dir = "./data/logs"
    config_dir = "./data/config"
  }

  seed {
    auth {
      initial_user = "root"
      initial_password = "root"

      reset_user = ""
      reset_password = ""
    }
  }
}
