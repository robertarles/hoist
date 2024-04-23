# hoist

Hoist duplicate files in a directory structure, removing duplication using symlinks.

Given a self contained directory structure, meaning no references outside the root directory, this util will find all duplicate files, hoist and replace with links.

This is tested on a very large (~18GB) directory with many subdirectories with duplicated html assets, which it reduced to about 6GB.

usage: `hoist <dirname>`

Will create a \<dirname\>/hoisted_resources directory for the hoisted files, leaving linked references to them in their origin.
