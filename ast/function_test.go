//  Copyright (c) 2013 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package ast

import (
	"reflect"
	"testing"

	"github.com/couchbaselabs/tuqtng/query"
)

func TestFunction(t *testing.T) {

	sampleContext := map[string]query.Value{
		"name": "will",
	}
	sampleMeta := map[string]query.Value{
		"id": "first",
	}

	context := query.NewMapItem(sampleContext, sampleMeta)

	tests := []struct {
		input  Expression
		output query.Value
		err    error
	}{
		{
			NewFunctionCall("META", FunctionArgExpressionList{}),
			map[string]query.Value{
				"id": "first",
			},
			nil,
		},
		{
			NewFunctionCall("LOWER", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralString("HELLO"))}),
			"hello",
			nil,
		},
		{
			NewFunctionCall("UPPER", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralString("hello"))}),
			"HELLO",
			nil,
		},
		{
			NewFunctionCall("CEIL", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(5.8))}),
			6.0,
			nil,
		},
		{
			NewFunctionCall("FLOOR", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(5.8))}),
			5.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(5.8))}),
			6.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36))}),
			60.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36)), NewFunctionArgExpression(NewLiteralNumber(0))}),
			60.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36)), NewFunctionArgExpression(NewLiteralNumber(1))}),
			60.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36)), NewFunctionArgExpression(NewLiteralNumber(2))}),
			55.0,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36)), NewFunctionArgExpression(NewLiteralNumber(3))}),
			55.4,
			nil,
		},
		{
			NewFunctionCall("ROUND", FunctionArgExpressionList{NewFunctionArgExpression(NewLiteralNumber(55.36)), NewFunctionArgExpression(NewLiteralNumber(4))}),
			55.36,
			nil,
		},
	}

	for _, x := range tests {
		result, err := x.input.Evaluate(context)
		if !reflect.DeepEqual(err, x.err) {
			t.Fatalf("Expected error %v, got %v", x.err, err)
		}
		if !reflect.DeepEqual(result, x.output) {
			t.Errorf("Expected %t %v, got %t %v", x.output, x.output, result, result)
		}
	}

}
