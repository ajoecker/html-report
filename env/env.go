// Copyright 2015 ThoughtWorks, Inc.

// This file is part of getgauge/html-report.

// getgauge/html-report is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// getgauge/html-report is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with getgauge/html-report.  If not, see <http://www.gnu.org/licenses/>.

package env

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/getgauge/common"
)

const (
	DefaultReportsDir           = "reports"
	GaugeReportsDirEnvName      = "gauge_reports_dir" // directory where reports are generated by plugins
	OverwriteReportsEnvProperty = "overwrite_reports"
	UseNestedSpecs              = "use_nested_specs"
	ReportLocalisation          = "report_localisation"
)

func GetCurrentExecutableDir() (string, string) {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalf(err.Error())
	}
	target, err := filepath.EvalSymlinks(ex)
	if err != nil {
		return path.Dir(ex), filepath.Base(ex)
	}
	return filepath.Dir(target), filepath.Base(ex)
}

// CreateDirectory creates given directory if it doesn't exist
func CreateDirectory(dir string) {
	if err := os.MkdirAll(dir, common.NewDirectoryPermissions); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", dir, err)
		os.Exit(1)
	}
}

func GetProjectRoot() string {
	projectRoot := os.Getenv(common.GaugeProjectRootEnv)
	if projectRoot == "" {
		fmt.Printf("Environment variable '%s' is not set. \n", common.GaugeProjectRootEnv)
		os.Exit(1)
	}
	return projectRoot
}

func GetLocale() string {
	locale := os.Getenv(ReportLocalisation)
	if locale == "" {
		return "en"
	}
	return locale
}

func AddDefaultPropertiesToProject() {
	defaultPropertiesFile := getDefaultPropertiesFile()

	reportsDirProperty := &(common.Property{
		Comment:      "The path to the gauge reports directory. Should be either relative to the project directory or an absolute path",
		Name:         GaugeReportsDirEnvName,
		DefaultValue: DefaultReportsDir})

	overwriteReportProperty := &(common.Property{
		Comment:      "Set as false if gauge reports should not be overwritten on each execution. A new time-stamped directory will be created on each execution.",
		Name:         OverwriteReportsEnvProperty,
		DefaultValue: "true"})

	reportLocalisationProperty := &(common.Property{
		Comment:      "The localisation used for the report. Defaults to 'en'",
		Name:         ReportLocalisation,
		DefaultValue: "en"})

	if !common.FileExists(defaultPropertiesFile) {
		fmt.Printf("Failed to setup html report plugin in project. Default properties file does not exist at %s. \n", defaultPropertiesFile)
		return
	}
	if err := common.AppendProperties(defaultPropertiesFile, reportsDirProperty, overwriteReportProperty, reportLocalisationProperty); err != nil {
		fmt.Printf("Failed to setup html report plugin in project: %s \n", err)
		return
	}
	fmt.Println("Succesfully added configurations for html-report to env/default/default.properties")
}

func getDefaultPropertiesFile() string {
	return filepath.Join(GetProjectRoot(), "env", "default", "default.properties")
}

func ShouldOverwriteReports() bool {
	envValue := os.Getenv(OverwriteReportsEnvProperty)
	if strings.ToLower(envValue) == "true" {
		return true
	}
	return false
}

func ShouldUseNestedSpecs() bool {
	envValue := os.Getenv(UseNestedSpecs)
	if strings.ToLower(envValue) == "true" {
		return true
	}
	return false
}
