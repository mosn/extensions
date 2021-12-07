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

package safe

import (
	"errors"
	"fmt"
	"sync"
)

// IntMap Used to store type mappings of string and uint64 and is thread safe.
// This is especially useful in protocol scenarios where string ID identifiers are used
type IntMap struct {
	table map[string]uint64 // id -> encoded stream id
	lock  sync.RWMutex      // protect table
}

func (m *IntMap) Get(key string) (val uint64, found bool) {

	m.lock.RLock()
	if len(m.table) <= 0 {
		return 0, false
	}

	val, found = m.table[key]

	m.lock.RUnlock()
	return
}

func (m *IntMap) Put(key string, val uint64) (err error) {

	m.lock.Lock()
	if m.table == nil {
		m.table = make(map[string]uint64, 8)
	}

	if v, found := m.table[key]; found {
		m.lock.Unlock()
		return errors.New(fmt.Sprintf("val conflict, exist key %s, val %d, current %d", key, v, val))
	}

	m.table[key] = val
	m.lock.Unlock()
	return
}

func (m *IntMap) Remove(key string) (err error) {

	m.lock.Lock()
	if m.table != nil {
		delete(m.table, key)
	}

	m.lock.Unlock()
	return
}
