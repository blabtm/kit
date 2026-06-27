package kit

import (
  "lab.kit/v1/dom/kit/env"
  "lab.kit/v1/dom/kit/svc"
  "lab.kit/v1/dom/kit/rp"
  "lab.kit/v1/dom/kit/mq"
  "lab.kit/v1/dom/kit/db"
)

#Kit: {
  cluster: env.#Cluster

  #validator: env.#ClusterValidator & {
    #cluster: cluster
  }

  service: {
    [string]: svc.#Config & {
      deployment: svc.#Deployment & {
        #AllowedNode: #validator.#AllowedNode
      }
    }
  }

  service: {
    "db": db.#Config
    "rp": rp.#Config & {
      #cluster: cluster
    }
    "mq": mq.#Config
  }
}
