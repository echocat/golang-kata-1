# Golang Kata - Part 1: Library

Code kata in Golang where you have to write a library of books and magazines.

[A code kata is an exercise in programming which helps a programmer hone their skills through practice and repetition.](https://en.wikipedia.org/wiki/Kata_(programming))

* [Topic](#topic)
* [Frame conditions](#frame-conditions)
* [Tasks](#tasks)
* [Procedure](#procedure)
* [License](#license)

## Topic

You have to implement an abstracted and simple library system.

## Frame conditions

1. You have exact 2.5 hours of time - no minute longer.

   If you reach this time limit stop your work immediately.
   It is one part of the kata to respect this time limit.
   
2. There are no restrictions on how to use the provided time.
   If you want to code the entire time, take a break or a cigaret - it’s up to you.

3. This is a real world situation. You are allowed to consult the Internet, use every library you want, call a friend ...

   **BUT:** You are not allowed to do [pair programming](https://en.wikipedia.org/wiki/Pair_programming).
   **AND** If you have already done this kata before to have a look at your previous implementation.
   
4. Develop your code based on Golang 1.11+.

5. Keep the following priorities in mind while you implementing - in the mentioned order:
   1. Code quality
   2. Structure
   3. Functionality

6. Given resources:

   > **Hint:** There is a reason why there are so many books and authors in german with [umlauts](https://en.wikipedia.org/wiki/Germanic_umlaut).

   * [`authors.csv`](resources/authors.csv): Contains authors with its `email`, `firstName` and `lastName`.
   * [`books.csv`](resources/books.csv): Contains books with its `title`, `description`, one or more `authors` and an `isbn`.
   * [`magazines.csv`](resources/magazines.csv): Contains magazines with its `title`, one or more `authors`, a `publishedAt` and an `isbn`.

## Tasks

* [Main tasks](#main-tasks)
* [Optional tasks](#optional-tasks)

### Main tasks

1. Your software should read all data from the given CSV files in a meaningful structure.

2. Print out all books and magazines (could be a GUI, console, …) with all their details (with a meaningful output format).

   > **Hint**: Do not call `printBooks(...)` first and then `printMagazines(...)` ;-)

3. Find a book or magazine by its `isbn`.

4. Find all books and magazines by their `authors`’ email.

5. Print out all books and magazines with all their details sorted by `title`.
   This sort should be done for books and magazines together.

### Optional tasks

> **Hint:** Optional means optional.

1. Write Unit tests for one or more methods.

2. Implement an interactive user interface for one or more of the main tasks mentioned above.
   This could be done by a website, on the console, etc.

3. Add a book and a magazine to the data structure of your software and export it to a new CSV files.

## Procedure

1. Get the code. There are several ways for it:

   1. With fork (makes it possible to preserve your work):
      1. [Fork this repository](https://github.com/echocat/golang-kata-1/fork)
      2. Clone this fork to your computer:
         ```bash
         git clone <your github url>
         # Example: git clone https://github.com/blaubaer/golang-kata-1.git
         ```

   2. Clone this repository with local branch:
      ```bash
      git clone https://github.com/echocat/golang-kata-1.git
      git checkout -b run-<yourname>-<iteration number> 
      # Example: git checkout -b run-blaubaer-1
      ```

   3. Just download it from [here](https://github.com/echocat/golang-kata-1/archive/master.zip)

2. Open in your favorite IDE.

   > **Hint**: We recommend [JetBrains GoLand](https://www.jetbrains.com/go/).

3. Start the kata.

4. Discuss with your friends and/or colleges your solution.

5. Repeat after some days/weeks.

## FAQ

##### How to run your application?

1. By IDE:<br/>
   Just click with right mouse on the [`main.go`](main.go) class and
   then on _Run_ or _Debug_.
2. By command line:
   ```bash
   go build . && ./golang-kata-1
   ```
##### How to run your tests?

1. By IDE:<br/>
   Just click with right mouse on the root of your project tree and click on
   then on _Run all tests_ or _Debug all tests_.
2. By command line:
   ```bash
   go test -v ./...
   ```

## License

See [LICENSE](LICENSE) file.
