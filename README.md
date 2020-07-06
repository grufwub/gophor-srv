This is a re-write of https://github.com/grufwub/gophor with the following in
mind:

- Separate much core functionality into `core/`

- With separated core functionality, using much of same codebase write a gemini
  protocol client (see https://gemini.circumlunar.space)

- Follow Unix philosophies on simplicity: do one thing, do it well. As such,
  features like reverse proxying, proxy protocol have been dropped (there may be
  others that I'm forgetting)

- In time add unit tests, and performance tests