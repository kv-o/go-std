// Package platform collects a list of platform names, codenames, and code
// characters.
//
// When targeting specific platforms, it is preferable to have a list of
// standard codenames associated with them. However, considering that programs
// typically target specific CPU architectures and operating systems, the
// resultant platform designation may still be slightly too verbose for
// frequent use. This is why, aside from providing codenames for each supported
// platform, a code character is provided as well.
package platform

import "runtime"

// Platform represents a platform, whether it be a CPU architecture or an OS.
type Platform struct {
	CodeChar rune
	CodeName string
	Name     string
}

// The following is a list of platform structures which provide code character,
// codename, and name associations for known platforms.
var (
	// CPU architectures
	Amd64    = Platform{'6', "amd64", "AMD64"}
	Arm      = Platform{'5', "arm", "little-endian ARM"}
	Arm64    = Platform{'7', "arm64", "little-endian ARM (64-bit)"}
	I386     = Platform{'8', "i386", "Intel 80386 and compatibles"}
	Loong64  = Platform{'l', "loong64", "Loongson (64-bit)"}
	M68k     = Platform{'m', "m68k", "Motorola 68000")
	Mips     = Platform{'0', "mips", "big-endian MIPS32"}
	Mips64   = Platform{'1', "mips64", "big-endian MIPS64"}
	Mips64le = Platform{'2', "mips64le", "MIPS64 (little-endian)"}
	Mipsle   = Platform{'3', "mipsle", "MIPS32 (little-endian)"}
	Ppc      = Platform{'p', "ppc", "Power PC"}
	Ppc64    = Platform{'4', "ppc64", "Power PC (64-bit)"}
	Ppc64le  = Platform{'9', "ppc64le", "Power PC (64-bit, little-endian)"}
	Riscv64  = Platform{'r', "riscv64", "RISC-V (64-bit)"}
	S390x    = Platform{'z', "s390x", "IBM z/Architecture"}
	Sparc64  = Platform{'s', "sparc64", "SPARC V9"}
	// Operating systems
	Aix       = Platform{'x', "aix", "IBM AIX"}
	Android   = Platform{'a', "android", "Android"}
	Bare      = Platform{'b', "bare", "Bare metal"}
	Darwin    = Platform{'d', "darwin", "Darwin and derivatives"}
	Dragonfly = Platform{'y', "dragonfly", "DragonFly BSD"}
	Freebsd   = Platform{'f', "freebsd", "FreeBSD"}
	Illumos   = Platform{'m', "illumos", "Illumos"}
	Ios       = Platform{'i', "ios", "iOS"}
	Linux     = Platform{'l', "linux", "Linux"}
	Netbsd    = Platform{'n', "netbsd", "NetBSD"}
	Openbsd   = Platform{'o', "openbsd", "OpenBSD"}
	Plan9     = Platform{'p', "plan9", "Plan 9 and derivatives"}
	Solaris   = Platform{'s', "solaris", "Oracle Solaris"}
	Windows   = Platform{'w', "windows", "Windows NT"}
)

// Arch is a slice of all known CPU architectures.
var Arch = []Platform{
	Amd64,
	Arm,
	Arm64,
	I386,
	Loong64,
	M68k,
	Mips,
	Mips64,
	Mips64le,
	Mipsle,
	Ppc,
	Ppc64,
	Ppc64le,
	Riscv64,
	S390x,
	Sparc64,
}

// OS is a slice of all known operating systems.
var OS = []Platform{
	Aix,
	Android,
	Bare,
	Darwin,
	Dragonfly,
	Freebsd,
	Illumos,
	Ios,
	Linux,
	Netbsd,
	Openbsd,
	Plan9,
	Solaris,
	Windows,
}

// Return the current CPU architecture and OS.
func Current() (arch, os Platform) {
	switch runtime.GOARCH {
	case "386":
		return WithCodeName(Arch, "i386"), WithCodeName(OS, runtime.GOOS)
	default:
		return WithCodeName(Arch, runtime.GOARCH), WithCodeName(OS, runtime.GOOS)
	}
}

// WithCodeChar returns the first platform in p with the given code character r.
// If there are no matching platforms, returns an empty Platform.
func WithCodeChar(p []Platform, r rune) Platform {
	for _, plat := range p {
		if plat.CodeChar == r {
			return plat
		}
	}
	return Platform{}
}

// WithCodeName returns the first platform in p with the given code name r.
// If there are no matching platforms, returns an empty Platform.
func WithCodeName(p []Platform, s string) Platform {
	for _, plat := range p {
		if plat.CodeName == s {
			return plat
		}
	}
	return Platform{}
}

// WithName returns the first platform in p with the given platform name s.
// If there are no matching platforms, returns an empty Platform.
func WithName(p []Platform, s string) Platform {
	for _, plat := range p {
		if plat.Name == s {
			return plat
		}
	}
	return Platform{}
}
