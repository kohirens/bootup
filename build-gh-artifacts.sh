#!/usr/bin/env bash

set -e

package_name="tmplpress"
platforms=("windows/amd64" "windows/386" "linux/amd64" "linux/386")

declare -a releaseFiles=()

for platform in "${platforms[@]}"; do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="${package_name}-${GOOS}-${GOARCH}"
    if [ "${GOOS}" = "windows" ]; then
        output_name="${output_name}"
    fi
    echo "building ${output_name}"
    env GOOS="${GOOS}" GOARCH="${GOARCH}" go build -o "${output_name}"
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    if [ "${GOOS}" = "windows" ]; then
        mv "${output_name}" "${output_name}.exe"
        zip "${output_name}.zip" "${output_name}.exe"
        releaseFiles+=("${output_name}.zip")
    else
        tar -zcvf "${output_name}.tar.gz" "${output_name}"
        releaseFiles+=("${output_name}.tar.gz")
    fi
done

if [ -z "${releaseFiles}" ]; then
    echo "error no release files found to upload"
    exit 1
fi

echo "releaseFiles=${releaseFiles}"

ls -ls .

# Switch to SSH to use the token stored in the environment variable GH_TOKEN.
gh config set git_protocol ssh --host github.com
gh auth status --hostname github.com

for releaseFile in "${releaseFiles[@]}"; do
    echo "uploading \"${releaseFile}\""
    gh release upload "${BUILD_VER}" "${releaseFile}"
done
