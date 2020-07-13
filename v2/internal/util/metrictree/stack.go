/*
 Copyright (c) 2020 SolarWinds Worldwide, LLC

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.
*/

package metrictree

type nodeStack []*Node

func (s *nodeStack) Push(v *Node) {
	*s = append(*s, v)
}

func (s *nodeStack) Pop() (*Node, bool) {
	if s.Len() == 0 {
		return nil, false
	}

	idx := s.Len() - 1
	node := (*s)[idx]
	(*s)[idx] = nil
	*s = (*s)[:idx]

	return node, true
}

func (s *nodeStack) Len() int {
	return len(*s)
}

func (s *nodeStack) Empty() bool {
	return s.Len() == 0
}
