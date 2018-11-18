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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
										Attribute: &Attribute{
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
			"Package:first-child",
			args{
				q: "Package:first-child",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											FirstChild: &PseudoFirstChild{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
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
			"Package:first-of-type",
			args{
				q: "Package:first-of-type",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											FirstOfType: &PseudoFirstOfType{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
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
			"Package:has(File)",
			args{
				q: "Package:has(File)",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Has: &PseudoHas{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
												},
												Selectors: []*Selector{
													&Selector{
														Pos: lexer.Position{
															Line:   1,
															Column: 13,
															Offset: 12,
														},
														SimpleSelectors: []*SimpleSelector{
															&SimpleSelector{
																Pos: lexer.Position{
																	Line:   1,
																	Column: 13,
																	Offset: 12,
																},
																Name: "File",
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
					},
				},
			},
			false,
		},
		{
			"Package:is(Package)",
			args{
				q: "Package:is(Package)",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Is: &PseudoIs{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
												},
												Selectors: []*Selector{
													&Selector{
														Pos: lexer.Position{
															Line:   1,
															Column: 12,
															Offset: 11,
														},
														SimpleSelectors: []*SimpleSelector{
															&SimpleSelector{
																Pos: lexer.Position{
																	Line:   1,
																	Column: 12,
																	Offset: 11,
																},
																Name: "Package",
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
					},
				},
			},
			false,
		},
		{
			"Package:last-child",
			args{
				q: "Package:last-child",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											LastChild: &PseudoLastChild{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
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
			"Package:last-of-type",
			args{
				q: "Package:last-of-type",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											LastOfType: &PseudoLastOfType{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
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
			"Package:not(:first-child)",
			args{
				q: "Package:not(:first-child)",
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
										Pseudo: &Pseudo{
											Pos: lexer.Position{
												Line:   1,
												Column: 9,
												Offset: 8,
											},
											Not: &PseudoNot{
												Pos: lexer.Position{
													Line:   1,
													Column: 9,
													Offset: 8,
												},
												Selectors: []*Selector{
													&Selector{
														Pos: lexer.Position{
															Line:   1,
															Column: 13,
															Offset: 12,
														},
														SimpleSelectors: []*SimpleSelector{
															&SimpleSelector{
																Pos: lexer.Position{
																	Line:   1,
																	Column: 13,
																	Offset: 12,
																},
																Options: []*SimpleSelectorOption{
																	&SimpleSelectorOption{
																		Pos: lexer.Position{
																			Line:   1,
																			Column: 13,
																			Offset: 12,
																		},
																		Pseudo: &Pseudo{
																			Pos: lexer.Position{
																				Line:   1,
																				Column: 14,
																				Offset: 13,
																			},
																			FirstChild: &PseudoFirstChild{
																				Pos: lexer.Position{
																					Line:   1,
																					Column: 14,
																					Offset: 13,
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
