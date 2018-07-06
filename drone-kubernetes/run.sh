#!/bin/sh

ok=1

# environment variables
if [[ -z ${KUBERNETES_SERVER} ]]; then echo KUBERNETES_SERVER not set; ok=0; fi
if [[ -z ${KUBERNETES_CERTIFICATE_AUTHORITY_DATA} ]]; then echo KUBERNETES_CERTIFICATE_AUTHORITY_DATA not set; ok=0; fi
if [[ -z ${KUBERNETES_CLIENT_CERTIFICATE_DATA} ]]; then echo KUBERNETES_CLIENT_CERTIFICATE_DATA not set; ok=0; fi
if [[ -z ${KUBERNETES_CLIENT_KEY_DATA} ]]; then echo KUBERNETES_CLIENT_KEY_DATA not set; ok=0; fi

echo; echo checking Environment variables
if [[ $ok -eq 0 ]]; then
    printf 'Invalid parameters!\n'
    exit 1;
fi

# optional variables
if [[ ! -z ${PLUGIN_KUBECTL} ]]; then echo option PLUGIN_KUBECTL has been set; fi
if [[ ! -z ${PLUGIN_HELM} ]]; then echo option PLUGIN_HELM has been set; fi
if [[ ! -z ${PLUGIN_DEPLOYMENT} ]]; then echo option PLUGIN_DEPLOYMENT has been set; fi
if [[ ! -z ${DRONE_WORKSPACE} ]]; then echo option DRONE_WORKSPACE has been set; fi


# create kubeconfig
envsubst < /template.kubeconfig > /kubeconfig

## Run kubectl command
if [[ ! -z ${PLUGIN_KUBECTL} ]]; then
  kubectl --kubeconfig=/kubeconfig ${PLUGIN_KUBECTL}
fi

# Run helm command
if [[ ! -z ${PLUGIN_HELM} ]]; then
  helm --kubeconfig=/kubeconfig ${PLUGIN_HELM}
fi

# Run helm command
if [[ ! -z ${PLUGIN_DEPLOYMENT} ]]; then
  kubectl --kubeconfig=/kubeconfig apply -f ${DRONE_WORKSPACE}/${PLUGIN_DEPLOYMENT}
fi


