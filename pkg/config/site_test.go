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

package config

import (
	"bytes"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"testing"
)

type SiteSuite struct {
	suite.Suite

	testdataPath string
}

func (s *SiteSuite) SetupTest() {
	v := Create()
	cwd, cwdErr := os.Getwd()
	s.Require().NoError(cwdErr)

	basedir := filepath.Dir(filepath.Dir(cwd))
	s.testdataPath = filepath.Join(basedir, "testdata")
	v.ConfigPath = s.testdataPath

	SetGlobal(v)
}

func TestSiteSuite(t *testing.T) {
	suite.Run(t, new(SiteSuite))
}

func (s *SiteSuite) TestParseTemplateString_Simple() {
	t, err := loadTemplate("test", "<div>{{.}}</div>")

	s.NoError(err)
	s.Require().NotNil(t)

	buf := bytes.Buffer{}
	t.Execute(&buf, "TEST")

	s.Equal("<div>TEST</div>", buf.String())
}

func (s *SiteSuite) TestParseTemplateString_BadSyntax() {
	t, err := loadTemplate("test", "<div>{{.</div>")

	s.Error(err)
	s.Nil(t)
}

func (s *SiteSuite) TestInitTemplateString_Simple() {
	t := initTemplate("test", "<div>{{.}}</div>")

	s.Require().NotNil(t)

	buf := bytes.Buffer{}
	t.Execute(&buf, "TEST")

	s.Equal("<div>TEST</div>", buf.String())
}

func (s *SiteSuite) TestINitTemplateString_BadSyntax() {
	s.Panics(func() {
		t := initTemplate("test", "<div>{{.</div>")

		s.NotNil(t)
	})
}

func (s *SiteSuite) TestRenderTemplateParsing_String() {
	rt := RenderTemplate{tpl: initTemplate("test", "<div>{{.}}</div>")}

	err := yaml.Unmarshal([]byte("<section>{{.}}</section>"), &rt)

	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = rt.tpl.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<section>TEST</section>", buf.String())
}

func (s *SiteSuite) TestResolveTemplate_String() {
	t, err := resolveTemplate("test", "<section>{{.}}</section>")
	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = t.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<section>TEST</section>", buf.String())
}

func (s *SiteSuite) TestResolveTemplate_AbsPath() {
	tmplPath := filepath.Join(Global().ConfigPath, "templates", "div_wrapper.tpl.html")
	t, err := resolveTemplate("test", tmplPath)
	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = t.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<div>TEST</div>", buf.String())
}

func (s *SiteSuite) TestResolveTemplate_RelPath() {
	tmplPath := filepath.Join("templates", "section_wrapper.tpl.html")
	t, err := resolveTemplate("test", tmplPath)
	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = t.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<section>TEST</section>", buf.String())
}

func (s *SiteSuite) TestResolveTemplate_BadString() {
	t, err := resolveTemplate("test", "<section>{{{{.}}</section>")

	s.Error(err)
	s.Nil(t)
}

func (s *SiteSuite) TestResolveTemplate_Unknown() {
	t, err := resolveTemplate("test", "<section></section>")

	s.Error(err)
	s.Nil(t)
}

func (s *SiteSuite) TestCreateRenderTemplate_Failure() {
	s.Panics(func() {
		createRenderTemplate("super-fail", "non-existent-dir/no-file.tmp")
	})
}

func (s *SiteSuite) TestLoadTemplateFile_Missing() {
	tmplPath := filepath.Join(Global().ConfigPath, "templates", "missing.tpl.html")
	t, err := loadTemplateFile("test", tmplPath)

	s.Error(err)
	s.Nil(t)
}

func (s *SiteSuite) TestUnmarshal_String() {
	tmpl := "<section>{{.}}</section>"

	rt := RenderTemplate{}

	err := yaml.Unmarshal([]byte(tmpl), &rt)

	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = rt.tpl.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<section>TEST</section>", buf.String())
}

func (s *SiteSuite) TestUnmarshal_Invalid() {
	rt := RenderTemplate{}

	data := "foo: bar\nbaz: nil"
	err := yaml.Unmarshal([]byte(data), &rt)

	s.Require().Error(err)
}

func (s *SiteSuite) TestUnmarshal_AbsFile() {
	tmplPath := filepath.Join(Global().ConfigPath, "templates", "div_wrapper.tpl.html")

	rt := RenderTemplate{}

	err := yaml.Unmarshal([]byte(tmplPath), &rt)

	s.Require().NoError(err)

	buf := bytes.Buffer{}
	err = rt.Execute(&buf, "TEST")

	s.Require().NoError(err)
	s.Equal("<div>TEST</div>", buf.String())
}

func (s *SiteSuite) TestLoadSiteConfig() {
	Global().ConfigPath = filepath.Join(Global().ConfigPath, "sites/test01/config")
	site, err := LoadSiteConfig()

	s.Require().NoError(err)
	s.Require().NotNil(site)

	s.Equal("Test01", site.Title)

	buf := bytes.Buffer{}
	execErr := site.Markdown.BlockTemplate.Execute(&buf, "TEST")
	s.NoError(execErr)
	s.Equal(`<section class="pageContent markdown">TEST</section>`, buf.String())
}

func (s *SiteSuite) TestLoadSiteConfig_BadTemplateFormat() {
	Global().ConfigPath = filepath.Join(Global().ConfigPath, "sites/fail01/config")
	site, err := LoadSiteConfig()

	s.Error(err)
	s.NotNil(site)
	s.Equal("Fail01", site.Title)
}

func (s *SiteSuite) TestLoadSiteConfig_MissingFile() {
	Global().ConfigPath = filepath.Join(Global().ConfigPath, "sites/fail02/config")
	site, err := LoadSiteConfig()

	s.Require().Error(err)
	s.NotNil(site)
}
