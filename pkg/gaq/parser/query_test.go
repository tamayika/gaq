package parser

import (
	"testing"

	"github.com/alecthomas/participle/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParseQuery(t *testing.T) {
	type args struct {
		q string
	}
	tests := []struct {
		name    string
		args    args
		want    *Query
		wantErr bool
	}{
		{
			"Empty",
			args{
				q: "",
			},
			&Query{
				Pos: lexer.Position{
					Line:   1,
					Column: 1,
				},
			},
			false,
		},
		{
			"Package",
			args{
				q: "Package",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Entry{
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 1,
						},
						Name: "Package",
					},
				},
			},
			false,
		},
		{
			"Package Ident",
			args{
				q: "Package Ident",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Entry{
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 1,
						},
						Name: "Package",
					},
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 9,
							Offset: 8,
						},
						Name: "Ident",
					},
				},
			},
			false,
		},
		{
			"Package > Ident",
			args{
				q: "Package > Ident",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Entry{
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 1,
						},
						Name: "Package",
					},
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 9,
							Offset: 8,
						},
						Combinator: ">",
						Name:       "Ident",
					},
				},
			},
			false,
		},
		{
			"Package + Ident",
			args{
				q: "Package + Ident",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Entry{
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 1,
						},
						Name: "Package",
					},
					&Entry{
						Pos: lexer.Position{
							Line:   1,
							Column: 9,
							Offset: 8,
						},
						Combinator: "+",
						Name:       "Ident",
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseQuery(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
