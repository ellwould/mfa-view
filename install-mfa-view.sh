#!/bin/bash

# Install Script for MFA View

#----------------------------------------------------------------------

# Check user is root otherwise exit script

if [ "$EUID" -ne 0 ]
then
  printf "\nPlease run as root\n\n";
  exit;
fi;

cd /root;

#----------------------------------------------------------------------

# Check MFA View has been cloned from GitHub

if [ ! -d "/root/mfa-view" ]
then
  printf "\nDirectory mfa-view does not exist in /root.\n";
  printf "Please run commands: \"cd /root; git clone https://github.com/ellwould/mfa-view\"\n";
  printf "and run install script again\n\n";
  exit;
fi;

#----------------------------------------------------------------------

# Copy unit file and reload systemd deamon

cp /root/mfa-view/systemd/mfaview.service /usr/lib/systemd/system/;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Remove any previous version of Go, download and install Go 1.24.4

wget -P /root https://go.dev/dl/go1.24.4.linux-amd64.tar.gz;
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.24.4.linux-amd64.tar.gz;

#----------------------------------------------------------------------

# Create directores used for HTML/CSS, configuration and CSV key file

mkdir -p /etc/mfaview{html-css,key};

# Copy HTML/CSS start and end files

cp /root/mfaview/html-css/* /etc/mfaview/html-css/;

# Copy MFA View configuration file

cp /root/mfa-view/env/mfaview.env /etc/mfaview/mfaview.env

# Copy key file

cp /root/mfa-view/key/mfaview-key.csv /etc/mfaview-key.csv

# Create Go directories in root home directory for compiling the source code

mkdir -p /root/go/{bin,pkg,src/mfaview};

# Copy MFA View source code

cp /root/mfa-view/go/mfaview.go /root/go/src/mfaview/mfaview.go;

# Create Go mod for mfaview

export PATH=$PATH:/usr/local/go/bin;
cd /root/go/src/mfaview;
go mod init root/go/src/mfaview;
go mod tidy;

# Compile mfaview.go

cd /root/go/src/mfaview;
go build mfaview.go;
cd /root;

# Create system user named mfaview with no shell, no home directory and lock the account

useradd -r -s /bin/false mfaview;
usermod -L mfaview;

# Change executables file permissions, owner, group and move executables

chown root:mfaview /root/go/src/mfaview/mfaview;
chmod 050 /root/go/src/mfaview/mfaview;
mv /root/go/src/mfaview/mfaview /usr/bin/mfaview;

# Change mfaview file permissions, owner and group

chown -R root:mfaview /etc/mfaview;
chmod 050 /etc/mfaview;
chmod 050 /etc/mfaview/*;
chmod 060 /etc/mfaview/mfaview.env;
chmod 060 /etc/mfaview/key/mfaview-key.csv;

# Enable mfaview on boot

systemctl enable mfaview;
