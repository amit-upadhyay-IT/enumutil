# enumutil

> A library for fetching enum keys and values.

Generally in Go, an enum is represented using constant variables of some primitive type, (they can't be non-primitive, however, their name can be aliased, eg: `type CustomType int`), and unlike other languages, Go doesn't provide us the options like `enum.Values()` for fetching out the values of enums.

To solve this problem you might wanna store the values in a `list` or a `map`, but then in later when the project grows it might
become difficult to maintain coz, while updating enum you will need to do changes at two different places: 1) enum itself 2) place where you are storing the values

But you can avoid this maintenance task by using this library.

**Dear visitor, It seemed to me that this code snippet is complex to understand, please let me know your opinion if you are looking at it.**
