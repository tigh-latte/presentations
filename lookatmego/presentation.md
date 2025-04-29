---
theme: ./theme.json
title: 'lookatmego: a story'
author: Tighearnán Carroll
styles:
  table:
      column_spacing: 3
extensions:
  - '*'
---

# Fáilte!

lookatmego: a story

<!-- stop -->

## Background

When your usecase is so niche that other tools just don't cut it, sometimes you have to make your own.

---

# In the beginning....

This is the fourth talk I've given. At the end of each one I would occassionally have someone mention the presentation software that I use. <!-- stop -->

- First talk: Someone asked me what am I using? <!-- stop -->

- Second talk: Someone said "that presentation software looked cool" <!-- stop -->

- Third talk: No one mentioned the software.

<!-- stop -->

Given the following formula:

1. Someone asked what was used: +1 point
1. Some said "that's cool": +0.5 points
1. No one cared : +0 points

Interested in terminal slideshows scores 1.5 points out of total of 3, so 50% of the time there is interest!! <!-- stop -->**(Investors hit me up)**

---

# So what do I use?

The truth is, the software I used has changed between presentations: <!-- stop -->

1. At first I used software called `lookatme`, it is written in python, and has/does everything you could want:

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [x] Supports widgets, one of which was a live terminal!

- [x] Supports user written plugins

- [ ] Ships as a standalone binary

<!--stop-->

## Why don't you use it now?

Because, it hasn't been updated in 2 years, and the python ecosystem has outpaced it.

---

# So what do I use?

The truth is, the software I used has changed between presentations:

2. After `lookatme` fell victim to the inevitable heat death of a scripting language's ecosystem, I looked for an alternative that ideally wouldn't be shipped as code you run, but instead a binary, so I came across `slides`. However:
<!-- stop -->

- [x] Presentation is written in markdown

- [ ] Supports gradually revealing a slide

- [x] Supports widgets,

- [ ] one of which was a live terminal!

- [ ] Supports user written plugins

- [x] Ships as a standalone binary
<!-- stop -->
- [x] Can be hosted on a server and rendered over ssh

<!-- stop -->

So while it was a good application, it was missing too many features I really liked to use.

---

# So what do I use?

The truth is, the software I used has changed between presentations:

3. Next I came across `presenterm`

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [x] Supports widgets

- [ ] one of which was a live terminal!

- [x] but it can execute arbitrary bash commands!

- [x] Shipped as a standalone binary
<!-- stop -->
- [x] Image rendering
<!-- stop -->
- [x] Gifs
<!-- stop -->
- [x] Columns
<!-- stop -->
- [x] Mermaid diagrams
<!-- stop -->
- [ ] Supports user written plugins

<!-- stop -->

This is what I used for my last talk, and it is honestly great software, if you're looking to give a terminal themed talk, I would recommend that you look at it.

However, it not supporting user written plugins did urk me a bit, so I thought "Golang has a plugin system, why not give it a go?"

---

## The minimum features I wanted

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [x] Supports widgets

- [x] Ships as a standalone binary

- [x] Image rendering

- [x] Supports user written plugins

---

## What I can show you today

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [x] Supports widgets

- [x] Ships as a standalone binary

- [ ] Image rendering

- [x] Supports user written plugins

---

# Just a quick note

I'm going to quite quickly through most of these bullet points because they just aren't all that interesting.

---

# Implementation

- [x] Presentation is written in markdown

Let me show you the source of this presentation.

---

# Implementation

- [x] Supports gradually revealing a slide

A presentation is made up of slides (`[]Slide`), and slides are made up of sections (`[]string`).

A library `goldmark` is used to parse the markdown file into a syntax tree. The syntax tree is iterated, building the slide.

We add content to the current working section. Once a `<!-- stop -->` is found, we create a new current working section. Once a `---` is found, the current slide is finalised and a new one is created.

---

# Implementation

- [x] Supports widgets

As well as discovering `<!-- stop -->`, helpers are usable

---

# Technology used

## The UI

The TUI is written using exclusively charm bracelet software, in particular:

- `bubbletea` for the UI and event handling
- `lipgloss` for styling
- `glamour` for rendering

---

# Technology used

## BubbleTea

Bubbletea is a MVC framework for TUIs which runs on an event loop. To work with it you just need to implement one interface:

```go
package tea

type Model interface {
	// Init is the first function that will be called. It returns an optional
	// initial command. To not perform an initial command return nil.
	Init() Cmd

	// Update is called when a message is received. Use it to inspect messages
	// and, in response, update the model and/or send a command.
	Update(Msg) (Model, Cmd)

	// View renders the program's UI, which is just a string. The view is
	// rendered after every Update.
	View() string
}
```

<!-- stop -->

### Demo

---

# Technology used

## Lipgloss

Think of lipgloss as providing stylesheets but for text, you can position text, colour it, wrap it in a border, etc

<!-- stop -->

### Demo

---

# Technology used

## Glamour

Automatic markdown rendering.

---

# Technology used

So I can render a UI, but this thing needs features! If I'm looking a markdown file to declare functionality, the we're going to need to parse it. For that I use `goldmark`.

<!-- stop -->

After that, we just implement a simple parser that builds some state:

```go
func (p *parser) parse(...) (ast.Node, error) {
    for ; v != nil; v = v.NextSibling() {
        switch n := v.(type) {
            case *ast.Paragraph:
                // paragraph specific pre handling
                p.parse(...)
                // paragraph specific post handling
            case *ast.Heading:
                // heading specific pre handling
                p.parse(...)
                // heading specific post handling
            // etc etc
        }
    }
}
```

<!-- stop -->

This allows me to spot features, such as:

- A render stop: `<!-- stop -->`<!-- stop -->
- End of a page: `---` <!-- stop -->
- A widget:

```html
<!-- plugin:mycoolplugin arg1=hello arg2=world -->
```

or

~~~plugin:raw
text: |-
  ~   ```plugin:mycoolplugin
  ~   arg1: hello
  ~   arg2: world
  ~   parent_arg:
  ~     sub_arg: holyhell
  ~   ```
~~~

---

# Technology used

## Piecing it together

Consider the following struct:

```go
type Presentation struct {
    Pages []Page
}

type Page struct {
    Segments []Segment
}
```

If you imagine that, using goldmark to travese the mardown document

Every time I hit a new `---` I start building a new `Page`

Every time I hit a `<!-- stop -->` I starting a new segment, then the actual meat of this wee program isn't all that hard or interesting. <!-- stop -->

Instead, the interesting part comes with one of my desired features:

- [x] Supports user written plugins

---

# Technology used

do a whole section on go plugins

---

# Technology used

So at this point I had decided that it would be easier to just implement any wee one-off features I wanted into the core application, than it would be write a golang plugin for them.

So, this project was parked for a few weeks until Boxing Day last year, I was sitting on my sofa and then the solution popped into my head.

<!-- stop -->

```plugin:reveal
text: |-
 ~ _     _   _    _
 ~| |   | | | |  / \
 ~| |   | | | | / _ \
 ~| |___| |_| |/ ___ \
 ~|_____|\___//_/   \_\
slide_length: 75
```


---

# Technology used

## Lua

Lua is a fast, _highly_ embeddable programming language with (similar to go) a very small built-in feature set. You are given the basics, and you're free to fire on.

Given it's speed, ease of embedding, and tiny size, you find it in more places than you would think:

- You can execute lua in `Redis` commands.
- You can configure `nginx` using lua via a plugin.
- `.rpm` files support lua scripting during installation.
- `MySQL Workbench` for user addons.
- `ScyllaDB` allows server-side functions to be written in Lua.
- It's a config language for a lot of software (`neovim`, `wezterm`, `awesomewm`).
- Some games (such as `Garry's Mod`, `Roblox`, `Hades 2`) have some-to-extensive portions of their code written in lua allowing for easy modding.
- and the [list goes on](https://en.wikipedia.org/wiki/List_of_applications_using_Lua)

<!-- stop -->

There are many packages offering to embed lua into a go program, but the one I chose is [gopher-lua](https://github.com/yuin/gopher-lua).

---

# Demo

So, I figured I would first show you a simple plugin being written, then explain how to do all of this.

I have defined a (yet to be written) plugin called `showcase`:

```html
<!-- plugin:showcase name=wow -->
```

<!-- stop -->

```plugin:showcase
name: |-
    a wise man once said
    > this is a quote
```
---

# That is the most impressive thing I've seen in my life!!

Thanks!

---

# Explain!

Ok!!!

---

# Embedding lua

---

# But wait there's more

<!-- plugin:echo -->

---

# Let's see what the weather is

<!-- plugin:query_weather prompt=> -->
