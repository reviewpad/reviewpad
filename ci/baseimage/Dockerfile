# Copyright (C) 2022 Explore.dev, Unipessoal Lda - All Rights Reserved
# Use of this source code is governed by a license that can be
# found in the LICENSE file.

# This dockerfile is used to create Reviewpad base image to run ci build on github actions.
# You can see the image being used at .github/workflows/build.yml under jobs > build > container > image.
# The image is publish on https://hub.docker.com/repository/docker/reviewpad/ci.

FROM golang:1.20

ENV LIBGIT2_ZIP v1.2.0.zip
ENV LIBGIT2 libgit2-1.2.0

WORKDIR /app

# Install necessary packages
RUN apt-get update && apt-get -y install unzip cmake libssl-dev && apt-get clean

# Install libgit2
RUN curl -OL https://github.com/libgit2/libgit2/archive/refs/tags/${LIBGIT2_ZIP} && \
    unzip -o $LIBGIT2_ZIP -d /tmp && \
    cd /tmp/${LIBGIT2} && \
    mkdir build && \
    cd build && \
    cmake .. && \
    cmake --build . --target install
