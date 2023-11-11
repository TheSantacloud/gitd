# mgtd - Task Management CLI

**mgtd** is a CLI tool designed for managing tasks seamlessly with popular task managers like Todoist and archive managers like Obsidian. It empowers users to streamline their workflow and enhance their Getting Things Done (GTD) methodology.

## Features

- **Review Phase:** Conduct a thorough review of all tasks and efficiently purge old and irrelevant items from the task manager.
- **Purge Tasks:** Purge tasks based on a specified timespan, ensuring your task manager remains focused on current and relevant activities.

## TODO

### Features

- [ ] add archive managers
- [ ] add obsidian
- [ ] add next actions control
- [ ] add articles support
- [ ] add "search context" support ("i wanna create a webapp" -> "you wanted to try out shadcn and nextjs14")
- [ ] setup taskmanager with all the relevant labels and such
- [ ] add weekly review process
- [ ] add monthly review process
- [ ] make review processes configurable
- [ ] alfred plugin
- [ ] daily planner/journal(?)

### Code

- [ ] add viper support and generally make config better
- [X] figure out how to use client secrets better
- [ ] tests
- [X] restructure the package
- [ ] rename mod
- [X] rework to support sync requests (Todoist)
- [ ] indicative error messages and exit codes
- [ ] cache results on initial run and update them in real time
- [ ] make it the entire thing a library

### Bugs

- [ ] purge table clipping on selection (+ fix height)

## Installation

```bash
go install github.com/dormunis/mgtd
```

## Usage

### Review Phase

```bash
mgtd review
```

During the review phase, **mgtd** allows you to evaluate all tasks and remove those that are outdated or no longer relevant.

### Purge Tasks

```bash
mgtd review purge --timespan="1 month"
```

The `purge` command helps clean up your task manager by removing tasks older than the specified timespan.

- Use the `--timespan` flag to set the timespan for reviewing tasks. The default is "1 month."

## Configuration

**mgtd** utilizes a configuration file to adapt to your preferences. Ensure that your settings are correctly configured for seamless integration with your task and archive managers.

## Notes

- This CLI currently supports Todoist as the default task manager.
- Archive actions, including creating Markdown files and archiving with archive managers, are not yet implemented.

## Contributing

Feel free to contribute to the development of **mgtd** by submitting issues or pull requests on the [GitHub repository](https://github.com/dormunis/mgtd).

## License

This project is licensed under the [MIT License](LICENSE).

---

**Note:** This README provides a basic overview of the functionality and usage of **mgtd**. Make sure to check for the latest updates and documentation in the [GitHub repository](https://github.com/dormunis/mgtd)

