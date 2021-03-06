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
	"github.com/apex/log"
	"github.com/stretchr/testify/suite"
	"github.com/zpxio/mdsite/pkg/config"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

type SiteSuite struct {
	suite.Suite
}

func TestSiteSuite(t *testing.T) {
	suite.Run(t, new(SiteSuite))
}

func (s *SiteSuite) SetupTest() {
	s.loadSite("test01")
}

func (s *SiteSuite) loadSite(siteName string) {
	v := config.Create()
	cwd, cwdErr := os.Getwd()
	s.Require().NoError(cwdErr)

	basedir := filepath.Dir(filepath.Dir(cwd))
	testdataPath := filepath.Join(basedir, "testdata/sites/"+siteName)
	v.ConfigPath = filepath.Join(testdataPath, "config")
	v.SitePath = filepath.Join(testdataPath, "site")
	config.SetGlobal(v)

	siteConf, siteErr := config.LoadSiteConfig()
	s.Require().NoError(siteErr)
	config.Global().SiteConfig = siteConf
}

func (s *SiteSuite) TestCreateIndex_Simple() {
	i, err := BuildIndex()

	s.Require().NoError(err)
	s.NotNil(i)

	s.Len(i.PageLookup, 4)
	log.Infof("PageLookup: %+v", i.PageLookup)
	s.Equal("/extra", i.PageLookup["/extra"].Url)
	s.Equal("extra.txt", i.PageLookup["/extra"].Path)

	s.Len(i.Pages, 4)
	s.Equal("/sample-01", i.Pages[0].Url)

	s.Equal(i.DefaultWeight, i.PageLookup["/supplement"].ListWeight)
	s.Equal(i.DefaultWeight, i.PageLookup["/extra"].ListWeight)
}

func (s *SiteSuite) TestCreateIndex_NoOrderYaml() {
	s.loadSite("fail03")
	i, err := BuildIndex()

	s.Require().NoError(err)
	s.NotNil(i)

	s.Len(i.PageLookup, 1)
}

func (s *SiteSuite) TestCreateIndex_BadSitePath() {
	s.loadSite("fail03")
	config.Global().SitePath = "/@@@@/BadPATH"
	i, err := BuildIndex()

	s.Error(err)
	s.Nil(i)
}

func (s *SiteSuite) TestCreateIndex_BadOrderFile() {
	s.loadSite("fail04")
	i, err := BuildIndex()

	s.Require().NoError(err)
	s.NotNil(i)

	s.Equal(float64(DefaultWeight), i.Pages[0].ListWeight)
}

func (s *SiteSuite) TestIndex_Fails() {
	s.loadSite("fail03")
	config.Global().SitePath = "/@@@@/BadPATH"

	s.Panics(func() {
		i := ReIndex()
		s.Nil(i)
	})
}

func (s *SiteSuite) TestReIndex() {
	i := ReIndex()

	s.NotNil(i)

	i2 := ReIndex()
	s.NotEqual(reflect.ValueOf(i).Pointer(), reflect.ValueOf(i2).Pointer())
}

func (s *SiteSuite) TestIndex_Singleton() {
	i := Index()

	s.NotNil(i)

	i2 := Index()
	s.Equal(reflect.ValueOf(i).Pointer(), reflect.ValueOf(i2).Pointer())
}

func (s *SiteSuite) TestIndex_PostReIndex() {
	i0 := Index()
	i := ReIndex()

	s.NotNil(i)

	i2 := Index()
	s.NotEqual(reflect.ValueOf(i).Pointer(), reflect.ValueOf(i0).Pointer())
	s.Equal(reflect.ValueOf(i).Pointer(), reflect.ValueOf(i2).Pointer())
}
