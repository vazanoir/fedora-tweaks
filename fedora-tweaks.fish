#!/usr/bin/env fish

### TWEAKS ----------------------------------------------------------------

# Set max parallel downloads to a bigger number
set content "max_parallel_downloads=20"
set file "/etc/dnf/dnf.conf"
if not grep -q $content $file
	echo -e $content | sudo tee -a $file > /dev/null
end

# Remove Fedora's flatpaks, install flathub and beta for system and user
if flatpak remotes | grep -q fedora
	flatpak remote-delete fedora
end

flatpak remote-add --if-not-exists flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists --user flathub https://dl.flathub.org/repo/flathub.flatpakrepo
flatpak remote-add --if-not-exists flathub-beta https://flathub.org/beta-repo/flathub-beta.flatpakrepo
flatpak remote-add --if-not-exists --user flathub-beta https://flathub.org/beta-repo/flathub-beta.flatpakrepo

# Set flatpak first in Gnome's Software
gsettings set org.gnome.software packaging-format-preference "['flatpak:flathub', 'flatpak:flathub-beta', 'rpm']"

# Load i2c-dev and i2c-piix4 kernel modules for hardware detection in software like OpenRGB
set file "/etc/modules-load.d/i2c.conf"
if not test -e $file
    echo -e "i2c-dev\ni2c-piix4" | sudo tee -a $file
end

# Install steam-devices for controller support in Steam's flatpak
if dnf list installed | grep -q steam-devices
    sudo dnf install steam-devices
end

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
