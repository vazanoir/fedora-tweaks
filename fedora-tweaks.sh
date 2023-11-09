#!/usr/bin/env bash

echo -e "Mets Flathub par défaut"
flatpak remote-delete fedora
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
gsettings set org.gnome.software packaging-format-preference "['flatpak', 'rpm']"

echo -e "Mets le nombre de téléchargements en parallèle de DNF à 10"
echo -e "max_parallel_downloads=10" | sudo tee -a /etc/dnf/dnf.conf > /dev/null

echo -e "Installation de Firefox Gnome Theme"
curl -s -o- https://raw.githubusercontent.com/rafaelmardojai/firefox-gnome-theme/master/scripts/install-by-curl.sh | bash

echo -e "Installation de MoreWaita"
dnf copr enable dusansimic/themes
dnf install morewaita-icon-theme
gsettings set org.gnome.desktop.interface icon-theme 'MoreWaita'
