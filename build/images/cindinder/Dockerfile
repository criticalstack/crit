FROM docker:dind
ARG KUBECTL_VERSION="v1.18.3"
WORKDIR /cinder
COPY build/ images/
RUN apk update
RUN apk add --no-cache curl make git go
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl
RUN chmod +x /usr/local/bin/kubectl
RUN curl https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | ash
COPY bin/cinder /usr/local/bin/cinder
COPY entrypoint.sh .
ENTRYPOINT ["/bin/ash"]
