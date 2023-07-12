#!/usr/bin/env sh
set -e

# check_sha is separated into a defined function so that we can
# capture the exit code effectively with `set -e` enabled
check_sha() {
  (
    cd /tmp/
    shasum -sc --ignore-missing "$1"
  )

  return $?
}

# check if jq is installed
if ! [ -x "$(command -v jq)" ]; then
  echo "Error: jq is not installed." >&2
  echo "Please install and try again. https://stedolan.github.io/jq/download"
  exit 1
fi

os=$(uname | tr '[:upper:]' '[:lower:]')
arch=$(uname -m | tr '[:upper:]' '[:lower:]' | sed -e s/x86_64/amd64/)
if [ "$arch" = "aarch64" ]; then
  arch="arm64"
fi

project_name="kse"

version=$(curl -s https://api.github.com/repos/leonardobiffi/$project_name/releases/latest | jq .name -r)
url="https://github.com/leonardobiffi/$project_name/releases/download/$version"
checksum_url="$url/checksums.txt"
binary_version=$(echo $version | cut -c 2-)
binary_name=$project_name
zip="${binary_name}_${binary_version}_${os}_${arch}.zip"
release_name="${binary_name}_${os}_${arch}"
echo "Downloading latest release of $release_name..."
curl -sL "$url/$zip" -o "/tmp/$zip"
echo

code=$(curl -s -L -o /dev/null -w "%{http_code}" "$checksum_url")
if [ "$code" = "404" ]; then
    echo "Skipping checksum validation as the sha for the release could not be found, no action needed."
else
  if [ -x "$(command -v shasum)" ]; then
    echo "Validating checksum for $release_name..."
    curl -sL "$checksum_url" -o "/tmp/$release_name.sha256"

    if ! check_sha "$release_name.sha256"; then
      echo
      read -r -p "Installation checksum failed. This could be a security issue. Would you like to continue? (y/n) " answer
      if [ "$answer" != "y" ]; then
        echo
        echo "Exiting ..."
        exit 1
      fi
    fi

    rm "/tmp/$release_name.sha256"
  else
    echo "Skipping checksum validation as the shasum command could not be found, no action needed."
  fi
fi
echo

unzip -q -o "/tmp/$zip" -d /tmp
rm "/tmp/$zip"

echo "Moving /tmp/$binary_name to /usr/local/bin/$binary_name (you might be asked for your password due to sudo)"
if [ -x "$(command -v sudo)" ]; then
  sudo mv "/tmp/$binary_name" "/usr/local/bin/$binary_name"
else
  mv "/tmp/$binary_name" "/usr/local/bin/$binary_name"
fi
echo
echo "Completed installing $binary_name $version"
