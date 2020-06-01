package pkgmanlib

// PkgRefresh refresh/update of packages before update
var PkgRefresh = map[string]string{
	"debian":   "apt-get update",
	"centos":   "yum check-update",
	"fedora":   "dnf check-update",
	"opensuse": "zypper refresh",
	"arch":     "pacman -Sy",
	"rhel":     "yum update",
	"sles":     "zypper refresh",
	"ubuntu":   "apt-get update",
}

// PkgUpdate map of Pkgs with update flags for every dist. Updates take different flags and arguments
var PkgUpdate = map[string]string{
	"debian":   "apt-get upgrade -y",
	"centos":   "yum update -y",
	"fedora":   "dnf upgrade -y",
	"opensuse": "zypper update -y",
	"arch":     "pacman -Syu",
	"rhel":     "yum update -y",
	"sles":     "zypper update -y",
	"ubuntu":   "apt-get upgrade -y",
}

// PkgUpdateOS map of Pkgs with update flags for every dist. Updates take different flags and arguments
var PkgUpdateOS = map[string]string{
	"debian":   "apt-get dist-upgrade -y",
	"centos":   "yum update -y",
	"fedora":   "dnf system-upgrade download --refresh --releasever=$(awk -v s=1 '{print $3+s}' /etc/fedora-release) -y",
	"opensuse": "zypper dup -y",
	"arch":     "pacman -Syyu",
	"rhel":     "yum update -y",
	"sles":     "zypper dup -y",
	"ubuntu":   "apt-get dist-upgrade -y",
}

// PkgInstall map of Pkgs for every OS. Installation flags differ from dist to dist
var PkgInstall = map[string]string{
	"debian":   "apt-get install ",
	"centos":   "yum install ",
	"fedora":   "dnf install ",
	"opensuse": "zypper install ",
	"arch":     "pacman -S ",
	"rhel":     "yum install ",
	"sles":     "zypper install ",
	"ubuntu":   "apt-get install ",
}

// PkgUninstall map of Pkgs for every OS. Installation flags differ from dist to dist
var PkgUninstall = map[string]string{
	"debian":   "apt-get remove ",
	"centos":   "yum remove ",
	"fedora":   "dnf remove ",
	"opensuse": "zypper remove ",
	"arch":     "pacman -R ",
	"rhel":     "yum remove ",
	"sles":     "zypper remove ",
	"ubuntu":   "apt-get remove ",
}

// PkgSearch map to search installed packages
var PkgSearch = map[string]string{
	"debian":   "dpkg-query -l | grep ",
	"centos":   "rpm -qa | grep ",
	"fedora":   "rpm -qa | grep ",
	"opensuse": "rpm -qa | grep ",
	"arch":     "pacman -Q | grep ",
	"rhel":     "rpm -qa | grep ",
	"sles":     "rpm -qa | grep ",
	"ubuntu":   "dpkg-query -l | grep ",
}

// OmniTools map of
var OmniTools = map[string]string{
	"": "",
}
