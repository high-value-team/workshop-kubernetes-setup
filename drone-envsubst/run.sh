#!/bin/bash

ok=1

# secrets
if [[ -z ${AWS_DEFAULT_REGION} ]]; then echo AWS_DEFAULT_REGION not set; ok=0; fi
if [[ -z ${ECR_REPOSITORY_ID} ]]; then echo ECR_REPOSITORY_ID not set; ok=0; fi

# plugin variables
if [[ -z ${PLUGIN_SOURCE} ]]; then echo PLUGIN_SOURCE not set; ok=0; fi
if [[ -z ${PLUGIN_DESTINATION} ]]; then echo PLUGIN_DESTINATION not set; ok=0; fi

# drone variables
if [[ -z ${DRONE_WORKSPACE} ]]; then echo DRONE_WORKSPACE not set; ok=0; fi

echo; echo checking Environment variables
if [[ $ok -eq 0 ]]; then
    printf 'Invalid parameters!\n'
    exit 1;
fi


source=${DRONE_WORKSPACE}/${PLUGIN_SOURCE}
destination=${DRONE_WORKSPACE}/${PLUGIN_DESTINATION}

envsubst < ${source} > ${destination}
