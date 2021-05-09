package parser

//go:generate pigeon -optimize-parser -optimize-grammar -alternate-entrypoints DocumentFragmentElementWithinParagraph,DocumentFragmentElementWithinDelimitedBlock,NormalGroup,HeaderGroup,AttributesGroup,MacrosGroup,QuotesGroup,NoneGroup,ReplacementsGroup,SpecialCharactersGroup,FileLocation,IncludedFileLine,BlockAttributes,InlineAttributes,LineRanges,TagRanges -o parser.go parser.peg
