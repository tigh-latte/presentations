---
theme: ./theme.json
title: Practicing Unsafe Stacks
author: Tighearnán Carroll
extensions:
  - terminal
  - file_loader
---
# Fáilte!

Practicing Unsafe Stacks.

<!-- stop -->

## Background

How using `unsafe` can reduce allocations (with hilarious consequences!)

<!-- stop -->

## Disclaimers

- I am by no means an expert in writing low allocation code.
<!-- stop -->

- I am by no means an expert in memory safety.
<!-- stop -->

- I am by no means recommending or advocating you do this.

---

# `unsafe`

Unsafe allows us to bypass certain golang memory safety operations by giving us a (limited) number of APIs to interface directly with the host machine's memory.


