# PizzaCLI

## Overview

PizzaCLI is a terminal-based pizza ordering system (mock orders). This project will help you understand how to work with memory in Go efficiently.

## Requirements

### 1. Core

- Hardcode a static menu of 7 pizzas.
- Support adding pizzas to a current order with quantities.
- Display the current order with itemized prices and quantities.
- Calculate and display the order total.
- Mock checkout that prints a summary and clears the order.

### 2. Memory Allocation

Define three main structs:

```go
type Pizza struct {
    ID    int
    Name  string
    Price float64
}

type OrderItem struct {
    Pizza    *Pizza     // pointer
    Quantity int
}

type Order struct {
    Items []OrderItem
}
```

**Stack allocation (must NOT escape):**

- Local `Pizza` variable when looking up a menu item (only used temporarily inside the function)
- Local variables used during total calculation (subtotal, loop variables, total accumulator – when not stored anywhere)
- Temporary string builders or format variables used only for printing

**Heap allocation (must escape):**

- The slice returned by the `getMenu` function (backing array escapes the function)
- The `Order.Items` slice itself (grows over time, lives as long as the order exists)
- The `Order` struct itself – allocated on the heap via `&Order{}` or `new(Order)` because it is used across multiple command calls
- Individual `Pizza` structs inside the menu slice (they live for the lifetime of the program)

## **3. Pointers and Memory Efficiency**

- Change `OrderItem.Pizza` from value to `Pizza`
- When adding an item, store a pointer to the existing `Pizza` from the menu instead of copying the whole struct

## **4. Common Caching Strategies**

Add a very simple in-memory cache with TTL to store the most recently used pizza lookups. Do this purely for education purposes. Of course, because we have a hardcoded menu, caching the single pizza won’t improve the performance, but if had to make the API call every time we’re finding the pizza, it would make a difference:

- Create a small struct-based cache that stores `Pizza` by normalized name (lowercase)
- Set a fixed TTL of 30 seconds for each cached entry
- When user runs `add <name>`, first check the cache; if found and not expired → use cached pointer
- If not found or expired → lookup in menu, store in cache with current time, then use it
- Cache is a simple map in package `menu` (global variable or passed around)

## **CLI Interface Specs**

- `menu` — Shows numbered list: `ID | Name | Price`
- `add <name> [quantity]` — Adds pizza by name (case insensitive), quantity defaults to 1
  - Example: `add margherita 2` or `add pepperoni`
  - If name not found → error message
- `order` — Shows current order:
  - Each line: `quantity × name @ price = item total`
  - Last line: `Total: $xx.xx`
- `checkout` — Prints full order summary + “Order placed! Thank you.” message, then clears the current order
- `help` — Lists all commands
- `exit` — Quit the interactive CLI

## **Success Criteria**

- All core commands work correctly
- Program does not panic on valid usage
- Adding an item stores only a pointer to `Pizza` → no repeated copying of `Pizza` struct
- Escape analysis (`go run -gcflags="-m" .`) shows clear separation:
  - **Expected stack (must NOT escape):**
    - temporary variables inside lookup function (loop index, lowered name string)
    - local `item` in `AddItem` before append
    - local variables inside `Total()` and printing functions
  - **Expected heap (must escape):**
    - `menu` slice and its `Pizza` elements
    - `Order.Items` slice backing array
    - `OrderItem` structs once appended (because they are stored in slice)
    - `Order` struct itself (because pointer escapes main loop)
    - `pizzaCache` map and its `cachedPizza` values
- The pointer in `OrderItem.Pizza` points to the original menu item → adding many of the same pizza creates no extra `Pizza` allocations
- After adding the same pizza twice quickly (< 30 s) → second lookup should hit cache (you can verify by adding debug print like "cache hit" / "cache miss")
- After waiting > 30 s and adding again → should miss cache and re-store

## **Constraints**

- Standard library only (no external packages)
- Use only `time`, `strings`, `fmt`, `bufio`, `os`, `strconv` etc.
- Modular structure:
  - `main.go` → CLI loop, command parsing
  - `menu.go` → menu data
  - `order.go` → `Order` type + methods (`AddItem`, `Total`, `Clear`, printing helpers)
