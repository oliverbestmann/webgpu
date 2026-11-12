#!/usr/bin/env bash

set -ex -o pipefail

# gets the prebuilt release files from wgpu-native and installs
# the header files to wgpu/lib.
# the library files are pushed to branches according to the arch.

VERSION="v29.0.1.1"

ZIPS="releases-$VERSION"

SYSTEMS="android darwin ios linux windows"

# download files and extract zip files.
# the resulting files are in $ZIPS/*/
function download-and-extract {
  if ! [[ -f $ZIPS/ok ]] ; then
    mkdir -p $ZIPS

    pushd $ZIPS
      rm -rf wgpu-native zips
      mkdir zips

      git clone https://github.com/gfx-rs/wgpu-native wgpu-native
      pushd wgpu-native
        gh release download $VERSION --dir ../zips -p "*-release.zip"
      popd

      pushd zips

        for ZIP in *.zip ; do
          if [[ "$ZIP" = *-x86_64-msvc-* ]] ; then
            continue
          fi

          ARCH=$(cut -d- -f2-3 <<< "$ZIP")
          if [[ "$ZIP" = *-aarch64-simulator-* ]] ; then
            # the iOS simulator on Apple silicon, not the device (ios-aarch64)
            ARCH="$ARCH-simulator"
          fi
          mkdir "$ARCH"
          pushd "$ARCH"
            unzip ../"$ZIP"
          popd
        done
      popd
    popd

    touch $ZIPS/ok
  fi
}

# create the libs repository and fetch latest branches
function fetch-libs-repository {
  mkdir -p libs

  pushd libs
    if ! [ -d repo/.git ] ; then
      mkdir repo
      pushd repo
        git init .
        git remote add origin https://github.com/oliverbestmann/webgpu
      popd
    fi

    for SYSTEM in $SYSTEMS ; do
      git -C repo fetch origin libs-"$SYSTEM"
      # start from the remote branch, it may have commits from merged pull requests
      git -C repo switch -C libs-"$SYSTEM" origin/libs-"$SYSTEM"
    done
  popd
}

# push all library branches to the git repository
function push-libs-repository {
  for SYSTEM in $SYSTEMS ; do
    git -C libs/repo push origin libs-"$SYSTEM"
  done
}

# Since Go 1.27.4 and 1.26.11 (CVE-2026-97038), cgo refuses libraries with LLVM bitcode
# sections, .drectve directives it can not verify, or objects without sections.
# Rust puts all of these into the archive, none are needed to link it.
function strip-for-go() {
  local FILE="$1"
  local TARGET="$2"

  local ARGS=(
    --remove-section=.llvmbc --remove-section=.llvmcmd
    --remove-section=__LLVM,__bitcode --remove-section=__LLVM,__cmdline
  )

  if [[ "$FILE" = *.lib ]] ; then
    # llvm-objcopy can not rewrite MSVC libraries because of their import members, so
    # strip the objects one by one and replace them in place. Extracting everything would
    # lose import members, many of them share a name.
    local TMP
    TMP=$(mktemp -d)
    (
      cd "$TMP"
      llvm-ar x "$OLDPWD/$FILE" $(llvm-ar t "$OLDPWD/$FILE" | grep '\.o$')
      for OBJ in *.o ; do
        llvm-objcopy "${ARGS[@]}" "$OBJ"
      done
      llvm-ar r --format=coff "$OLDPWD/$FILE" *.o
    )
    rm -rf "$TMP"
    return
  fi

  if [[ "$TARGET" = windows ]] ; then
    # only holds -exclude-symbols, which matters when exporting from a DLL
    ARGS+=(--remove-section=.drectve)
  fi

  llvm-objcopy "${ARGS[@]}" "$FILE"

  # empty split dwarf objects
  local DWO
  DWO=$(llvm-ar t "$FILE" | grep '\.dwo$' || true)
  if [[ -n "$DWO" ]] ; then
    llvm-ar d "$FILE" $DWO
  fi
}

# copy files into target directories
function copy-to-target() {
  local ZIP="$1"
  local TARGET="$2"
  local ARCH="$3"
  local LIB="$4"

  # copy headers
  mkdir -p wgpu/lib/"$TARGET"/"$ARCH"
  cp $ZIPS/zips/"$ZIP"/include/webgpu/*.h wgpu/lib/"$TARGET"/"$ARCH"
  git add "wgpu/lib/$TARGET/$ARCH"


  # switch to the correct library branch
  git -C libs/repo switch "libs-$TARGET"

  # copy lib to branch
  mkdir -p libs/repo/libs-"$TARGET"/"$ARCH"
  cp $ZIPS/zips/"$ZIP"/lib/"$LIB" libs/repo/libs-"$TARGET"/"$ARCH"
  strip-for-go libs/repo/libs-"$TARGET"/"$ARCH"/"$LIB" "$TARGET"
  git -C libs/repo add libs-"$TARGET"/"$ARCH"

  # commit the file if something has changed
  if [ -n "$(git -C libs/repo status --porcelain)" ] ; then
    git -C libs/repo commit -m "update library to $VERSION"
  fi
}

function write-gitattributes() {
  cat > wgpu/lib/.gitattributes <<EOF
# See
# https://github.com/github/linguist/blob/249bbd1c2ffc631ca2ec628da26be5800eec3d48/docs/overrides.md#vendored-code

webgpu.h linguist-vendored
wgpu.h linguist-vendored
EOF
}

if ! command -v llvm-objcopy >/dev/null ; then
  echo "llvm-objcopy and llvm-ar are required, try: nix shell nixpkgs#llvm -c $0" >&2
  exit 1
fi

download-and-extract

fetch-libs-repository

copy-to-target "macos-aarch64"    "darwin"    "arm64"   libwgpu_native.a
copy-to-target "macos-x86_64"     "darwin"    "amd64"   libwgpu_native.a
copy-to-target "ios-aarch64"      "ios"       "arm64"   libwgpu_native.a
copy-to-target "ios-x86_64"       "ios"       "amd64"   libwgpu_native.a
copy-to-target "ios-aarch64-simulator" "ios" "arm64-simulator" libwgpu_native.a
copy-to-target "windows-aarch64"  "windows"   "arm64"   wgpu_native.lib
copy-to-target "windows-x86_64"   "windows"   "amd64"   libwgpu_native.a
copy-to-target "windows-i686"     "windows"   "386"     wgpu_native.lib
copy-to-target "linux-aarch64"    "linux"     "arm64"   libwgpu_native.a
copy-to-target "linux-x86_64"     "linux"     "amd64"   libwgpu_native.a
copy-to-target "android-aarch64"  "android"   "arm64"   libwgpu_native.a
copy-to-target "android-x86_64"   "android"   "amd64"   libwgpu_native.a
copy-to-target "android-armv7"    "android"   "arm"     libwgpu_native.a
copy-to-target "android-i686"     "android"   "386"     libwgpu_native.a

write-gitattributes

push-libs-repository

for SYSTEM in $SYSTEMS ; do
  env GOPROXY=direct go get -u github.com/oliverbestmann/webgpu/libs-$SYSTEM@libs-$SYSTEM
done

# re-generate enums & bindings
go generate ./wgpu

git add go.mod go.sum wgpu/lib/
git commit -m "chore: update wgpu to $VERSION"
