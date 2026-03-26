package locations

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/0ghny/gitconfig/internal/adapter/filesystem"
	"github.com/0ghny/gitconfig/internal/domain/gitconfig"
	"github.com/0ghny/gitconfig/internal/domain/location"
	"github.com/0ghny/gitconfig/internal/templates"
	"github.com/spf13/afero"
)

const (
	regexStrGitConfigLocationGroups string = `(?m)^#[\s]*gitconfig.location.key[\s]*(?P<key>.*)[\s]*\[includeIf[\s]* \"gitdir:(?P<dir>.*)\"\][\s]*[\s]*path[\s]*=[\s]*(?P<path>.*)$`
)

// LocationManager orchestrates operations on the .gitconfig file (via GitConfig)
// and on the individual location config files.
type LocationManager struct {
	gitcfg *gitconfig.GitConfig
	fs     *afero.Afero
	regex  *regexp.Regexp
}

// Ensure LocationManager satisfies the Service interface at compile time.
var _ location.Service = (*LocationManager)(nil)

// NewLocationManager creates a LocationManager that reads and writes to the
// given gitconfig file using the provided filesystem. When fs is nil the OS
// filesystem is used; when gitconfigPath is empty the user's ~/.gitconfig is used.
func NewLocationManager(gitconfigPath string, fs *afero.Afero) *LocationManager {
	if fs == nil {
		fs = filesystem.NewOsFs()
	}
	return &LocationManager{
		gitcfg: gitconfig.NewGitConfig(gitconfigPath, fs),
		fs:     fs,
		regex:  regexp.MustCompile(regexStrGitConfigLocationGroups),
	}
}

// GetLocations returns all configured locations found in the gitconfig file.
func (lm LocationManager) GetLocations() ([]location.Location, error) {
	locations := []location.Location{}

	fileContent, err := lm.gitcfg.GetContent()
	if err != nil {
		return nil, err
	}

	result := lm.regex.FindAllStringSubmatch(fileContent, -1)
	for _, group := range result {
		locations = append(locations, location.Location{
			// TrimSuffix normalises paths captured from the file (e.g. "~/code/")
			// so that ToSection() does not produce a double trailing slash.
			Key:        strings.TrimSpace(group[1]),
			Path:       strings.TrimSuffix(strings.TrimSpace(group[2]), "/"),
			ConfigFile: strings.TrimSpace(group[3]),
		})
	}

	return locations, nil
}

// FindLocationByKey searches for a location by key. Returns nil, nil if not found.
func (lm *LocationManager) FindLocationByKey(key string) (*location.Location, error) {
	locations, err := lm.GetLocations()
	if err != nil {
		return nil, err
	}

	for _, l := range locations {
		if strings.EqualFold(l.Key, key) {
			return &l, nil
		}
	}
	return nil, nil
}

// SaveLocation creates a new location entry for the given key and path, or
// updates the existing one if the key is already present.
func (lm LocationManager) SaveLocation(key string, location_path string) error {
	l, err := lm.FindLocationByKey(key)
	if err != nil {
		return err
	}

	newLocation := location.NewLocation(key, location_path)
	newSection, err := newLocation.ToSection()
	if err != nil {
		return err
	}

	// TODO: Currently if issues moving/creating new config file for location
	// in case gitconfig is updated but an error happens on config file write
	// gitconfig and location config will be unsync
	if l != nil {
		// UPDATE MODE
		actualContent, err := lm.gitcfg.GetContent()
		if err != nil {
			return err
		}
		oldSection, err := l.ToSection()
		if err != nil {
			return err
		}
		newContent := strings.ReplaceAll(actualContent, oldSection, newSection)
		if err := lm.gitcfg.WriteContent(newContent); err != nil {
			return err
		}
		if l.ConfigFile != newLocation.ConfigFile {
			if err := lm.fs.Rename(l.ConfigFile, newLocation.ConfigFile); err != nil {
				return err
			}
		}
	} else {
		// NEW MODE
		if err := lm.gitcfg.AppendSection(newSection); err != nil {
			return err
		}
		f, err := lm.fs.OpenFile(newLocation.ConfigFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(templates.GitConfigTemplateFileContent); err != nil {
			return err
		}
	}

	return nil
}

// DeleteLocation removes the location identified by key from the gitconfig file
// and deletes its associated config file.
func (lm *LocationManager) DeleteLocation(key string) error {
	l, err := lm.FindLocationByKey(key)
	if err != nil {
		return err
	}
	if l == nil {
		return fmt.Errorf("location %s not found", key)
	}

	actualContent, err := lm.gitcfg.GetContent()
	if err != nil {
		return err
	}

	// Remove the raw matched section for this key directly from the file content,
	// avoiding round-trip through ToSection() which would produce a mismatched path.
	newContent := lm.regex.ReplaceAllStringFunc(actualContent, func(match string) string {
		groups := lm.regex.FindStringSubmatch(match)
		if len(groups) > 1 && strings.EqualFold(strings.TrimSpace(groups[1]), key) {
			return ""
		}
		return match
	})

	if err := lm.gitcfg.WriteContent(newContent); err != nil {
		return err
	}

	return lm.fs.Remove(l.ConfigFile)
}
