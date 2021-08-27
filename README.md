# Chainconfig

This is an example of a design pattern for passing configuration through a system,
but without breaking API contracts when new configurations are added (and, to an
extent, when old configurations are removed). I learned it from Lee Boynton back
in the WebMD days. While we did it in Java then (and he used Python before that), this implementation is in Go.

The idea is that you use a linked list to pass configuration around. The keys are
well known, though not necessarily discoverable. The values are unstructured.
This puts the impetus on the developer to "assert" that the values are appropriate
in their specific context.

THIS IS JUST AN EXAMPLE, AND IS NOT PRODUCTION/LIBRARY-GRADE CODE.

## Strengths

The benefit is that one standard configuration type can be passed throughout the codebase, but without stipulating what is in that configuration.

- It makes testing easy -- you can create configs with _just_ the things needed for the test, and in so doing flag any case where a new config dependency is added.
- It is easier to work with than a mega-struct. You don't need a specialized builder or giant struct declarations.
- Config can be changed locally for a certain case, but not globally, simply by doing something like `localConfig = globalConfig.add(MyOverride, "new value")` (In that case, you get the global config with just `MyOverride` changed for the local config)

In more advanced cases (not modeled in the code here), you can even add methods to
detect if there are multiple versions of a param. E.g. you can check to see if a
parameter has been overridden. This was a nice debugging feature, and also allowed
us to unwind config stuff. (Keep in mind, this was used on a HUGE sprawling codebase
with dozens and dozens of engineers all working on it. On occasion, we had elaborate
hacks.)

## Weaknesses

This way of doing things had some _definite_ drawbacks.

1. It is essentially a way to work around the type system
    - It is always incumbent on the developer _using_ the config to know the type of the data in the value field.
    - Convention enforcement becomes important, because changing a type could break the
      code in wild and difficult-to-debug ways
    - One time, an int was changed to a float on insertion, but all over the code we had
      places that assumed int (and returned a default int value). It caused a major headache
    - This can be partially mitigated by adding some type-sensitivity to the config
      object itself, and forcing a declaration on type during the insertion phase
    - I haven't looked, but using the `reflect` package might be helpful
2. Consts should not be removed. And that's not a drawback. If a key is no longer set, everyone will just end up using the default. This is a design of the system, but can end up being surprising.

It is possible to use strings as keys instead of consts. But they are definitely harder 
to discover. Essentially, you have to audit all calls to `New()` and `Add()` to find 
what params are available. It's also easy to accidentally override one because you 
didn't realize, for example, that `title` was already used somewhere.