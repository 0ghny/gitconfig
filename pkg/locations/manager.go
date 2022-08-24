package locations

import (
	"io/fs"
	"os"
	"regexp"
	"strings"

	"github.com/0ghny/gitconfig/internal/filesystem"
	"github.com/0ghny/gitconfig/internal/templates"
	"github.com/0ghny/gitconfig/pkg/gitconfig"
	"github.com/spf13/afero"
)

const (
	regexStrGitConfigLocationGroups string      = `(?m)^#[\s]*gitconfig.location.key[\s]*(?P<key>.*)[\s]*\[includeIf[\s]* \"gitdir:(?P<dir>.*)\"\][\s]*[\s]*path[\s]*=[\s]*(?P<path>.*)$`
	configFilePerms                 fs.FileMode = fs.FileMode(int(0644))
)

type LocationManager struct {
	gitconfigPath                string
	fs                           *afero.Afero
	regexGitConfigLocationGroups *regexp.Regexp
}

func NewLocationManager(gitconfigPath string, fs *afero.Afero) *LocationManager {
	// If not path, use user default
	if gitconfigPath == "" {
		gitconfigPath = gitconfig.GetUserGitConfigPath()
	}
	// if not afero specified, use default (OS)
	if fs == nil {
		fs = filesystem.NewOsFs()
	}
	return &LocationManager{
		gitconfigPath:                gitconfigPath,
		fs:                           fs,
		regexGitConfigLocationGroups: regexp.MustCompile(regexStrGitConfigLocationGroups),
	}
}

// Get all configured locations in specified gitconfig file
func (t LocationManager) GetLocations() ([]Location, error) {
	locations := []Location{}

	fileContent, err := t.getGitConfigFileContent()
	if err != nil {
		return nil, err
	}

	// Use regex to get all matches with groups
	result := t.regexGitConfigLocationGroups.FindAllStringSubmatch(fileContent, -1)
	for _, group := range result {
		locations = append(locations, Location{
			Key:        group[1],
			Path:       group[2],
			ConfigFile: group[3],
		})
	}

	return locations, nil
}

// FindLocationByKey search the specified location in configured gitconfig file by key.
// Accepts the key to search as string as parameter
// It returns if found, the location information but also an error in case of any issue
func (t *LocationManager) FindLocationByKey(key string) (*Location, error) {
	locations, err := t.GetLocations()

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

func (t LocationManager) SaveLocation(key string, location string) error {
	// Seach if key is already present in file
	l, err := t.FindLocationByKey(key)
	if err != nil {
		return err
	}

	// Get the content
	actualContent, err := t.getGitConfigFileContent()
	if err != nil {
		return err
	}

	// Creates the new location object
	newLocation := NewLocation(key, location)
	newSection, err := newLocation.ToSection()
	if err != nil {
		return err
	}
	// TODO: Currently if issues moving/creating new config file for location
	// in case gitconfig is updated but an error happens on config file write
	// gitconfig and location config will be unsync
	if l != nil {
		// UPDATE MODE
		oldSection, err := l.ToSection()
		if err != nil {
			return err
		}
		newContent := strings.ReplaceAll(actualContent, oldSection, newSection)
		t.fs.WriteFile(gitconfig.GetUserGitConfigPath(), []byte(newContent), 0644)

		// Now time to (if location is different from previous, to move the file)
		if l.ConfigFile != newLocation.ConfigFile {
			err := t.fs.Rename(l.ConfigFile, newLocation.ConfigFile)
			if err != nil {
				return err
			}
		}
	} else {
		// NEW MODE
		f, err := t.fs.OpenFile(t.gitconfigPath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			// error opening gitconfig file in append mode
			return err
		}
		defer f.Close()
		if _, err := f.WriteString(newSection); err != nil {
			// error writing location section to gitconfig
			return err
		}

		// Write new file into location from template
		f2, err := t.fs.OpenFile(newLocation.ConfigFile, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			// error opening gitconfig file in append mode
			return err
		}
		defer f2.Close()
		if _, err := f2.WriteString(templates.GitConfigTemplateFileContent); err != nil {
			// error writing location gitconfig
			return err
		}
	}

	return nil
}

// TODO: Not sure if this method should be on gitconfig.go file
// with a proper DI, to support afero
// In this first version i will keep this simply here
func (t LocationManager) getGitConfigFileContent() (string, error) {
	fileBytes, err := t.fs.ReadFile(t.gitconfigPath)
	if err != nil {
		return "", err
	}
	return string(fileBytes), nil
}
