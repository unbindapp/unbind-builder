# unbind-builder

![GitHub License](https://img.shields.io/github/license/unbindapp/unbind-builder) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/unbindapp/unbind-builder)

Unbind build container image.

# Run a BuildKit instance as a container

docker run --rm --privileged -d --name buildkit moby/buildkit

# Set the buildkit host to the container

export BUILDKIT_HOST=docker-container://buildkit
