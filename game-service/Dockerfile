# Build binary.
FROM 713614461671.dkr.ecr.eu-central-1.amazonaws.com/golang:1.21.3-ubuntu22.04 AS build-env
RUN apt update && \
    apt upgrade -y && \
    apt install -y git ssh build-essential

ARG BUILD_NR
ARG BUILD_DATE
ARG GIT_HASH

ENV BUILD_NR=$BUILD_NR
ENV BUILD_DATE=$BUILD_DATE
ENV GIT_HASH=$GIT_HASH

RUN echo "BUILD_NR=$BUILD_NR" && \
    echo "BUILD_DATE=$BUILD_DATE" && \
    echo "GIT_HASH=$GIT_HASH"

# Create appuser.
ENV USER=appuser
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN
RUN adduser --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" --no-create-home --uid "${UID}" "${USER}"

# Copy codecommit access key from pipelines base image.
RUN mkdir /root/.ssh/
ADD codecommit_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa

# Skip Host verification for git and add config for IAM role
RUN echo "StrictHostKeyChecking no " > /root/.ssh/config && \
    echo "Host git-codecommit.eu-central-1.amazonaws.com" >> /root/.ssh/config && \
    echo "  User APKA2MJV4RLT2F34QSPV" >> /root/.ssh/config && \
    echo "  Port 22" >> /root/.ssh/config && \
    git config --global url."ssh://git-codecommit.eu-central-1.amazonaws.com:".insteadOf "https://git-codecommit.eu-central-1.amazonaws.com/"

COPY . /go/src/git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service
WORKDIR /go/src/git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service

RUN mkdir -p /go/src/git-codecommit.eu-central-1.amazonaws.com/v1/repos && \
    cd /go/src/git-codecommit.eu-central-1.amazonaws.com/v1/repos && \
    git clone ssh://git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng-binary && \
    cp ./prng-binary/libprng.so /usr/local/lib && \
    cp ./prng-binary/libprng.h /usr/local/include && \
    chmod +x /usr/local/lib/libprng.so && \
    sha1sum /usr/local/lib/libprng.so

RUN export GOPRIVATE="git-codecommit.eu-central-1.amazonaws.com" && \
    export TARGET_HASH="git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes.GameService" && \
    export TARGET_BUILD="git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts.SemBuild" && \
    export TARGET_DATE="git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts.SemDate" && \
    go mod tidy && \
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -tags DEBUG -ldflags="-s -w -X ${TARGET_HASH}=${GIT_HASH} -X ${TARGET_BUILD}=${BUILD_NR} -X ${TARGET_DATE}=${BUILD_DATE}" -o dist/game-service cmd/*.go

# Build image.
FROM 713614461671.dkr.ecr.eu-central-1.amazonaws.com/ubuntu:22.04
RUN apt update && \
    apt upgrade -y && \
    apt install -y ca-certificates net-tools

COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /etc/group /etc/group
COPY --from=build-env /usr/local/lib/libprng.so /usr/local/lib/
COPY --from=build-env /go/src/git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/dist/game-service /app/

EXPOSE 8080

USER appuser:appuser

WORKDIR /app
ENTRYPOINT /app/game-service
