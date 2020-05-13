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
	"errors"
	"fmt"
	"github.com/apex/log"
	"github.com/zpxio/mdsite/pkg/util"
	"gopkg.in/yaml.v2"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Site struct {
	Title    string               `yaml:"title"`
	Global   GlobalRenderConfig   `yaml:"global"`
	Markdown MarkdownRenderConfig `yaml:"markdown"`
	Html     HtmlRenderConfig     `yaml:"html"`
	Contents ContentsRenderConfig `yaml:"toc"`
}

type GlobalRenderConfig struct {
	PageTemplate *RenderTemplate `yaml:"pageTemplate"`
}

type HtmlRenderConfig struct {
	BlockTemplate *RenderTemplate `yaml:"blockTemplate"`
}

type MarkdownRenderConfig struct {
	BlockTemplate *RenderTemplate `yaml:"blockTemplate"`
}

type ContentsRenderConfig struct {
	PageTemplate *RenderTemplate `yaml:"pageTemplate"`
}

type RenderTemplate struct {
	tpl *template.Template
}

var templateSerial = 0

func nextTemplateName() string {
	templateSerial++
	return fmt.Sprintf("template-%05d", templateSerial)
}

func createRenderTemplate(name string, tpl string) *RenderTemplate {
	rt := RenderTemplate{}

	pt, err := resolveTemplate(name, tpl)
	if err != nil {
		panic(fmt.Sprintf("Failed to resolve template [%s]: %s", name, err))
	}
	rt.tpl = pt

	return &rt
}

func (t *RenderTemplate) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tmplSt string
	err := unmarshal(&tmplSt)
	if err != nil {
		return err
	}
	log.Infof("Unmarshaling: %s", tmplSt)

	rt, err := resolveTemplate(nextTemplateName(), tmplSt)
	if err != nil {
		return err
	}

	t.tpl = rt

	return nil
}

func (t *RenderTemplate) Execute(w io.Writer, data interface{}) error {
	return t.tpl.Execute(w, data)
}

func (t RenderTemplate) String() string {
	return fmt.Sprintf("Template[%s]", t.tpl.Name())
}

func resolveTemplate(name string, tmpl string) (*template.Template, error) {
	// Try an absolute path
	if filepath.IsAbs(tmpl) {
		// Load absolute file path
		return loadTemplateFile(name, tmpl)
	}

	// Try a relative path
	path := filepath.Join(Global().ConfigPath, tmpl)
	if _, err := os.Stat(path); err == nil {
		// Load the relative path
		return loadTemplateFile(name, path)
	}

	// See if the string looks like a template
	if strings.Contains(tmpl, "{{") {
		// Load the template as a string
		t, err := loadTemplate("", tmpl)
		if err != nil {
			return nil, err
		}
		return t, nil
	}

	return nil, errors.New(fmt.Sprintf("Could not resolve template: %s", util.StringPrefix(tmpl, 16)))
}

func loadTemplateFile(name string, tmplPath string) (*template.Template, error) {
	templateData, err := ioutil.ReadFile(tmplPath)
	if err != nil {
		return nil, err
	}

	log.Infof("Reading template file [%s]: %s", name, tmplPath)
	return loadTemplate(name, string(templateData))
}

func loadTemplate(name string, tmpl string) (*template.Template, error) {
	t := template.New(name)
	_, err := t.Parse(tmpl)
	log.Infof("Parsed: %s", tmpl)
	if err != nil {
		log.Errorf("Could not parse template [%s]: %s", name, err)
		return nil, err
	}

	return t, nil
}

func initTemplate(name string, tmpl string) *template.Template {
	t, err := loadTemplate(name, tmpl)

	if err != nil {
		log.Errorf("Failed to initialize template [%s]: %s", name, err)
		panic("Failed to initialize template[" + name + "]")
	}

	return t
}

func defaultSiteConfig() Site {
	s := Site{
		Title: "Default",
		Markdown: MarkdownRenderConfig{
			BlockTemplate: createRenderTemplate("md-default", `<div id="content markdown">{{.}}</div>"`),
		},
		Html: HtmlRenderConfig{
			BlockTemplate: createRenderTemplate("html-default", `<div id="content html">{{.}}</div>"`),
		},
	}

	return s
}

func LoadSiteConfig() (Site, error) {
	base := defaultSiteConfig()

	// Try to load file data
	siteFile := filepath.Join(Global().ConfigPath, "site.yml")
	data, err := ioutil.ReadFile(siteFile)
	if err != nil {
		return base, err
	}

	log.Infof("Loading site config: %s", siteFile)
	err = yaml.Unmarshal(data, &base)
	if err != nil {
		log.Errorf("Failed to load config: %s", err)
		return base, err
	}
	log.Infof("Config: %+v", base)

	return base, nil
}
