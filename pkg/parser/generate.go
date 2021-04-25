package parser

//go:generate pigeon -optimize-parser -optimize-grammar -alternate-entrypoints NormalGroup,HeaderGroup,AttributesGroup,MacrosGroup,QuotesGroup,NoneGroup,ReplacementsGroup,SpecialCharactersGroup,FileLocation,IncludedFileLine,BlockAttributes,InlineAttributes,LineRanges,TagRanges -o parser.go parser.peg
