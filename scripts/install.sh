#!/bin/sh

npm i -g \
    semantic-release \
    @semantic-release/changelog \
    @semantic-release/commit-analyzer \
    @semantic-release/exec \
    @semantic-release/git \
    @semantic-release/npm \
    @semantic-release/release-notes-generator

apk update && apk add git
