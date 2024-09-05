#!/usr/bin/env fish

### APPARENCES -----------------------------------------------------------------

# Installation de Betterfox et de Firefox Gnome Theme
curl -s -o- https://raw.githubusercontent.com/vazanoir/update-betterfox/main/install-by-curl.sh | bash
curl -s -o- https://raw.githubusercontent.com/rafaelmardojai/firefox-gnome-theme/master/scripts/install-by-curl.sh | bash

# Installation de MoreWaita
if not dnf list installed | grep -q morewaita
	sudo dnf copr enable dusansimic/themes -y
	sudo dnf install morewaita-icon-theme -y
	gsettings set org.gnome.desktop.interface icon-theme 'MoreWaita'
end



### AJUSTEMENTS ----------------------------------------------------------------

# Mets le nombre de téléchargements parallèles de DNF à 10
set content "max_parallel_downloads=10"
set file "/etc/dnf/dnf.conf"
if not grep -q $content $file
	echo -e $content | sudo tee -a $file > /dev/null
end

# Enlève les flatpaks de Fedora, installe flathub, mets flatpak en premier dans Logiciels
if flatpak remotes | grep -q fedora
	flatpak remote-delete fedora
end

flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists flathub-beta https://flathub.org/beta-repo/flathub-beta.flatpakrepo
gsettings set org.gnome.software packaging-format-preference "['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']"



### CORRECTIONS DE BUGS --------------------------------------------------------

# Règle le problème entre SELinux et les jeux Source
if getsebool allow_execheap | grep -q off
	sudo setsebool -P allow_execheap 1
end

# Règle un problème pour certains gros jeux
set content "vm.max_map_count = 16777216"
set file "/etc/sysctl.conf"
if not grep -q $content $file
	echo -e $content | sudo tee -a $file > /dev/null
	sudo sysctl -p
end
