ARG USER_NAME='app'
ARG USER_UID='1000'
ARG USER_GID='1000'
ARG USER_GROUP='app_users'
ARG REPO='github.com/kohirens/tmplpress'

FROM kohirens/go:1.21.0 AS dev

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

ENV GOPATH /home/${USER_NAME}
ENV CGO_ENABLED=0
ENV WORK_DIR=/home/${USER_NAME}/src/${REPO}
ENV SHELL /bin/bash

COPY .docker/start.sh /usr/local/bin/start.sh

# VS Code Requirements: openssh, musl, libgcc, libstdc++, git
RUN apk --no-progress --purge --no-cache upgrade \
 && apk --no-progress --purge --no-cache add --upgrade \
    curl \
    git \
    gnupg \
    libgcc \
    libstdc++ \
    openssh \
 && apk --no-progress --purge --no-cache upgrade \
 && rm -vrf /var/cache/apk/* \
 && rm -rf /tmp/* \
 && curl --version \
 && git --version

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin ~/src/${REPO} && chmod 0774 -R ~/src/${REPO}

WORKDIR "/home/${USER_NAME}"

# # Install Go helpful dev tools.
RUN go install golang.org/x/tools/gopls@latest
#  && go install github.com/go-delve/delve/cmd/dlv@latest \
#  && go install github.com/ramya-rao-a/go-outline@latest \
#  && go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest \
#  && go install honnef.co/go/tools/cmd/staticcheck@latest \
#  && go install github.com/cweill/gotests/...@latest

COPY --chown=${USER_NAME}:${USER_GROUP} .devcontainer/download-vs-code-server.sh .

# Install VS Code Server and Requirements for pre-installing extensions
RUN chmod +x ./download-vs-code-server.sh && ./download-vs-code-server.sh \
 && mkdir -p /home/${USER_NAME}/.vscode-server/extensions \
 && rm ./download-vs-code-server.sh

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

WORKDIR "/home/${USER_NAME}/src/${REPO}"

ENTRYPOINT [ "start.sh" ]

CMD [ "" ]

FROM dev AS build

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP
ARG REPO

COPY --chown=${USER_NAME}:${USER_GROUP} . .

RUN go mod tidy && go build && go install

FROM alpine:3.18 AS release

ARG USER_NAME
ARG USER_UID
ARG USER_GID
ARG USER_GROUP

# Add a non-root group and user, helpful if you dev on Linux.
RUN addgroup --system --gid ${USER_GID} ${USER_GROUP} \
 && adduser --system \
    --disabled-password \
    --ingroup ${USER_GROUP} \
    --uid ${USER_UID} \
    ${USER_NAME}

USER ${USER_NAME}

# Make directories with correct permissions
RUN mkdir -p ~/bin

COPY --from=build --chown=${USER_NAME}:${USER_GROUP} /home/${USER_NAME}/bin/tmplpress /home/${USER_NAME}/bin

ENV PATH="/home/${USER_NAME}/bin:${PATH}"

ENTRYPOINT [ "tmplpress" ]

CMD [ "-h" ]
