#!/bin/bash

# Uninstall script for MFA View

#----------------------------------------------------------------------

# Check user is root otherwise exit script

if [ "$EUID" -ne 0 ]
then
  printf "\nPlease run as root\n\n";
  exit;
fi;

cd /root;

#----------------------------------------------------------------------

# Stop MFA View automatically starting on boot

systemctl stop mfaview.service;
systemctl disable mfaview.service;

# Remove MFA View unit file and reload systemd deamon

rm /usr/lib/systemd/system/mfaview.service;
systemctl daemon-reload;

#----------------------------------------------------------------------

# Move mfaview-key.csv file to /root and change the owner and group to root

mv /etc/mfaview/key/mfaview-key.csv /root/mfaview-key.csv;
chown root:root /root/mfaview-key.csv;

# Remove MFA View binary

rm /usr/bin/mfaview;

# Remove all other directores and files used by MFA View

rm -r /etc/mfaview;

# Remove MFA View source code in root home directory

rm -r /root/go/src/mfaview;

# Remove the user and group mfaview from the system

userdel mfaview;
