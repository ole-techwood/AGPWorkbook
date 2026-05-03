---
name: golang-idioms-best-practices
description: "Apply Go idioms and best practices for design, errors, naming, interfaces, and maintainability. Use when implementing or reviewing Go changes, refactoring handlers/services, or improving code quality before merge."
argument-hint: "Target area, goal, and constraints. Example: internal/analysis add ftp scheme with idiomatic design"
---

# Golang Idioms and Best Practices

## Principles of Code Readability and Maintainability

1. **Clarity and Simplicity:** Write code that is easy to understand and straightforward. Avoid overly complex or convoluted solutions, and favor simplicity and clarity over cleverness.
2. **Consistency:** Follow consistent naming conventions, coding styles, and formatting guidelines throughout the codebase. Consistency improves readability and reduces cognitive overhead when navigating and understanding the code.
3. **Modularity and Encapsulation:** Decompose complex systems into smaller, modular components with well-defined interfaces. Encapsulate implementation details and minimize dependencies between modules to improve code maintainability and reusability.
4. **Documentation:** Write clear, concise, and meaningful comments and documentation to explain the purpose, behavior, and usage of code. Document public APIs, interfaces, and complex algorithms to facilitate understanding and usage by other developers.
5. **Multi-Option Values via Maps:** When one variable can contain two or more expected options, do not chain direct equality or inequality checks (`==`, `!=`) against each option. Prefer map lookups where keys are expected values and values are booleans or associated metadata, similar to service routing maps.
6. **PascalCase for Named Types:** Use PascalCase naming for interfaces, types, and structs. Avoid lower camelCase or prefixed forms for these declarations.

Reference code and usage examples: [examples](./references/examples.md)

## Interfaces and Abstraction

Use interfaces to define contracts between different parts of your code. This allows for flexibility and enables you to swap implementations easily without changing the interface.

### When should you use interfaces?

**1. When you need multiple implementations**

If your application should work with different databases (for example PostgreSQL and MongoDB) or different message delivery mechanisms.

**2. For unit testing**

This is the most common case. If your function calls an external API or deletes files from disk, you do not want to do that in every test. An interface lets you substitute a mock or stub implementation.

**3. For decoupling modules**

If you are building a library for other developers, provide an interface so they can plug in their own logic that fits your framework.

### When should you NOT use interfaces?

**1. If you only have one implementation and it is unlikely to change**

If you are writing a small internal script or helper function that formats a string, an interface usually adds unnecessary complexity.

**2. For "future flexibility" only (YAGNI)**

Follow the YAGNI principle (You Ain't Gonna Need It). Do not create interfaces just in case. In Go, introducing an interface later via refactoring is usually straightforward.

**3. If it harms readability**

When navigating code in an IDE, you usually want to jump directly to behavior. Excessive interfaces can force extra navigation through method lists before finding the actual implementation.
