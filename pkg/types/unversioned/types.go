// Copyright 2018 Google Inc. All rights reserved.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package unversioned

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

type EnvVar struct {
	Key     string
	Value   string
	IsRegex bool `yaml:"isRegex"`
}

type Label struct {
	Key     string
	Value   string
	IsRegex bool `yaml:"isRegex"`
}

type Config struct {
	Env          map[string]string
	Entrypoint   []string
	Cmd          []string
	Volumes      []string
	Workdir      string
	ExposedPorts []string
	Labels       map[string]string
	User         string
}

type ContainerRunOptions struct {
	User         string
	Privileged   bool
	TTY          bool     `yaml:"allocateTty"`
	EnvVars      []string `yaml:"envVars"`
	EnvFile      string   `yaml:"envFile"`
	Capabilities []string
	BindMounts   []string `yaml:"bindMounts"`
}

func (opts *ContainerRunOptions) IsSet() bool {
	return len(opts.User) != 0 ||
		opts.Privileged ||
		opts.TTY ||
		len(opts.EnvFile) > 0 ||
		(opts.EnvVars != nil && len(opts.EnvVars) > 0) ||
		(opts.Capabilities != nil && len(opts.Capabilities) > 0) ||
		(opts.BindMounts != nil && len(opts.BindMounts) > 0)
}

type TestResult struct {
	Name     string        `xml:"name,attr"`
	Pass     bool          `xml:"-"`
	Stdout   string        `json:",omitempty" xml:"-"`
	Stderr   string        `json:",omitempty" xml:"-"`
	Errors   []string      `json:",omitempty" xml:"failure"`
	Duration time.Duration `xml:"time,attr"`
}

func (t *TestResult) String() string {
	strRepr := fmt.Sprintf("\nTest Name:%s", t.Name)
	testStatus := "Fail"
	if t.IsPass() {
		testStatus = "Pass"
	}
	strRepr += fmt.Sprintf("\nTest Status:%s", testStatus)
	if t.Stdout != "" {
		strRepr += fmt.Sprintf("\nStdout:%s", t.Stdout)
	}
	if t.Stderr != "" {
		strRepr += fmt.Sprintf("\nStderr:%s", t.Stderr)
	}
	strRepr += fmt.Sprintf("\nErrors:%s\n", strings.Join(t.Errors, ","))
	strRepr += fmt.Sprintf("\nDuration:%s\n", t.Duration.String())
	return strRepr
}

func (t *TestResult) Error(s string) {
	t.Errors = append(t.Errors, s)
}

func (t *TestResult) Errorf(s string, args ...interface{}) {
	t.Errors = append(t.Errors, fmt.Sprintf(s, args...))
}

func (t *TestResult) Fail() {
	t.Pass = false
}

func (t *TestResult) IsPass() bool {
	return t.Pass
}

type SummaryObject struct {
	XMLName  xml.Name      `json:"-" xml:"testsuites"`
	Pass     int           `xml:"-"`
	Fail     int           `xml:"failures,attr"`
	Total    int           `xml:"tests,attr"`
	Duration time.Duration `xml:"time,attr"`
	Results  []*TestResult `json:",omitempty" xml:"testsuite>testcase"`
}

type JUnitTestSuite struct {
	Name    string           `xml:"name,attr"`
	Results []*JUnitTestCase `xml:"testcase"`
}

type JUnitTestCase struct {
	Name     string   `xml:"name,attr"`
	Errors   []string `xml:"failure"`
	Duration float64  `xml:"time,attr"`
	Stdout   string   `xml:"system-out"`
	Stderr   string   `xml:"system-err"`
}

type OutputValue int

const (
	Text OutputValue = iota
	Json
	Junit
)

func (o OutputValue) String() string {
	return [...]string{"text", "json", "junit"}[o]
}

func (o OutputValue) Type() string {
	return "string"
}

func (o *OutputValue) Set(value string) error {
	switch value {
	case "text":
		*o = Text
	case "json":
		*o = Json
	case "junit":
		*o = Junit
	default:
		return fmt.Errorf("unsupported format %s: please select from `text`, `json`, or `junit`", value)
	}

	return nil
}
