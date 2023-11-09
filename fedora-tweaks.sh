#!/usr/bin/env bash

# Enlève les flatpaks de Fedora, installe flathub, mets flatpak en premier dans Logiciels
flatpak remote-delete fedora
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
gsettings set org.gnome.software packaging-format-preference "['flatpak', 'rpm']"

# Mets le nombre de téléchargements parallèles de DNF à 10
echo -e "max_parallel_downloads=10" | sudo tee -a /etc/dnf/dnf.conf > /dev/null

# Installation de Firefox Gnome Theme
curl -s -o- https://raw.githubusercontent.com/rafaelmardojai/firefox-gnome-theme/master/scripts/install-by-curl.sh | bash

# Installation de MoreWaita
sudo dnf copr enable dusansimic/themes
sudo dnf install morewaita-icon-theme
gsettings set org.gnome.desktop.interface icon-theme 'MoreWaita'
