/*
Beside language-specific parsing this package performs the analysis of the source code.

The main goal is to present the user with a nice dependency graph and/or hint at cycles.
We should not expect the user to have much knowledge about the code to analyze.
Otherwise she or he would not need this tool. This we must provide decent default
behavior without asking many questions.

The implemented algorithm in a nutshell (see Glossary in README):

Step 1) Find all source-units
We collect all source-units and determine the longest shared package prefix.

Step 2) Collect all dependencies
We collect all dependencies of every source-unit. Dependencies to stuff outside
the root package is dropped.
All other dependencies are cropped according to the detail level (default is
length of root package plus one). We also crop the source-unit.
Remaining dependencies from a caller to itself are dropped as well.
Source-Units become nodes and dependencies become edges in a dependency graph.

Step 3) Analyse dependency graph
Now we search for cycles and analyse the graph. Each node is weighted with the number
of descendants and flaged if it is part of a cycle. We use the weights and flags
for rendering.
*/
package analysis
