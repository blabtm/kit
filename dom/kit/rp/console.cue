package rp

#ConsoleConfig: {
  #kafka: {
    host: string
    port: int
  }

  #secret: #Secret

  kafka: {
    brokers: ["\(#kafka.host):\(#kafka.port)"]
  }

  schemaRegistry: {
    enabled: false
  }

  redpanda: {
    adminApi: {
      enabled: true
      urls: ["http://\(#kafka.host):9644"]
      authentication: {
        basic: {
          username: #secret.user
          password: #secret.pass
        }
        impersonateUser: false
      }
    }
  }
}
