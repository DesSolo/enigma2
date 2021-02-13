VERSION=$(cat ../VERSION)
PROJECT_NAME=$(basename $(dirname "$PWD"))

echo "building $PROJECT_NAME:$VERSION"
podman build -t $PROJECT_NAME:$VERSION ../