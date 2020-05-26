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
	"github.com/zpxio/mdsite/pkg/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

type PageIndex struct {
	PageLookup    map[string]*PageEntry
	Pages         []*PageEntry
	WeightLookup  map[string]float64
	DefaultWeight float64
	Title         string
}

var indexInit sync.Once
var globalIndex *PageIndex

func Index() *PageIndex {
	indexInit.Do(func() {
		ReIndex()
	})

	return globalIndex
}

func ReIndex() *PageIndex {
	var err error
	globalIndex, err = BuildIndex()

	if err != nil {
		panic("failed to build site index")
	}

	return globalIndex
}

func BuildIndex() (*PageIndex, error) {
	i := PageIndex{
		PageLookup:    make(map[string]*PageEntry),
		Pages:         []*PageEntry{},
		WeightLookup:  make(map[string]float64),
		DefaultWeight: DefaultWeight,
		Title:         "Contents",
	}

	// Read order data
	i.readOrder()

	err := filepath.Walk(config.Global().SitePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				// Ignore directories
				return nil
			}

			relPath, _ := filepath.Rel(config.Global().SitePath, path)
			log.Infof("Looking at file: %s", relPath)
			i.addResource(relPath)

			return nil
		})
	if err != nil {
		return nil, err
	}

	i.calculateOrder()

	return &i, nil
}

type OrderInfo struct {
	DefaultWeight float64  `yaml:"default"`
	OrderOrigin   int      `yaml:"orderMin"`
	Order         []string `yaml:"order"`
}

func (i *PageIndex) readOrder() {
	var order = OrderInfo{
		DefaultWeight: DefaultWeight,
		OrderOrigin:   1,
		Order:         []string{},
	}

	// Zero the order
	i.WeightLookup = make(map[string]float64)
	i.DefaultWeight = order.DefaultWeight

	// Read the file
	orderFile := filepath.Join(config.Global().ConfigPath, "order.yml")
	orderData, err := ioutil.ReadFile(orderFile)
	if err != nil {
		// Do nothing.
		return
	}
	log.Infof("Read file order from: %s", orderFile)

	err = yaml.Unmarshal(orderData, &order)
	if err != nil {
		// Do nothing, but log the abnormal error
		log.Errorf("Error while parsing order info: %s", err)
		return
	}

	// Assign weights by index
	for wd, u := range order.Order {
		i.WeightLookup[u] = float64(order.OrderOrigin + wd)
	}
}

func (i *PageIndex) addResource(path string) {
	p, err := LoadPageEntry(path)
	if err != nil {
		log.Errorf("Failed to load resource [%s]: %s", path, err)
	}

	w, ok := i.WeightLookup[p.Path]
	if ok {
		p.ListWeight = w
	} else {
		p.ListWeight = DefaultWeight
	}

	i.PageLookup[p.Url] = p
}

func (i *PageIndex) calculateOrder() {
	ordered := make([]*PageEntry, len(i.PageLookup))

	j := 0
	for _, pe := range i.PageLookup {
		ordered[j] = pe
		j++
	}

	sort.Slice(ordered, func(a, b int) bool {
		if ordered[a].ListWeight == ordered[b].ListWeight {
			return ordered[a].Url < ordered[b].Url
		}

		return ordered[a].ListWeight < ordered[b].ListWeight
	})

	i.Pages = ordered
}
