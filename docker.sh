#!use/bin/env sh

# docker deamon commands reference: https://docs.docker.com/get-started/part2/#recap-and-cheat-sheet-optional

# Create image using this directory's Dockerfile
# This will fetch the golang base image from Docker Hub, copy the package source to it, build the package inside it, and tag the resulting image as outyet.
docker build -t outyet .

# To run a container from the resulting image:
# The --publish flag or -p tells docker to publish the container's port 8080 on the external port 6060.
# The --detach flag or -d runs in detached mode
# The --name flag gives our container a predictable name to make it easier to work with.
# The --rm flag tells docker to remove the container image when the outyet server exits.
docker run --publish 6060:8080 --detach --name test --rm outyet

# Gracefully stop the specified container.
docker stop test
# or
docker container stop <hash>
