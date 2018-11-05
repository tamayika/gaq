package query

import (
	"testing"

	"github.com/alecthomas/participle/lexer"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
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
			"*",
			args{
				q: "*",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "*",
							},
						},
					},
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
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package, Package",
			args{
				q: "Package, Package",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
						},
					},
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 10,
							Offset: 9,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 10,
									Offset: 9,
								},
								Name: "Package",
							},
						},
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
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 9,
									Offset: 8,
								},
								Name: "Ident",
							},
						},
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
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
							&SimpleSelector{
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
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
							&SimpleSelector{
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
				},
			},
			false,
		},
		{
			"Package ~ Ident",
			args{
				q: "Package ~ Ident",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
							},
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 9,
									Offset: 8,
								},
								Combinator: "~",
								Name:       "Ident",
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name]",
			args{
				q: "Package[Name]",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "Name",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name='foo']",
			args{
				q: "Package[Name='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			`Package[Name="foo"]`,
			args{
				q: `Package[Name="foo"]`,
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name~='foo']",
			args{
				q: "Package[Name~='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "~=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name~='foo']",
			args{
				q: "Package[Name~='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "~=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name|='foo']",
			args{
				q: "Package[Name|='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "|=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name^='foo']",
			args{
				q: "Package[Name^='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "^=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name$='foo']",
			args{
				q: "Package[Name$='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "$=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package[Name*='foo']",
			args{
				q: "Package[Name*='foo']",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Attribute: &SimpleSelectorOptionAttribute{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name:     "Name",
											Operator: "*=",
											Value:    "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package:foo",
			args{
				q: "Package:foo",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Pseudo: &SimpleSelectorOptionPseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "foo",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package:foo-bar",
			args{
				q: "Package:foo-bar",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Pseudo: &SimpleSelectorOptionPseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "foo-bar",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package:foo(bar)",
			args{
				q: "Package:foo(bar)",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Pseudo: &SimpleSelectorOptionPseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "foo",
											Expressions: []*Expression{
												&Expression{
													Pos: lexer.Position{
														Line:   1,
														Column: 13,
														Offset: 12,
													},
													Name: "bar",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package:foo('bar')",
			args{
				q: "Package:foo('bar')",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Pseudo: &SimpleSelectorOptionPseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "foo",
											Expressions: []*Expression{
												&Expression{
													Pos: lexer.Position{
														Line:   1,
														Column: 13,
														Offset: 12,
													},
													String: "bar",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"Package:foo(10)",
			args{
				q: "Package:foo(10)",
			},
			&Query{
				lexer.Position{
					Line:   1,
					Column: 1,
				},
				[]*Selector{
					&Selector{
						lexer.Position{
							Line:   1,
							Column: 1,
						},
						[]*SimpleSelector{
							&SimpleSelector{
								Pos: lexer.Position{
									Line:   1,
									Column: 1,
								},
								Name: "Package",
								Options: []*SimpleSelectorOption{
									&SimpleSelectorOption{
										Pos: lexer.Position{
											Line:   1,
											Column: 8,
											Offset: 7,
										},
										Pseudo: &SimpleSelectorOptionPseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Name: "foo",
											Expressions: []*Expression{
												&Expression{
													Pos: lexer.Position{
														Line:   1,
														Column: 13,
														Offset: 12,
													},
													Number: "10",
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args.q)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseQuery() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
