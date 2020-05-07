package tags

import (
	"fmt"
	"strings"

	"github.com/coreos/go-semver/semver"
)

// DefaultTagSuffix returns a set of default suggested tags
// based on the commit ref with an attached suffix.
func VersionTagsTagSuffix(ref, suffix string) []string {
	tags := VersionTags(ref)
	if len(suffix) == 0 {
		return tags
	}
	for i, tag := range tags {
		if tag == "latest" {
			tags[i] = suffix
		} else {
			tags[i] = fmt.Sprintf("%s-%s", tag, suffix)
		}
	}
	return tags
}

func splitOff(input string, delim string) string {
	parts := strings.SplitN(input, delim, 2)

	if len(parts) == 2 {
		return parts[0]
	}

	return input
}

// VersionTags returns a set of default suggested tags based on
// the version.
func VersionTags(rawVersion string) []string {
	v := stripVersionPrefix(rawVersion)

	version, err := semver.NewVersion(v)
	if err != nil {
		return []string{}
	}
	if version.PreRelease != "" || version.Metadata != "" {
		return []string{
			version.String(),
		}
	}

	v = splitOff(splitOff(v, "+"), "-")
	dotParts := strings.SplitN(v, ".", 3)

	if version.Major == 0 {
		return []string{
			fmt.Sprintf("%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor),
			fmt.Sprintf("%0*d.%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor, len(dotParts[2]), version.Patch),
		}
	}
	return []string{
		fmt.Sprintf("%0*d", len(dotParts[0]), version.Major),
		fmt.Sprintf("%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor),
		fmt.Sprintf("%0*d.%0*d.%0*d", len(dotParts[0]), version.Major, len(dotParts[1]), version.Minor, len(dotParts[2]), version.Patch),
	}
}

func stripVersionPrefix(ref string) string {
	ref = strings.TrimPrefix(ref, "v")
	return ref
}
