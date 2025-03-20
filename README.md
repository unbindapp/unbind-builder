# unbind-builder

One off jobs for building code bases. ...

# Run a BuildKit instance as a container

docker run --rm --privileged -d --name buildkit moby/buildkit

# Set the buildkit host to the container

export BUILDKIT_HOST=docker-container://buildkit
