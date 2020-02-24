#!/usr/bin/env bash
# Shell script to document MSI build process.

out() { echo -ne "[\e[32m+\e[0m] $1\n"; }
err() { echo -ne "[\e[31m!\e[0m] $1\n"; }

VERSION=0.7.1.0
ARCH=amd64
PROMU="${GOPATH}/bin/promu.exe"
PROJECT_ROOT="${GOPATH}/github.com/martinlindeh/wmi_exporter"
MSINAME="wmi_exporter-${VERSION}-${ARCH}.msi"
MSI_PATH="${PROJECT_ROOT}/installer/Output/${MSINAME}"

if ! [[ -e "${PROMU}" ]]; then
	err "Missing promu utility! (go install github.com/prometheus/promu)"
	exit 1
fi

if ! cd "$PROJECT_ROOT"; then
	err "Can't cd to ${PROJECT_ROOT}"
	exit 1
fi

out "Building WMI Exporter \e[33m${VERSION}\e[0m"
if ! ${PROMU} build --config "${PROJECT_ROOT}/.promu.yml"; then
	err "failed to build wmi_exporter.exe"
	exit 1
fi

out "Creating MSI installer.."
cd "${PROJECT_ROOT}/installer"
./build.ps1 -PathToExecutable "${PROJECT_ROOT}/wmi_exporter.exe" -Arch ${ARCH} -Version ${VERSION}

cp "${MSI_PATH}" "$(PWD)/${MSINAME}"
out "Created \e[35m${MSINAME}\e[0m"

