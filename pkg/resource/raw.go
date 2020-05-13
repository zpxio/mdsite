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

package resource

import (
	"io"
)

type RawResource struct {
}

func (r RawResource) MediaType() string {
	return "text/plain"
}

func (r RawResource) ResourceMode() string {
	return "text"
}

func (r RawResource) Render(w io.Writer, data *RenderData) error {
	_, err := io.WriteString(w, "RAW: ")
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, data.Resource)
	if err != nil {
		return err
	}

	return nil
}
