# [](https://github.com/bytesparadise/libasciidoc/compare/v0.2.0...v) (2019-02-24)



# [0.2.0](https://github.com/bytesparadise/libasciidoc/compare/v0.1.0...v0.2.0) (2019-02-24)


### Bug Fixes

* **renderer:** avoid double encoding of document attributes ([#296](https://github.com/bytesparadise/libasciidoc/issues/296)) ([00c0132](https://github.com/bytesparadise/libasciidoc/commit/00c0132)), closes [#43](https://github.com/bytesparadise/libasciidoc/issues/43) [#43](https://github.com/bytesparadise/libasciidoc/issues/43) [#295](https://github.com/bytesparadise/libasciidoc/issues/295)
* **renderer:** do not always render preamble withing wrapper ([#299](https://github.com/bytesparadise/libasciidoc/issues/299)) ([76ea3f7](https://github.com/bytesparadise/libasciidoc/commit/76ea3f7)), closes [#298](https://github.com/bytesparadise/libasciidoc/issues/298)
* **renderer:** fix panic in ToC when doc has no section ([#285](https://github.com/bytesparadise/libasciidoc/issues/285)) ([f7ee178](https://github.com/bytesparadise/libasciidoc/commit/f7ee178)), closes [#284](https://github.com/bytesparadise/libasciidoc/issues/284)
* **renderer:** wrap continuing element in `<div>` ([#281](https://github.com/bytesparadise/libasciidoc/issues/281)) ([f94e69a](https://github.com/bytesparadise/libasciidoc/commit/f94e69a)), closes [#270](https://github.com/bytesparadise/libasciidoc/issues/270)
* **types:** attach child elements to correct parent in ordered list ([#294](https://github.com/bytesparadise/libasciidoc/issues/294)) ([8d72ae0](https://github.com/bytesparadise/libasciidoc/commit/8d72ae0)), closes [#293](https://github.com/bytesparadise/libasciidoc/issues/293)


### Features

* **parser:** support spaces and tabs ahead of single line comment ([#289](https://github.com/bytesparadise/libasciidoc/issues/289)) ([061eb82](https://github.com/bytesparadise/libasciidoc/commit/061eb82)), closes [#288](https://github.com/bytesparadise/libasciidoc/issues/288)
* **parser/renderer:** attach list item to ancestor ([#291](https://github.com/bytesparadise/libasciidoc/issues/291)) ([6d9eb0e](https://github.com/bytesparadise/libasciidoc/commit/6d9eb0e)), closes [#264](https://github.com/bytesparadise/libasciidoc/issues/264)
* **parser/renderer:** support checklists ([#262](https://github.com/bytesparadise/libasciidoc/issues/262)) ([34598af](https://github.com/bytesparadise/libasciidoc/commit/34598af)), closes [#244](https://github.com/bytesparadise/libasciidoc/issues/244)
* **parser/renderer:** support element ID prefix ([#302](https://github.com/bytesparadise/libasciidoc/issues/302)) ([9780fef](https://github.com/bytesparadise/libasciidoc/commit/9780fef)), closes [#300](https://github.com/bytesparadise/libasciidoc/issues/300)
* **parser/renderer:** support list separation  ([#274](https://github.com/bytesparadise/libasciidoc/issues/274)) ([d2945ab](https://github.com/bytesparadise/libasciidoc/commit/d2945ab)), closes [#263](https://github.com/bytesparadise/libasciidoc/issues/263)
* **renderer:** support 'start' attribute when rendering an ordered list ([#279](https://github.com/bytesparadise/libasciidoc/issues/279)) ([e7f692d](https://github.com/bytesparadise/libasciidoc/commit/e7f692d)), closes [#271](https://github.com/bytesparadise/libasciidoc/issues/271)
* **renderer:** support "Q and A" labeled lists ([#280](https://github.com/bytesparadise/libasciidoc/issues/280)) ([6be62cd](https://github.com/bytesparadise/libasciidoc/commit/6be62cd)), closes [#271](https://github.com/bytesparadise/libasciidoc/issues/271)
* **renderer:** support dropping of principal text in ordered list item ([#283](https://github.com/bytesparadise/libasciidoc/issues/283)) ([2387309](https://github.com/bytesparadise/libasciidoc/commit/2387309)), closes [#266](https://github.com/bytesparadise/libasciidoc/issues/266) [#265](https://github.com/bytesparadise/libasciidoc/issues/265)
* **renderer:** support predefined attributes ([#282](https://github.com/bytesparadise/libasciidoc/issues/282)) ([99581b5](https://github.com/bytesparadise/libasciidoc/commit/99581b5)), closes [#266](https://github.com/bytesparadise/libasciidoc/issues/266)
* **renderer:** support title on labeled lists ([#278](https://github.com/bytesparadise/libasciidoc/issues/278)) ([a50c637](https://github.com/bytesparadise/libasciidoc/commit/a50c637)), closes [#267](https://github.com/bytesparadise/libasciidoc/issues/267)



# [0.1.0](https://github.com/bytesparadise/libasciidoc/compare/39964e8...v0.1.0) (2019-01-02)


### Bug Fixes

* **build:** remove support for golang1.8 on travis-ci and appveyor ([#241](https://github.com/bytesparadise/libasciidoc/issues/241)) ([9afc556](https://github.com/bytesparadise/libasciidoc/commit/9afc556)), closes [#240](https://github.com/bytesparadise/libasciidoc/issues/240)
* **build:** update codecov config ([#135](https://github.com/bytesparadise/libasciidoc/issues/135)) ([d207759](https://github.com/bytesparadise/libasciidoc/commit/d207759)), closes [#134](https://github.com/bytesparadise/libasciidoc/issues/134)
* **build:** use optimized grammar when building/installing ([#159](https://github.com/bytesparadise/libasciidoc/issues/159)) ([8c08ab8](https://github.com/bytesparadise/libasciidoc/commit/8c08ab8)), closes [#158](https://github.com/bytesparadise/libasciidoc/issues/158)
* **cli:** command hangs when no arg is provided ([#239](https://github.com/bytesparadise/libasciidoc/issues/239)) ([af34129](https://github.com/bytesparadise/libasciidoc/commit/af34129)), closes [#236](https://github.com/bytesparadise/libasciidoc/issues/236)
* **doc:** fix broken links in README.adoc ([#92](https://github.com/bytesparadise/libasciidoc/issues/92)) ([cdf4e1c](https://github.com/bytesparadise/libasciidoc/commit/cdf4e1c)), closes [#91](https://github.com/bytesparadise/libasciidoc/issues/91)
* **parser:** avoid too much parsing for sections ([#129](https://github.com/bytesparadise/libasciidoc/issues/129)) ([6cc6f51](https://github.com/bytesparadise/libasciidoc/commit/6cc6f51)), closes [#121](https://github.com/bytesparadise/libasciidoc/issues/121)
* **parser:** broken literal block parsing ([#200](https://github.com/bytesparadise/libasciidoc/issues/200)) ([7012d2b](https://github.com/bytesparadise/libasciidoc/commit/7012d2b)), closes [#197](https://github.com/bytesparadise/libasciidoc/issues/197)
* **parser:** fix admonition paragraph parsing ([#90](https://github.com/bytesparadise/libasciidoc/issues/90)) ([b1adfb3](https://github.com/bytesparadise/libasciidoc/commit/b1adfb3)), closes [#88](https://github.com/bytesparadise/libasciidoc/issues/88)
* **parser:** fix parser failure on lists ([#233](https://github.com/bytesparadise/libasciidoc/issues/233)) ([7713b7a](https://github.com/bytesparadise/libasciidoc/commit/7713b7a)), closes [#230](https://github.com/bytesparadise/libasciidoc/issues/230) [#234](https://github.com/bytesparadise/libasciidoc/issues/234) [#235](https://github.com/bytesparadise/libasciidoc/issues/235)
* **parser:** fix parsing issue when processing 'article.adoc' ([#164](https://github.com/bytesparadise/libasciidoc/issues/164)) ([29a8985](https://github.com/bytesparadise/libasciidoc/commit/29a8985)), closes [#153](https://github.com/bytesparadise/libasciidoc/issues/153)
* **parser:** incorrect default image alt ([#201](https://github.com/bytesparadise/libasciidoc/issues/201)) ([d60a28c](https://github.com/bytesparadise/libasciidoc/commit/d60a28c)), closes [#198](https://github.com/bytesparadise/libasciidoc/issues/198)
* **parser:** increase bench timeout on Travis-ci ([#163](https://github.com/bytesparadise/libasciidoc/issues/163)) ([a3aca2e](https://github.com/bytesparadise/libasciidoc/commit/a3aca2e)), closes [#162](https://github.com/bytesparadise/libasciidoc/issues/162)
* **parser:** misapplied ordered list on paragraph ([#208](https://github.com/bytesparadise/libasciidoc/issues/208)) ([44ee222](https://github.com/bytesparadise/libasciidoc/commit/44ee222)), closes [#207](https://github.com/bytesparadise/libasciidoc/issues/207)
* **parser:** missing sublists ([#206](https://github.com/bytesparadise/libasciidoc/issues/206)) ([39c0af8](https://github.com/bytesparadise/libasciidoc/commit/39c0af8)), closes [#203](https://github.com/bytesparadise/libasciidoc/issues/203)
* **parser:** parse blank lines ([#13](https://github.com/bytesparadise/libasciidoc/issues/13)) ([9c84e23](https://github.com/bytesparadise/libasciidoc/commit/9c84e23))
* **parser:** support for quoted text in list items ([#167](https://github.com/bytesparadise/libasciidoc/issues/167)) ([d4fe363](https://github.com/bytesparadise/libasciidoc/commit/d4fe363)), closes [#161](https://github.com/bytesparadise/libasciidoc/issues/161)
* **parser:** Support line starting with `.` in delimited blocks ([#120](https://github.com/bytesparadise/libasciidoc/issues/120)) ([efbdd39](https://github.com/bytesparadise/libasciidoc/commit/efbdd39)), closes [#116](https://github.com/bytesparadise/libasciidoc/issues/116)
* **parser:** support multiple sections with level 0 ([#124](https://github.com/bytesparadise/libasciidoc/issues/124)) ([bf43f4c](https://github.com/bytesparadise/libasciidoc/commit/bf43f4c))
* **parser:** support unclosed delimited blocks ([#101](https://github.com/bytesparadise/libasciidoc/issues/101)) ([#104](https://github.com/bytesparadise/libasciidoc/issues/104)) ([3029837](https://github.com/bytesparadise/libasciidoc/commit/3029837))
* **parser:** support unordered lists on multiple levels ([#145](https://github.com/bytesparadise/libasciidoc/issues/145)) ([4554793](https://github.com/bytesparadise/libasciidoc/commit/4554793)), closes [#137](https://github.com/bytesparadise/libasciidoc/issues/137)
* **parser:** unrecognized footnote in paragraph ([#211](https://github.com/bytesparadise/libasciidoc/issues/211)) ([d659997](https://github.com/bytesparadise/libasciidoc/commit/d659997)), closes [#210](https://github.com/bytesparadise/libasciidoc/issues/210)
* **parser/renderer:** avoid extra spaces in literal blocks ([#193](https://github.com/bytesparadise/libasciidoc/issues/193)) ([e8a26b0](https://github.com/bytesparadise/libasciidoc/commit/e8a26b0)), closes [#188](https://github.com/bytesparadise/libasciidoc/issues/188)
* **parser/renderer:** unique section id ([#209](https://github.com/bytesparadise/libasciidoc/issues/209)) ([0adc6a1](https://github.com/bytesparadise/libasciidoc/commit/0adc6a1)), closes [#184](https://github.com/bytesparadise/libasciidoc/issues/184)
* **project:** remove `.test` files ([#132](https://github.com/bytesparadise/libasciidoc/issues/132)) ([644b4eb](https://github.com/bytesparadise/libasciidoc/commit/644b4eb)), closes [#130](https://github.com/bytesparadise/libasciidoc/issues/130)
* **renderer:** element IDs and document header ([#156](https://github.com/bytesparadise/libasciidoc/issues/156)) ([c3e3fbd](https://github.com/bytesparadise/libasciidoc/commit/c3e3fbd)), closes [#155](https://github.com/bytesparadise/libasciidoc/issues/155)
* **renderer:** fix table numbering when title is included ([#166](https://github.com/bytesparadise/libasciidoc/issues/166)) ([7f3a6e0](https://github.com/bytesparadise/libasciidoc/commit/7f3a6e0))
* **renderer:** infinite recursive call ([#80](https://github.com/bytesparadise/libasciidoc/issues/80)) ([daed6fc](https://github.com/bytesparadise/libasciidoc/commit/daed6fc))
* **renderer:** missing '</head>' tag ([#202](https://github.com/bytesparadise/libasciidoc/issues/202)) ([3e3ca78](https://github.com/bytesparadise/libasciidoc/commit/3e3ca78)), closes [#199](https://github.com/bytesparadise/libasciidoc/issues/199)
* **types:** tidy up initials func ([#81](https://github.com/bytesparadise/libasciidoc/issues/81)) ([9448be5](https://github.com/bytesparadise/libasciidoc/commit/9448be5))


### Features

* **build:** add makefile goal to verify the generated parser ([#126](https://github.com/bytesparadise/libasciidoc/issues/126)) ([15b4680](https://github.com/bytesparadise/libasciidoc/commit/15b4680))
* **build:** add windows and osx builds to Travis ([#224](https://github.com/bytesparadise/libasciidoc/issues/224)) ([94a8009](https://github.com/bytesparadise/libasciidoc/commit/94a8009)), closes [#223](https://github.com/bytesparadise/libasciidoc/issues/223) [#225](https://github.com/bytesparadise/libasciidoc/issues/225)
* **build:** use golangci-lint for all linting ([b07c3a7](https://github.com/bytesparadise/libasciidoc/commit/b07c3a7)), closes [#61](https://github.com/bytesparadise/libasciidoc/issues/61)
* **build/cmd:** include commit/tag and time in 'version' cmd ([#114](https://github.com/bytesparadise/libasciidoc/issues/114)) ([96409c3](https://github.com/bytesparadise/libasciidoc/commit/96409c3)), closes [#113](https://github.com/bytesparadise/libasciidoc/issues/113)
* **cli:** add arg to specify the output file ([#122](https://github.com/bytesparadise/libasciidoc/issues/122)) ([d402c2d](https://github.com/bytesparadise/libasciidoc/commit/d402c2d)), closes [#119](https://github.com/bytesparadise/libasciidoc/issues/119)
* **cmd:** add command line interface ([#78](https://github.com/bytesparadise/libasciidoc/issues/78)) ([2f6ae3b](https://github.com/bytesparadise/libasciidoc/commit/2f6ae3b)), closes [#60](https://github.com/bytesparadise/libasciidoc/issues/60)
* **cmd:** add flag to suppress header/footer ([#95](https://github.com/bytesparadise/libasciidoc/issues/95)) ([4a31775](https://github.com/bytesparadise/libasciidoc/commit/4a31775))
* **cmd:** add support to specify log level ([#85](https://github.com/bytesparadise/libasciidoc/issues/85)) ([47e6e3c](https://github.com/bytesparadise/libasciidoc/commit/47e6e3c))
* **cmd:** allow reading input from stdin ([#86](https://github.com/bytesparadise/libasciidoc/issues/86)) ([add3287](https://github.com/bytesparadise/libasciidoc/commit/add3287))
* **make:** add goal to build executable ([#94](https://github.com/bytesparadise/libasciidoc/issues/94)), show help by default ([#99](https://github.com/bytesparadise/libasciidoc/issues/99)) ([#103](https://github.com/bytesparadise/libasciidoc/issues/103)) ([3ea969a](https://github.com/bytesparadise/libasciidoc/commit/3ea969a))
* **parser:** add support for meta-elements: ID, link and title ([c08a7f3](https://github.com/bytesparadise/libasciidoc/commit/c08a7f3))
* **parser:** allow id and title on paragraphs ([#16](https://github.com/bytesparadise/libasciidoc/issues/16)) ([c499d94](https://github.com/bytesparadise/libasciidoc/commit/c499d94))
* **parser:** support double punctuation in quoted text ([#39](https://github.com/bytesparadise/libasciidoc/issues/39)) ([f7f82e9](https://github.com/bytesparadise/libasciidoc/commit/f7f82e9))
* **parser:** support front-matter in YAML format ([#28](https://github.com/bytesparadise/libasciidoc/issues/28)) ([b69fe01](https://github.com/bytesparadise/libasciidoc/commit/b69fe01))
* **parser:** support italic and monospace quotes, as well as nested quotes ([bd58fd1](https://github.com/bytesparadise/libasciidoc/commit/bd58fd1))
* **parser:** support relative links ([#65](https://github.com/bytesparadise/libasciidoc/issues/65)) ([5e47b65](https://github.com/bytesparadise/libasciidoc/commit/5e47b65)), closes [#56](https://github.com/bytesparadise/libasciidoc/issues/56)
* **parser:** support substitution prevention ([#40](https://github.com/bytesparadise/libasciidoc/issues/40)) ([8e59c45](https://github.com/bytesparadise/libasciidoc/commit/8e59c45))
* **parser:** use the `memoize` option in the parser to improve perfs ([#123](https://github.com/bytesparadise/libasciidoc/issues/123)) ([491dbdd](https://github.com/bytesparadise/libasciidoc/commit/491dbdd)), closes [#117](https://github.com/bytesparadise/libasciidoc/issues/117)
* **parser/renderer:** image blocks with metadata and paragraphs with multiple lines ([8ff1125](https://github.com/bytesparadise/libasciidoc/commit/8ff1125))
* **parser/renderer:** list item continuation ([#53](https://github.com/bytesparadise/libasciidoc/issues/53)) ([613a112](https://github.com/bytesparadise/libasciidoc/commit/613a112))
* **parser/renderer:** parse and render inline images ([#17](https://github.com/bytesparadise/libasciidoc/issues/17)) ([65f8ac7](https://github.com/bytesparadise/libasciidoc/commit/65f8ac7))
* **parser/renderer:** parse and render unordered list items ([#12](https://github.com/bytesparadise/libasciidoc/issues/12)) ([868e95a](https://github.com/bytesparadise/libasciidoc/commit/868e95a))
* **parser/renderer:** support admonitions ([#70](https://github.com/bytesparadise/libasciidoc/issues/70)) ([6c221f1](https://github.com/bytesparadise/libasciidoc/commit/6c221f1)), closes [#67](https://github.com/bytesparadise/libasciidoc/issues/67)
* **parser/renderer:** support block and paragraph quotes ([#157](https://github.com/bytesparadise/libasciidoc/issues/157)) ([9f1e394](https://github.com/bytesparadise/libasciidoc/commit/9f1e394)), closes [#141](https://github.com/bytesparadise/libasciidoc/issues/141)
* **parser/renderer:** support cross-references with Element ID ([#47](https://github.com/bytesparadise/libasciidoc/issues/47)) ([65f9c9c](https://github.com/bytesparadise/libasciidoc/commit/65f9c9c))
* **parser/renderer:** support example blocks ([#72](https://github.com/bytesparadise/libasciidoc/issues/72)) ([230febb](https://github.com/bytesparadise/libasciidoc/commit/230febb)), closes [#71](https://github.com/bytesparadise/libasciidoc/issues/71)
* **parser/renderer:** support explicit line breaks ([#195](https://github.com/bytesparadise/libasciidoc/issues/195)) ([f5f87cc](https://github.com/bytesparadise/libasciidoc/commit/f5f87cc)), closes [#189](https://github.com/bytesparadise/libasciidoc/issues/189)
* **parser/renderer:** support for delimited source blocks ([4cb7c14](https://github.com/bytesparadise/libasciidoc/commit/4cb7c14))
* **parser/renderer:** support for document attributes ([#22](https://github.com/bytesparadise/libasciidoc/issues/22)) ([362892a](https://github.com/bytesparadise/libasciidoc/commit/362892a))
* **parser/renderer:** support for document attributes reset and substitutions ([#23](https://github.com/bytesparadise/libasciidoc/issues/23)) ([f24fbd5](https://github.com/bytesparadise/libasciidoc/commit/f24fbd5))
* **parser/renderer:** Support for Document Author and Revision, and Preamble ([#36](https://github.com/bytesparadise/libasciidoc/issues/36)) ([99b1fd9](https://github.com/bytesparadise/libasciidoc/commit/99b1fd9))
* **parser/renderer:** support for labeled list ([#51](https://github.com/bytesparadise/libasciidoc/issues/51)) ([5e758c6](https://github.com/bytesparadise/libasciidoc/commit/5e758c6))
* **parser/renderer:** support for literal blocks ([#29](https://github.com/bytesparadise/libasciidoc/issues/29)) ([51f4897](https://github.com/bytesparadise/libasciidoc/commit/51f4897))
* **parser/renderer:** support inline footnotes ([#183](https://github.com/bytesparadise/libasciidoc/issues/183)) ([28e43c7](https://github.com/bytesparadise/libasciidoc/commit/28e43c7)), closes [#138](https://github.com/bytesparadise/libasciidoc/issues/138)
* **parser/renderer:** support links to section title ([#58](https://github.com/bytesparadise/libasciidoc/issues/58)) ([1900b10](https://github.com/bytesparadise/libasciidoc/commit/1900b10))
* **parser/renderer:** support listing blocks ([#42](https://github.com/bytesparadise/libasciidoc/issues/42)) ([2fb5fe6](https://github.com/bytesparadise/libasciidoc/commit/2fb5fe6))
* **parser/renderer:** support literal block attributes ([#186](https://github.com/bytesparadise/libasciidoc/issues/186)) ([4ef1381](https://github.com/bytesparadise/libasciidoc/commit/4ef1381)), closes [#185](https://github.com/bytesparadise/libasciidoc/issues/185)
* **parser/renderer:** support optional label in cross-references ([#174](https://github.com/bytesparadise/libasciidoc/issues/174)) ([ec85fd2](https://github.com/bytesparadise/libasciidoc/commit/ec85fd2)), closes [#143](https://github.com/bytesparadise/libasciidoc/issues/143)
* **parser/renderer:** support ordered lists ([#77](https://github.com/bytesparadise/libasciidoc/issues/77)) ([416e9ea](https://github.com/bytesparadise/libasciidoc/commit/416e9ea)), closes [#64](https://github.com/bytesparadise/libasciidoc/issues/64)
* **parser/renderer:** support passthrough ([#41](https://github.com/bytesparadise/libasciidoc/issues/41)) ([aa501da](https://github.com/bytesparadise/libasciidoc/commit/aa501da))
* **parser/renderer:** support role attributes, refactor attributes and image type ([#171](https://github.com/bytesparadise/libasciidoc/issues/171)) ([d2b6e95](https://github.com/bytesparadise/libasciidoc/commit/d2b6e95)), closes [#151](https://github.com/bytesparadise/libasciidoc/issues/151)
* **parser/renderer:** support sidebar blocks ([#182](https://github.com/bytesparadise/libasciidoc/issues/182)) ([e34547c](https://github.com/bytesparadise/libasciidoc/commit/e34547c)), closes [#139](https://github.com/bytesparadise/libasciidoc/issues/139)
* **parser/renderer:** support single line and block comments ([#146](https://github.com/bytesparadise/libasciidoc/issues/146)) ([c6549d3](https://github.com/bytesparadise/libasciidoc/commit/c6549d3)), closes [#144](https://github.com/bytesparadise/libasciidoc/issues/144)
* **parser/renderer:** support source code blocks with language ([#255](https://github.com/bytesparadise/libasciidoc/issues/255)) ([293761e](https://github.com/bytesparadise/libasciidoc/commit/293761e)), closes [#229](https://github.com/bytesparadise/libasciidoc/issues/229)
* **parser/renderer:** support subscript and superscript quotes ([#237](https://github.com/bytesparadise/libasciidoc/issues/237)) ([97e8929](https://github.com/bytesparadise/libasciidoc/commit/97e8929)), closes [#228](https://github.com/bytesparadise/libasciidoc/issues/228)
* **parser/renderer:** support tables (basic) ([#165](https://github.com/bytesparadise/libasciidoc/issues/165)) ([9956517](https://github.com/bytesparadise/libasciidoc/commit/9956517)), closes [#57](https://github.com/bytesparadise/libasciidoc/issues/57)
* **parser/renderer:** support TOC placement in preamble ([#45](https://github.com/bytesparadise/libasciidoc/issues/45)) ([b1a6a74](https://github.com/bytesparadise/libasciidoc/commit/b1a6a74))
* **parser/renderer:** support verses ([#149](https://github.com/bytesparadise/libasciidoc/issues/149)) ([ec67024](https://github.com/bytesparadise/libasciidoc/commit/ec67024)), closes [#140](https://github.com/bytesparadise/libasciidoc/issues/140)
* **parser/renderer:** verifies that `article.adoc` renders as expected ([#227](https://github.com/bytesparadise/libasciidoc/issues/227)) ([399d127](https://github.com/bytesparadise/libasciidoc/commit/399d127)), closes [#215](https://github.com/bytesparadise/libasciidoc/issues/215)
* **project:** first draft of the asciidoc grammar and parser ([39964e8](https://github.com/bytesparadise/libasciidoc/commit/39964e8))
* **renderer:** convert ellipsis ([#178](https://github.com/bytesparadise/libasciidoc/issues/178)) ([4733cfc](https://github.com/bytesparadise/libasciidoc/commit/4733cfc)), closes [#8230](https://github.com/bytesparadise/libasciidoc/issues/8230) [#8203](https://github.com/bytesparadise/libasciidoc/issues/8203) [#175](https://github.com/bytesparadise/libasciidoc/issues/175)
* **renderer:** render external links ([#48](https://github.com/bytesparadise/libasciidoc/issues/48)) ([1154a87](https://github.com/bytesparadise/libasciidoc/commit/1154a87))
* **renderer:** render external links without description ([#50](https://github.com/bytesparadise/libasciidoc/issues/50)) ([8457fa5](https://github.com/bytesparadise/libasciidoc/commit/8457fa5))
* **renderer:** render full document ([#18](https://github.com/bytesparadise/libasciidoc/issues/18)) ([bcdccfc](https://github.com/bytesparadise/libasciidoc/commit/bcdccfc))
* **renderer:** render headings with support for default and custom ID ([#10](https://github.com/bytesparadise/libasciidoc/issues/10)) ([76a05d4](https://github.com/bytesparadise/libasciidoc/commit/76a05d4))
* **renderer:** render section preamble ([#15](https://github.com/bytesparadise/libasciidoc/issues/15)) ([a897a73](https://github.com/bytesparadise/libasciidoc/commit/a897a73))
* **renderer:** render table of content ([#44](https://github.com/bytesparadise/libasciidoc/issues/44)) ([847f6a2](https://github.com/bytesparadise/libasciidoc/commit/847f6a2))
* **renderer:** render whole document ([baee941](https://github.com/bytesparadise/libasciidoc/commit/baee941))
* **renderer:** support 'imagesdir' attribute ([#170](https://github.com/bytesparadise/libasciidoc/issues/170)) ([852cca4](https://github.com/bytesparadise/libasciidoc/commit/852cca4)), closes [#160](https://github.com/bytesparadise/libasciidoc/issues/160)
* **renderer:** support icons in admonition blocks ([#218](https://github.com/bytesparadise/libasciidoc/issues/218)) ([aeef974](https://github.com/bytesparadise/libasciidoc/commit/aeef974)), closes [#214](https://github.com/bytesparadise/libasciidoc/issues/214)
* **renderer:** support ID and title on delimited blocks ([#213](https://github.com/bytesparadise/libasciidoc/issues/213)) ([8993045](https://github.com/bytesparadise/libasciidoc/commit/8993045)), closes [#212](https://github.com/bytesparadise/libasciidoc/issues/212)
* **renderer:** support inline attribute substitutions ([#179](https://github.com/bytesparadise/libasciidoc/issues/179)) ([d2f398e](https://github.com/bytesparadise/libasciidoc/commit/d2f398e)), closes [#142](https://github.com/bytesparadise/libasciidoc/issues/142)
* **rendering:** first draft of HTML5 rendering ([#3](https://github.com/bytesparadise/libasciidoc/issues/3)) ([b53b3a2](https://github.com/bytesparadise/libasciidoc/commit/b53b3a2))
* **rendering:** render italic and monospace quotes ([9ce2a48](https://github.com/bytesparadise/libasciidoc/commit/9ce2a48))


