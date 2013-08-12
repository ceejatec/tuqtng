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
	"fmt"
)

type ExpressionFormalNotationConverter struct {
	forbiddenAliases []string
	aliases          []string
	defaultAlias     string
}

// this visitor takes a list of valid aliases
// if there are any forbidden Aliases
// all property refereces MUST NOT start with one of these aliases
// if there is more than 1 alias in the list
// all property references MUST start with one of these aliases
// if not, an appropriate error is returned
// if there is only 1 alias, and the reference can be converted
// a new expression with the proper reference is returned
// it is up to the caller to update any references it may have
func NewExpressionFormalNotationConverter(forbiddenAliases []string, aliases []string, defaultAlias string) *ExpressionFormalNotationConverter {
	return &ExpressionFormalNotationConverter{
		forbiddenAliases: forbiddenAliases,
		aliases:          aliases,
		defaultAlias:     defaultAlias,
	}
}

func (this *ExpressionFormalNotationConverter) Visit(expr Expression) (Expression, error) {
	var err error
	switch expr := expr.(type) {
	// override recursive behavior for some types
	case *DotMemberOperator:
		_, err = this.VisitDotMemberChild(expr)
		if err != nil {
			return expr, err
		}
	case *CollectionAnyOperator:
		_, err = this.VisitCollectionAnyChild(expr)
		if err != nil {
			return expr, err
		}
	case *CollectionAllOperator:
		_, err = this.VisitCollectionAllChild(expr)
		if err != nil {
			return expr, err
		}
	// and add custom behavior to others (with default recursion)
	case *FunctionCall:
		// default processing of children
		_, err = VisitChildren(this, expr)
		if err != nil {
			return expr, err
		}
		// now custom processing at this node
		return this.VisitFunctionCall(expr)
	case *Property:
		// has no children anyway
		return this.VisitProperty(expr)
	default:
		return VisitChildren(this, expr)
	}
	return expr, nil
}

func (this *ExpressionFormalNotationConverter) VisitDotMemberChild(expr *DotMemberOperator) (Expression, error) {
	var err error
	expr.Left, err = expr.Left.Accept(this)
	if err != nil {
		return expr, err
	}
	return expr, nil
}

func (this *ExpressionFormalNotationConverter) VisitCollectionAnyChild(expr *CollectionAnyOperator) (Expression, error) {
	var err error
	expr.Over, err = expr.Over.Accept(this)
	if err != nil {
		return expr, err
	}

	updatedAliases := make([]string, len(this.aliases)+1)
	copy(this.aliases, updatedAliases)
	updatedAliases[len(this.aliases)] = expr.As
	childVisitor := NewExpressionFormalNotationConverter(this.forbiddenAliases, updatedAliases, this.defaultAlias)
	expr.Condition, err = expr.Condition.Accept(childVisitor)
	return expr, err
}

func (this *ExpressionFormalNotationConverter) VisitCollectionAllChild(expr *CollectionAllOperator) (Expression, error) {
	var err error
	expr.Over, err = expr.Over.Accept(this)
	if err != nil {
		return expr, err
	}

	updatedAliases := make([]string, len(this.aliases)+1)
	copy(this.aliases, updatedAliases)
	updatedAliases[len(this.aliases)] = expr.As
	childVisitor := NewExpressionFormalNotationConverter(this.forbiddenAliases, updatedAliases, this.defaultAlias)
	expr.Condition, err = expr.Condition.Accept(childVisitor)
	return expr, err
}

func (this *ExpressionFormalNotationConverter) VisitFunctionCall(expr *FunctionCall) (Expression, error) {
	// two specific checks need to be made here for special functions
	if expr.Name == "VALUE" {
		// VALUE() with 0 args is converted to VALUE(defaultAlias) when there is one
		if len(expr.Operands) == 0 && this.defaultAlias != "" {
			expr.Operands = append(expr.Operands, NewFunctionArgExpression(NewProperty(this.defaultAlias)))
		}
	}
	if expr.Name == "META" {
		// META() with 0 args is converted to META(defaultAlias) when there is one
		if len(expr.Operands) == 0 && this.defaultAlias != "" {
			expr.Operands = append(expr.Operands, NewFunctionArgExpression(NewProperty(this.defaultAlias)))
		} else if len(expr.Operands) > 0 {
			// check to see that the correct bucket is referenced (currently always aliases[0])
			switch operexpr := expr.Operands[0].Expr.(type) {
			case *Property:
				if operexpr.Path != this.aliases[0] {
					return expr, fmt.Errorf("invalid argument to META() function, must be bucket name/alias")
				}
			default:
				return expr, fmt.Errorf("invalid argument to META() function, must be bucket name/alias")
			}
		}
	}
	return expr, nil
}

func (this *ExpressionFormalNotationConverter) VisitProperty(expr *Property) (Expression, error) {
	// check to see if any of the forbiddenAliases are mentioned
	if len(this.forbiddenAliases) > 0 {
		for _, forbiddenAlias := range this.forbiddenAliases {
			if expr.Path == forbiddenAlias {
				return expr, fmt.Errorf("Alias %s cannot be referenced", expr.Path)
			}
		}
	}
	// this test is not needed when there are no aliases in the from clause (expression evaluation only)
	if len(this.aliases) > 0 {
		for _, alias := range this.aliases {
			if expr.Path == alias {
				return expr, nil
			}
		}
		if this.defaultAlias != "" {
			return NewDotMemberOperator(NewProperty(this.defaultAlias), expr), nil
		} else {
			return expr, fmt.Errorf("Property reference %s missing qualifier bucket/alias", expr.Path)
		}
	}
	return expr, nil
}