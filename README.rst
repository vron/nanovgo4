MODIFICATION
============

This is a modification to the original library by Yoshiki Shibukawa to use a different OpenGL library and OpenGL 4.1 in particular. See the commit history for modifications. All other items have been left as is - but please note that this is not the original unmodified version of the code.


NanoVGo
=============

Pure golang implementation of `NanoVG <https://github.com/memononen/nanovg>`_. NanoVG is a vector graphics engine inspired by HTML5 Canvas API.

`DEMO <https://shibukawa.github.io/nanovgo/>`_

API Reference
---------------

See `GoDoc <https://godoc.org/github.com/vron/nanovgo4>`_

Porting Memo
--------------

* Root folder ``.go`` files

  Ported from NanoVG.

* ``fontstashmini/fontstash_mini.go``

  Ported from `fontstash <https://github.com/memononen/fontstash>`_. It includes only needed functions.

* ``fontstashmini/truetype``

  Copy from ``https://github.com/TheOnly92/fontstash.go`` (Public Domain)

License
----------

zlib license

Original (NanoVG) Author
---------------------------

* `Mikko Mononen <https://github.com/memononen>`_

Author
---------------

* `Yoshiki Shibukawa <https://github.com/shibukawa>`_

Contribution
----------------

* Moriyoshi Koizumi
* @hnakamur2
* @mattn_jp
* @hagat
* @h_doxas
* FSX
