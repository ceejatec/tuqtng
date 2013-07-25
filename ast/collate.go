//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

// this code originall taken from walrus
// https://github.com/couchbaselabs/walrus

package ast

import (
	"fmt"
	"sort"

	"code.google.com/p/go.exp/locale/collate"
	"code.google.com/p/go.text/locale"
	"github.com/couchbaselabs/tuqtng/query"
)

var icuCollator = collate.New(locale.Make("icu"))

// CouchDB-compatible collation/comparison of JSON values.
// See: http://wiki.apache.org/couchdb/View_collation#Collation_Specification
func CollateJSON(key1, key2 interface{}) int {
	type1 := collationType(key1)
	type2 := collationType(key2)
	if type1 != type2 {
		return type1 - type2
	}
	switch type1 {
	case 0, 1, 2:
		return 0
	case 3:
		n1 := collationToFloat64(key1)
		n2 := collationToFloat64(key2)
		if n1 < n2 {
			return -1
		} else if n1 > n2 {
			return 1
		}
		return 0
	case 4:
		s1 := key1.(string)
		s2 := key2.(string)
		return icuCollator.CompareString(s1, s2)
	case 5:
		array1 := key1.([]query.Value)
		array2 := key2.([]query.Value)
		for i, item1 := range array1 {
			if i >= len(array2) {
				return 1
			}
			if cmp := CollateJSON(item1, array2[i]); cmp != 0 {
				return cmp
			}
		}
		return len(array1) - len(array2)
	case 6:
		obj1 := key1.(map[string]query.Value)
		obj2 := key2.(map[string]query.Value)

		// first see if one object is larger than the other
		if len(obj1) < len(obj2) {
			return -1
		} else if len(obj1) > len(obj2) {
			return 1
		}

		// if not, proceed to do key by ke comparision

		// collect all the keys
		allkeys := make(sort.StringSlice, 0)
		for k, _ := range obj1 {
			allkeys = append(allkeys, k)
		}
		for k, _ := range obj2 {
			allkeys = append(allkeys, k)
		}

		// sort the keys
		allkeys.Sort()

		// now compare the values associated with each key
		for _, key := range allkeys {
			val1, ok := obj1[key]
			if !ok {
				// obj1 didn't have this key, so it is smaller
				return -1
			}
			val2, ok := obj2[key]
			if !ok {
				// ojb2 didnt have this key, so its smaller
				return 1
			}
			// key was in both objects, need to compare them
			comp := CollateJSON(val1, val2)
			if comp != 0 {
				// if this decided anything, return
				return comp
			}
			//otherwise continue to compare next key
		}
		// if got here, both objects are the same
		return 0
	}
	panic("bogus collationType")
}

func collationType(value interface{}) int {
	if value == nil {
		return 0
	}
	switch value := value.(type) {
	case bool:
		if !value {
			return 1
		}
		return 2
	case float64, uint64:
		return 3
	case string:
		return 4
	case []query.Value:
		return 5
	case map[string]query.Value:
		return 6
	}
	panic(fmt.Sprintf("collationType doesn't understand %+v of type %T", value, value))
}

func collationToFloat64(value interface{}) float64 {
	if i, ok := value.(uint64); ok {
		return float64(i)
	}
	if n, ok := value.(float64); ok {
		return n
	}
	panic(fmt.Sprintf("collationToFloat64 doesn't understand %+v", value))
}
