


# $GOOS	$GOARCH
# android	arm
# darwin	386
# darwin	amd64
# darwin	arm
# darwin	arm64
# dragonfly	amd64
# freebsd	386
# freebsd	amd64
# freebsd	arm
# linux	386
# linux	amd64
# linux	arm
# linux	arm64
# linux	ppc64
# linux	ppc64le
# linux	mips64
# linux	mips64le
# netbsd	386
# netbsd	amd64
# netbsd	arm
# openbsd	386
# openbsd	amd64
# openbsd	arm
# plan9	386
# plan9	amd64
# solaris	amd64
# windows	386
# windows	amd64

# env GOOS=linux GOARCH=arm go build -v .
# env GOOS=darwin GOARCH=amd64 go build -v .

#  https://github.com/aktau/github-release
#
#  expects GITHUB_TOKEN in env
#
github-release info -u lytics -r lytics

# TAG=$(git describe $(git rev-list --tags --max-count=1))

TAG="latest"
echo "releasing $TAG"

# if we are re-running, lets delete it first
github-release delete --user lytics --repo lytics --tag $TAG

# create a formal release
github-release release \
    --user lytics \
    --repo lytics \
    --tag $TAG \
    --name "Lytics CLI Latest" \
    --description "
Scripts to download and save the binary and rename to lytics

\`\`\`
# linux/amd64
curl -Lo lytics https://github.com/lytics/lytics/releases/download/latest/lytics_linux && chmod +x lytics && sudo mv lytics /usr/local/bin/

# OS X/amd64 
curl -Lo lytics https://github.com/lytics/lytics/releases/download/lytics/lytics_mac && chmod +x lytics && sudo mv lytics /usr/local/bin/


\`\`\`
"

# upload a file, the mac osx amd64 binary
echo "Creating and uploading mac client"
env GOOS=darwin GOARCH=amd64 go build
github-release upload \
    --user lytics \
    --repo lytics \
    --tag $TAG \
    --name "lytics_mac" \
    --file lytics

# now linux
go build
echo "Creating and uploading linux client"
github-release upload \
    --user lytics \
    --repo lytics \
    --tag $TAG \
    --name "lytics_linux" \
    --file lytics