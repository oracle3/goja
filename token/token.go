// Package token defines constants representing the lexical tokens of JavaScript (ECMA5).
package token

import (
	"strconv"
)

// Token is the set of lexical tokens in JavaScript (ECMA5).
type Token int

// String returns the string corresponding to the token.
// For operators, delimiters, and keywords the string is the actual
// token string (e.g., for the token PLUS, the String() is
// "+"). For all other tokens the string corresponds to the token
// name (e.g. for the token IDENTIFIER, the string is "IDENTIFIER").
//
//String返回与标记对应的字符串。
//对于运算符、分隔符和关键字，字符串是实际的标记字符串（例如，对于标记PLUS，string（）是“+”）。
//对于所有其他标记，字符串对应于标记名称（例如，对于标记IDENTIFIER，字符串是“IDENTIFIER”）。
func (tkn Token) String() string {
	if tkn == 0 {
		return "UNKNOWN"
	}
	if tkn < Token(len(token2string)) {
		return token2string[tkn]
	}
	return "token(" + strconv.Itoa(int(tkn)) + ")"
}

// This is not used for anything
// 未使用的函数
func (tkn Token) precedence(in bool) int {

	switch tkn {
	case LOGICAL_OR:
		return 1

	case LOGICAL_AND:
		return 2

	case OR, OR_ASSIGN:
		return 3

	case EXCLUSIVE_OR:
		return 4

	case AND, AND_ASSIGN, AND_NOT, AND_NOT_ASSIGN:
		return 5

	case EQUAL,
		NOT_EQUAL,
		STRICT_EQUAL,
		STRICT_NOT_EQUAL:
		return 6

	case LESS, GREATER, LESS_OR_EQUAL, GREATER_OR_EQUAL, INSTANCEOF:
		return 7

	case IN:
		if in {
			return 7
		}
		return 0

	case SHIFT_LEFT, SHIFT_RIGHT, UNSIGNED_SHIFT_RIGHT:
		fallthrough
	case SHIFT_LEFT_ASSIGN, SHIFT_RIGHT_ASSIGN, UNSIGNED_SHIFT_RIGHT_ASSIGN:
		return 8

	case PLUS, MINUS, ADD_ASSIGN, SUBTRACT_ASSIGN:
		return 9

	case MULTIPLY, SLASH, REMAINDER, MULTIPLY_ASSIGN, QUOTIENT_ASSIGN, REMAINDER_ASSIGN:
		return 11
	}
	return 0
}

type _keyword struct {
	token         Token
	futureKeyword bool // 将来是否是关键字
	strict        bool // 在严格模式下是否是关键字
}

// IsKeyword returns the keyword token if literal is a keyword, a KEYWORD token
// if the literal is a future keyword (const, let, class, super, ...), or 0 if the literal is not a keyword.
//
// If the literal is a keyword, IsKeyword returns a second value indicating if the literal
// is considered a future keyword in strict-mode only.
//
// 7.6.1.2 Future Reserved Words:
//
//       const
//       class
//       enum
//       export
//       extends
//       import
//       super
//
// 7.6.1.2 Future Reserved Words (strict):
//
//       implements
//       interface
//       let
//       package
//       private
//       protected
//       public
//       static
//
//IsKeyword返回关键字标记，如果literal是关键字，则返回关键字标记
//如果是将来的关键字（const，let，class，super，…），或者如果文本不是关键字，则为0。
//如果是关键字，则返回的第二个值指示literal仅在严格模式下被视为未来关键字。

func IsKeyword(literal string) (Token, bool) {
	if keyword, exists := keywordTable[literal]; exists {
		if keyword.futureKeyword {
			return KEYWORD, keyword.strict
		}
		return keyword.token, false
	}
	return 0, false
}
