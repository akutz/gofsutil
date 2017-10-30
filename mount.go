package mount

import (
	"context"
	"errors"
)

// Most of this file is based on k8s.io/pkg/util/mount

// ErrNotImplemented is returned when a platform does not implement
// the contextual function.
var ErrNotImplemented = errors.New("not implemented")

// Info is information about a single mount point.
type Info struct {
	// Device is the device on which the filesystem is mounted.
	Device string

	// Path is the filesystem path to which the device is mounted.
	Path string

	// Source may be set to one of two values:
	//
	//   1. If this is a bind mount created with "bindfs" then Source
	//      is set to the filesystem path bind mounted to Path.
	//
	//   2. If this is any other type of mount then Source is set to
	//      a concatenation of the mount source and the root of
	//      the mount within the file system (fields 10 & 4 from
	//      the section on /proc/<pid>/mountinfo at
	//      https://www.kernel.org/doc/Documentation/filesystems/proc.txt).
	//
	// It is not possible to diffentiate a native bind mount from a
	// non-bind mount after the native bind mount has been created.
	// Therefore, while the Source field will be set to the filesystem
	// path bind mounted to Path for native bind mounts, the value of
	// the Source field can in no way be used to determine *if* a mount
	// is a bind mount.
	Source string

	// Type is the filesystem type.
	Type string

	// Opts are the mount options used to create this mount point.
	Opts []string
}

// Entry is a superset of Info and maps to the fields of a mount table
// entry:
//
//   (1) mount ID:  unique identifier of the mount (may be reused after umount)
//   (2) parent ID:  ID of parent (or of self for the top of the mount tree)
//   (3) major:minor:  value of st_dev for files on filesystem
//   (4) root:  root of the mount within the filesystem
//   (5) mount point:  mount point relative to the process's root
//   (6) mount options:  per mount options
//   (7) optional fields:  zero or more fields of the form "tag[:value]"
//   (8) separator:  marks the end of the optional fields
//   (9) filesystem type:  name of filesystem of the form "type[.subtype]"
//   (10) mount source:  filesystem specific information or "none"
//   (11) super options:  per super block options
type Entry struct {
	// Root of the mount within the filesystem.
	Root string

	// MountPoint relative to the process's root
	MountPoint string

	// MountOpts are per-mount options.
	MountOpts []string

	// FSType is the name of filesystem of the form "type[.subtype]".
	FSType string

	// MountSource is filesystem specific information or "none"
	MountSource string
}

// EntryScanFunc defines the signature of the function that is optionally
// provided to the functions in this package that scan the mount table.
// The mount entry table is ignored when this function returns a false
// value or error.
type EntryScanFunc func(
	ctx context.Context,
	entry Entry,
	cache map[string]Entry) (Info, bool, error)

// GetDiskFormat uses 'lsblk' to see if the given disk is unformatted.
func GetDiskFormat(ctx context.Context, disk string) (string, error) {

	return getDiskFormat(ctx, disk, nil)
}

// GetDiskFormatWithEntryScanFunc is GetDiskFormat with a custom EntryScanFunc.
func GetDiskFormatWithEntryScanFunc(
	ctx context.Context,
	disk string,
	scanEntry EntryScanFunc) (string, error) {

	return getDiskFormat(ctx, disk, scanEntry)
}

// FormatAndMount uses unix utils to format and mount the given disk.
func FormatAndMount(
	ctx context.Context,
	source, target, fsType string,
	options ...string) error {

	return formatAndMount(ctx, source, target, fsType, options)
}

// Mount mounts source to target as fstype with given options.
//
// The parameters 'source' and 'fstype' must be empty strings in case they
// are not required, e.g. for remount, or for an auto filesystem type where
// the kernel handles fstype automatically.
//
// The 'options' parameter is a list of options. Please see mount(8) for
// more information. If no options are required then please invoke Mount
// with an empty or nil argument.
func Mount(source, target, fsType string, options ...string) error {
	return mount(source, target, fsType, options)
}

// BindMount behaves like Mount was called with a "bind" flag set
// in the options list.
func BindMount(source, target string, options ...string) error {
	if options == nil {
		options = []string{"bind"}
	} else {
		options = append(options, "bind")
	}
	return mount(source, target, "", options)
}

// Unmount unmounts the target.
func Unmount(target string) error {
	return unmount(target)
}

// GetMounts returns a slice of all the mounted filesystems.
//
// * Linux hosts use mount_namespaces to obtain mount information.
//
//   Support for mount_namespaces was introduced to the Linux kernel
//   in 2.2.26 (http://man7.org/linux/man-pages/man5/proc.5.html) on
//   2004/02/04.
//
//   The kernel documents the contents of "/proc/<pid>/mountinfo" at
//   https://www.kernel.org/doc/Documentation/filesystems/proc.txt.
//
// * Darwin hosts parse the output of the "mount" command to obtain
//   mount information.
func GetMounts(ctx context.Context) ([]Info, error) {

	return getMounts(ctx, nil)
}

// GetMountsWithEntryScanFunc is GetMounts with a custom EntryScanFunc.
func GetMountsWithEntryScanFunc(
	ctx context.Context,
	scanEntry EntryScanFunc) ([]Info, error) {

	return getMounts(ctx, scanEntry)
}

// GetDevMounts returns a slice of all mounts for the provided device.
func GetDevMounts(ctx context.Context, dev string) ([]Info, error) {

	return getDevMounts(ctx, dev, nil)
}

// GetDevMountsWithEntryScanFunc is GetDevMounts with a custom EntryScanFunc.
func GetDevMountsWithEntryScanFunc(
	ctx context.Context,
	dev string,
	scanEntry EntryScanFunc) ([]Info, error) {

	return getDevMounts(ctx, dev, scanEntry)
}
