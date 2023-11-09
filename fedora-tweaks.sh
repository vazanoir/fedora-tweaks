#!/usr/bin/env bash

echo -e "Mise en place de Flathub par défaut"
flatpak remote-delete fedora
flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
gsettings set org.gnome.software packaging-format-preference "['flatpak', 'rpm']"

echo -e "Installation de MoreWaita"
dnf copr enable dusansimic/themes
dnf install morewaita-icon-theme
gsettings set org.gnome.desktop.interface icon-theme 'MoreWaita'
