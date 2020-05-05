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

import "github.com/gin-gonic/gin"

type MarkdownResource struct {
}

func (r *MarkdownResource) Render(c *gin.Context, path string) {
	c.Writer.Header().Set("X-QuickSite-Mode", "Markdown")

	c.Writer.WriteString("Path: ")
	c.Writer.WriteString(path)
}
