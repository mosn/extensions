/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package xr

import (
	"mosn.io/api"
	"strconv"
	"strings"
)

type Matcher struct{}

func (m *Matcher) XrProtocolMatcher(data []byte) api.MatchResult {
	if len(data) < RequestHeaderLen {
		return api.MatchAgain
	}

	rawLen := strings.TrimLeft(string(data[0:8]), "0")
	if rawLen == "" {
		return api.MatchFailed
	}

	packetLen, err := strconv.Atoi(rawLen)
	// invalid packet length or not number
	if packetLen <= 0 || err != nil {
		return api.MatchFailed
	}

	return api.MatchSuccess
}
