## Thoughts

Go is perfect for creating packages. Just write a `.go` file, publish, and done! But making a concurrency package for the first time comes with hidden problems, upon which I stumbled during the building of this project.

Making a package for concurrency patterns is not straightforward - but not for the reasons you might expect. When I talked in Go developer communities, some of the responses to my idea were "why would I need a package? I'll just write a few lines myself". At the time, I thought "if there is a recognizable pattern, it **ought** to be abstracted for convenience!"

So as not to forget the lessons from this humble experience, I've decided to write this small text (it's mostly for myself) - and here's what I found about writing a Go package (a concurrency package anyway)

#### Abstracting concurrency

But Go has been built for concurrency, and its concurrency primitives are very easy to use and compose into patterns. When I was building the package I noticed that, with practice, I could just write my own simple worker pool from memory. What's more, having one's own implementation opens up opportunities to tweak behavior according to needs.
How and when do you need to communicate between goroutines? How do you pass the work and in which way? How do you handle errors? All those questions must be taken into account when creating an abstraction. A generalized API hides these details, often limiting optimization and control. Unsurprisingly, many Gophers are skeptical of that trade-off.

Go is easy to build with, but not abstract with. For instance, a glaring issue - Go doesn’t allow methods to introduce their own type parameters, which makes certain pipeline-style abstractions awkward. Without this, a pipeline cannot easily contain a change between types _within itself_. Divide the pipe into several parts which take and return channels each time (like I've done in this package) - workaround found! But is there a reason to having a pipe abstraction at all at this point? Maybe. But also maybe not.

#### What's next

Despite these drawbacks, there is a place for concurrency pattern abstractions. Perhaps a simple general implementation is all you need? Or maybe you're starting out and want to leverage Go's concurrency without getting too far into the weeds? Although I'd say that in Go's case, just using a concurrency abstraction is no substitute for actually knowing how those abstractions work, and being able to make them yourself on demand.

Overall, despite the limited use for "packaged" concurrency (along with having other solutions like [`conc`](https://github.com/sourcegraph/conc) or [`ants`](https://github.com/panjf2000/ants)), I can't say building this was a waste of time. Jumping in headfirst into a well-designed but difficult core feature of a language is a great way to learn and practice. Highly recommend!

- Kirill
