giostorage
------

That package was created to cache some information, because nobody likes to wait 
for a http response. This package was designed to save small data, like user
preferences, profile pictures and related names.

Due to privacy reasons, giostorage also provides an encrypted implementation. That
is useful to keep way others users (or others unprivileged apps) from reading the content.

It's currently being used on WebAssembly, Windows and Android. _It may work on Linux, Darwin 
and iOS, but it's not tested._

-------

**Known-Bug**: You can't open files on Android before the `gio` initialize. So, you can't open
any file on `func init(){}`, for instance. The best idea is to open files on the first `system.StageEvent`.

------

### Why not use os.Open?

Because you can't, in some OSes. It's not supported on WebAssembly. If WASM isn't a concern,
you can use `os.Open`, just fine. (: