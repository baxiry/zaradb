package parser

import (
	"testing"
)

func Test_parse(t *testing.T) {
	query := "_select * _from users _where age = 23 _orderby age _offset 10 _limit 20"
	parse(query)

	query2 := " _select    * _from users _where age = 23 _orderby age _offset 10 _limit 20"
	parse(query2)
}
