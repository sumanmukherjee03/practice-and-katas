#!/usr/bin/env bash

export TMPDIR
TMPDIR="$(mktemp -d /tmp/provision.XXXX)"
trap 'echo "Cleaning up"; rm -rf "$TMPDIR"; echo "Removed temp dir";' EXIT

export DEBIAN_FRONTEND="noninteractive"
export MAIN_PID=$$

set -e

die() {
  echo "Error : $*" >/dev/stderr
  kill -s TERM $MAIN_PID
  exit 1
}
export -f die

info() {
  echo -n "Info : $*"
}
export -f info

ok() {
  echo -n ''
}
export -f ok

disable_automatic_updates() {
  # We need to run disable_automatic_updates at the very first because by defaults auto updates are on in ubuntu
  # So, as soon as the AMI builder starts the instance, ubuntu starts upgrading in the background and creates a dpkg lock.
  # The script continues and eventually times out waiting to acquire the dpkg lock when running apt-get update or install.
  local apt_conf=/etc/apt/apt.conf.d/20auto-upgrades
  systemctl stop apt-daily.service
  systemctl stop apt-daily.timer
  systemctl stop apt-daily-upgrade.service
  systemctl stop apt-daily-upgrade.timer
  systemctl disable apt-daily.service
  systemctl disable apt-daily.timer
  systemctl disable apt-daily-upgrade.service
  systemctl disable apt-daily-upgrade.timer
  sleep 60
  systemctl list-unit-files apt* --all
  sed -i -e 's#Unattended-Upgrade "1"#Unattended-Upgrade "0"#g;s#Update-Package-Lists "1"#Update-Package-Lists "0"#g' "$apt_conf"
  echo 'APT::Periodic::Enable "0";' >> /etc/apt/apt.conf.d/10periodic
  ok
}

log_debug_info() {
  # `lsb_release -a` is also an available option for getting OS info, but sometimes that binary might break due to python symlink problems
  cat /etc/*-release
  uname -a
  echo "User id : $(id -u)"
  echo "Group id : $(id -g)"
  echo "User name : $(id -un)"
  echo "Group name : $(id -gn)"
  printenv
  ok
}

validate_running_as_root() {
  [[ $(id -u) -eq 0 ]] \
    || die "This script needs to be run as user root"
  ok
}

install_essential_deps() {
  apt-get -y --no-install-recommends install \
    software-properties-common \
    apt-transport-https \
    ca-certificates \
    gnupg-agent \
    make \
    curl \
    unzip \
    build-essential \
    libssl-dev \
    libffi-dev \
    libz-dev \
    libcurl4-gnutls-dev \
    libexpat1-dev \
    gettext \
    cmake \
    gcc \
    zlib1g-dev \
    libncurses5-dev \
    libgdbm-dev \
    libnss3-dev \
    libsqlite3-dev \
    libreadline-dev \
    wget \
    libbz2-dev \
    gperf \
    git
  ok
}

install_python3_and_essentials() {
  # We must update the package list so that `python3.8` package is available for installation
  apt-get -y update
  apt-get -y --no-install-recommends install python3.8
  apt-get -y --no-install-recommends install python3-dev python3-pip
  # Set python as python3 so that when other package install, they install referring to the correct version of python we want ie python3
  # Changing system python is not ideal but such is the nature of python. 2 of them cant coexist peacefully and breaks packages
  update-alternatives --install /usr/bin/python python "$(command -v python3.8)" 10
  update-alternatives --set python "$(command -v python3.8)"
  # pushd .
  # cd /usr/lib/python3/dist-packages
  # ln -s apt_pkg.cpython-{36m,37m}-x86_64-linux-gnu.so
  # popd
  # Dont let pip install it's cache in users' HOME dir
  python3 -m pip --no-cache-dir install setuptools wheel pyyaml
  # Package libgnome-keyring-dev is not available in ubuntu focal 20.04. This is the replacement for python's keyring
  apt-get -y --no-install-recommends install python3-dbusmock
  ok
}

install_nodejs_and_essentials() {
  # Refer to official documentation here for installing LTS node version - https://github.com/nodesource/distributions/blob/master/README.md
  curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -
  apt-get install -y nodejs
  npm install -g npm@8.5.3
  ok
}

create_non_root_user() {
  useradd -m builder
  usermod -aG sudo builder
  echo '%sudo ALL=(ALL) NOPASSWD: ALL' >> /etc/sudoers
  usermod --shell /bin/bash builder
  mkdir /home/builder/build
}

main() {
  pushd .
  cd "$TMPDIR"
  disable_automatic_updates
  log_debug_info
  validate_running_as_root
  apt-get clean
  install_python3_and_essentials
  install_essential_deps
  install_nodejs_and_essentials
  create_non_root_user
  popd
  ok
}

[[ "${BASH_SOURCE[0]}" = "$0" ]] && main "$@"
