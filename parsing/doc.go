/*
All language specific code to read packages and dependencies form source files resides in this packages.
Grouping, cropping, filtering and analysis is done elsewhere.

The main entry point is the codeParser.go which delegates the work to the language-specific parsers.
New languages must be added there or the code won't be called. 
*/
package parsing
