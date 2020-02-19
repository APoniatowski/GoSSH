package pkgmanlib

// PackageManagerUpdate map of packagemanagers with update flags for every dist. Updates take different flags and arguments
var PackageManagerUpdate = map[string]string{
	"debian":   "apt-get",
	"centos":   "yum",
	"fedora":   "dnf",
	"opensuse": "zypper",
	"arch":     "pacman",
}

// PackageManagerUpdateOS map of packagemanagers with update flags for every dist. Updates take different flags and arguments
var PackageManagerUpdateOS = map[string]string{
	"debian":   "apt-get",
	"centos":   "yum",
	"fedora":   "dnf",
	"opensuse": "zypper",
	"arch":     "pacman",
}

// PackageManagerInstall map of packagemanagers for every OS. Installation flags differ from dist
var PackageManagerInstall = map[string]string{
	"debian":   "apt-get",
	"centos":   "yum",
	"fedora":   "dnf",
	"opensuse": "zypper",
	"arch":     "pacman",
}
