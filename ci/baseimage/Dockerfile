# Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
# Use of this source code is governed by a license that can be
# found in the LICENSE file.

# This dockerfile is used to create Reviewpad base image to run ci build on github actions.
# You can see the image being used at .github/workflows/build.yml under jobs > build > container > image.
# The image is publish on https://hub.docker.com/repository/docker/reviewpad/ci.

FROM golang:1.21

ENV LIBGIT2_ZIP v1.2.0.zip
ENV LIBGIT2 libgit2-1.2.0

WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get -y install unzip cmake libssl-dev ca-certificates curl gnupg lsb-release && apt-get clean

# Install docker
RUN mkdir -m 0755 -p /etc/apt/keyrings

RUN curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg

RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/debian $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

RUN apt-get update && apt-get -y install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip -o awscliv2.zip && \
    ./aws/install

# Install libgit2
RUN curl -OL https://github.com/libgit2/libgit2/archive/refs/tags/${LIBGIT2_ZIP} && \
    unzip -o $LIBGIT2_ZIP -d /tmp && \
    cd /tmp/${LIBGIT2} && \
    mkdir build && \
    cd build && \
    cmake .. && \
    cmake --build . --target install
