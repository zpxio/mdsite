/*
 * Copyright 2020 zpxio (Jeff Sharpe)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package site

import (
	"github.com/stretchr/testify/suite"
	"github.com/zpxio/mdsite/pkg/config"
	"os"
	"path/filepath"
	"testing"
)

type PageSuite struct {
	suite.Suite
}

func TestPageSuite(t *testing.T) {
	suite.Run(t, new(PageSuite))
}

func (s *PageSuite) SetupTest() {
	v := config.Create()
	cwd, cwdErr := os.Getwd()
	s.Require().NoError(cwdErr)

	basedir := filepath.Dir(filepath.Dir(cwd))
	testdataPath := filepath.Join(basedir, "testdata/sites/test01")
	v.ConfigPath = filepath.Join(testdataPath, "config")
	v.SitePath = filepath.Join(testdataPath, "site")
	config.SetGlobal(v)

	siteConf, siteErr := config.LoadSiteConfig()
	s.Require().NoError(siteErr)
	config.Global().SiteConfig = siteConf
}

func (s *PageSuite) TestGenerateLabel_Simple() {
	s.Equal("Simple", generateLabel("simple.md"))
}

func (s *PageSuite) TestGenerateLabel_SimpleSubdir() {
	s.Equal("Simple", generateLabel("sub/dir/simple.md"))
}

func (s *PageSuite) TestGenerateLabel_CapStart() {
	s.Equal("Caps", generateLabel("Caps.md"))
}

func (s *PageSuite) TestGenerateLabel_DashedWords() {
	s.Equal("Simple Page File", generateLabel("simple-page-file.md"))
}

func (s *PageSuite) TestGenerateLabel_MixedWords() {
	s.Equal("Simple Page File", generateLabel("simple page-file.md"))
}

func (s *PageSuite) TestGenerateLabel_CamelCase() {
	s.Equal("Simple Page File", generateLabel("simplePageFile.md"))
}

func (s *PageSuite) TestGenerateLabel_MixedCamelCase() {
	s.Equal("Simple Page File", generateLabel("simple pageFile.md"))
}

func (s *PageSuite) TestGenerateLabel_MultiCapCamelCase() {
	s.Equal("This Is ASCII", generateLabel("ThisIsASCII.md"))
}

func (s *PageSuite) TestLoadPageEntry_Happy() {
	filePath := "sample-01.md"
	fileStats, fsErr := os.Stat(filepath.Join(config.Global().SitePath, filePath))
	s.Require().NoError(fsErr)

	p, err := LoadPageEntry(filePath)

	s.Require().NoError(err)
	s.Equal(filePath, p.Path)
	s.Equal("/sample-01", p.Url)
	s.Equal("md", p.Extension)
	s.Equal("Sample 01", p.Label)
	s.Equal(fileStats.ModTime(), p.Modified)
	s.Greater(p.Id, uint64(0))
}

func (s *PageSuite) TestLoadPageEntry_SubDirFile() {
	filePath := "info/deep-file.txt"
	fileStats, fsErr := os.Stat(filepath.Join(config.Global().SitePath, filePath))
	s.Require().NoError(fsErr)

	p, err := LoadPageEntry(filePath)

	s.Require().NoError(err)
	s.Equal(filePath, p.Path)
	s.Equal("/info/deep-file", p.Url)
	s.Equal("txt", p.Extension)
	s.Equal("Deep File", p.Label)
	s.Equal(fileStats.ModTime(), p.Modified)
	s.Greater(p.Id, uint64(0))
}

func (s *PageSuite) TestLoadPageEntry_Missing() {
	filePath := "missing-01.txt"
	p, err := LoadPageEntry(filePath)

	s.Error(err)
	s.Nil(p)
}

func (s *PageSuite) TestLoadPageEntry_Dir() {
	filePath := "info"
	p, err := LoadPageEntry(filePath)

	s.Error(err)
	s.Nil(p)
}
