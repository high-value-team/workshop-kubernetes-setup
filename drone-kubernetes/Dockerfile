FROM alpine
LABEL maintainer="Florian Boehmak <florian@fnbk.cc>"

ENV HELM_VERSION="v2.9.1"
ENV KUBECTL_VERSION="v1.9.6"

RUN \
  apk add --update ca-certificates && \
  apk add -t deps curl && \
  apk add bash && \
  apk add --update libintl && \
  apk add --virtual build_deps gettext &&  \
  cp /usr/bin/envsubst /usr/local/bin/envsubst && \
  apk del build_deps

RUN \
  curl -Lo /tmp/helm.tar.gz \
    https://storage.googleapis.com/kubernetes-helm/helm-${HELM_VERSION}-linux-amd64.tar.gz && \
  tar -xvf /tmp/helm.tar.gz -C /tmp && \
  mv /tmp/linux-amd64/helm /usr/local/bin && \
  rm -rf /tmp/*

RUN \
  curl -Lo /usr/local/bin/kubectl \
    https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl && \
  chmod +x /usr/local/bin/kubectl && \
  apk del --purge deps && \
  rm /var/cache/apk/*

COPY template.kubeconfig /template.kubeconfig

COPY run.sh /run.sh
RUN chmod +x /run.sh

CMD ["/run.sh"]