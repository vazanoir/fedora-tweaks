#!/usr/bin/env fish

### LOOKS -----------------------------------------------------------------

# Install Betterfox and Firefox Gnome Theme
curl -s -o- https://raw.githubusercontent.com/vazanoir/update-betterfox/main/install-by-curl.sh | bash
curl -s -o- https://raw.githubusercontent.com/rafaelmardojai/firefox-gnome-theme/master/scripts/install-by-curl.sh | bash

# Install MoreWaita
if not dnf list installed | grep -q morewaita
	sudo dnf copr enable dusansimic/themes -y
	sudo dnf install morewaita-icon-theme -y
	gsettings set org.gnome.desktop.interface icon-theme 'MoreWaita'
end



### TWEAKS ----------------------------------------------------------------

set content "max_parallel_downloads=10"
set file "/etc/dnf/dnf.conf"
if not grep -q $content $file
	echo -e $content | sudo tee -a $file > /dev/null
end

# Remove Fedora's flatpaks, install flathub for system, user and beta branch
if flatpak remotes | grep -q fedora
	flatpak remote-delete fedora
end

flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists flathub-beta https://flathub.org/beta-repo/flathub-beta.flatpakrepo

# Set flatpak first in Gnome's Software
gsettings set org.gnome.software packaging-format-preference "['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']"



### BUGFIXES --------------------------------------------------------

# Fix issue between SELinux and Source games
if getsebool allow_execheap | grep -q off
	sudo setsebool -P allow_execheap 1
end

# Fix issue with big games
set content "vm.max_map_count = 16777216"
set file "/etc/sysctl.conf"
if not grep -q $content $file
	echo -e $content | sudo tee -a $file > /dev/null
	sudo sysctl -p
end
