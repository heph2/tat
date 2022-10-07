# tat

Tat is a very basic Static site generator that use Gemtext, it
supports a simple *serve* useful for testing, an *init* command that
let you create a directory structure that `tat` can understand and a
*build* command that generate an out directory with static files.

## Install

The easiest way to use `tat` is through `nix run`

	$ nix run 'github:heph2/tat'

But you can also build it and use the executable locally using:

	$ nix build

## How To Use

	$ tat -init

Create a basic directory structure:
 - assets: for CSS, img and other static files
 - pages: pages like about, home, contacts, ecc. ecc
 - posts: where to put your posts
 - layouts: where to put HTML layouts
 
Then you have to modify the html templates on layouts according to
your needs, add posts and PROFIT!

	$ tat -serve
	
Listen on localhost:8080

	$ tat -build

Create an out directory with the files generated

## License
tat is licensed under the terms of the GNU General Public License, version3. For more information, see [LICENSE][].

[LICENSE]: LICENSE
