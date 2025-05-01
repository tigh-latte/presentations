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

lookatmego

<!-- stop -->

## Background

When your usecase is so niche that other tools just don't cut it, sometimes you feel like making your own.

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

However, it not supporting user written plugins did irk me a bit, so I thought "Golang has a plugin system, why not give it a go?".

And that is what I'm using here today, `lookatmego`.

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

- [x] Ships as a standalone binary

`go build` and that's it.

---

# Implementation

- [x] Presentation is written in markdown

Let me show you the source of this presentation.

<!-- stop -->

We use `goldmark` for parsing, more on that later.

---


# Implementation

- [x] Presentation is written in markdown

## The UI

The TUI is written using exclusively charm bracelet software, in particular:

- `bubbletea` for the UI and event handling
- `lipgloss` for styling
- `glamour` for rendering

---

# Implementation

- [x] Presentation is written in markdown

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

# Implementation

- [x] Presentation is written in markdown

## Lipgloss

Think of lipgloss as providing stylesheets but for text, you can position text, colour it, wrap it in a border, etc

<!-- stop -->

### Demo

---

# Implementation

- [x] Presentation is written in markdown

## Glamour

Automatic markdown rendering.

---

# Implementation

- [x] Presentation is written in markdown

## Goldmark

So we can render a UI, but this thing needs features! If we're looking a markdown file to declare functionality, the we're going to need to parse it. For that we can use `goldmark`.

<!-- stop -->

Using that, we just implement a simple parser that builds some state:

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

This allows us to identify and act on content, such as:

- A render stop: (`*ast.RawHTML`) `<!-- stop -->`<!-- stop -->
- End of a page: (`*ast.ThemeaticBreak`) `---` <!-- stop -->
- A widget: (`*ast.RawHTML` or `*ast.FencedCodeBlock`)

<!-- stop -->

`<!-- plugin:greeter name=john-->`

~~~plugin:not_markdown
text: |-
 ```plugin:greeter
 name: john
 ```
~~~

---

# Implementation

- [x] Supports gradually revealing a slide

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

If you imagine that, using goldmark to traverse the markdown document

Every time we hit any text, we add it to the current `Segment`.

Every time we hit a `<!-- stop -->` we start a new `Segment`.

Every time we hit a new `---` we start a new `Segment` and a new `Page`.

After that the actual meat of this wee program isn't all that hard or interesting. <!-- stop -->

Instead, the interesting part comes from a couple of the desired features:

- [x] Supports widgets

- [x] Supports user written plugins

---

# Implementation

- [x] Supports user written plugins

## Go's plugin system

If you didn't know, golang supports plugins by way of:

```go
import "plugin"
```

You define an interface that you want a user's plugin to implement, and they just adhere to that.

<!-- stop -->

So imagine a simple interface:

```go
type Plugin interface {
	Render(input []byte) string
}
```

Imagine you wanted a plugin that could take a file path and return its text, so that your presentation can just reference demo code without having to update it in two places, and it takes as input:

```yaml
file: demo_code.go
lang: go
lines:
  start: 5
```

---

# Implementation

- [x] Supports user written plugins

Well, all you would need to do is:

```go
import "github.com/tigh-latte/lookatmego"

func New() lookatmego.Plugin {
	return &filePlugin{}
}

type filePlugin struct{}

type fileArgs struct {
	Path  string `yaml:"path"`
	Lang  string `yaml:"lang"`
	Lines struct {
		Start int `yaml:"start"`
		End   int `yaml:"end"`
	} `yaml:"lines"`
}


func (f *filePlugin) Render(input []byte) string {
	var args fileArgs
	err := yaml.Unmarshal(f.Args, &args)
	// handle err


	if args.Path == "" {}// handle

	bb, err := readFile(args.Path, args.Lines.Start, args.Lines.End)
	// handle err

	return "\n```" + args.Lang + "\n" + string(bb) + "\n```\n"
}
```

And then compile:

```sh
go build -buildmode=plugin -o fileplugin.so
```

---

# Implementation

- [x] Supports user written plugins

Then all I need to do is:

```go
func ExecPlugin(input []string) string {
	plug, err := plugin.Open("fileplugin.so")
	// handle err

	builder, err := plug.Lookup("New")
	// handle err

	fn, ok := builder.(func() lookatmego.Plugin)
	// handle not ok

	return fn().Render(input)
}
```

<!-- stop -->

Pretty cool!

## Demo

---

# Implementation

- [x] Supports user written plugins

# However

I ended up not using go plugins. The constant compilation of the plugin was annoying. `.so` files trigger a lot of `WAF` rules, but, worst of all:

- If a plugin is compiled a different version of go, then it can't loaded.
- If a two plugins import competing version of the same package, then they can't be loaded.

This meant users would have to run the same go version as me, and the same version of all my deps.<!-- stop --> Unacceptable.

---

# Implementation

- [x] Supports user written plugins

# A different approach...

<!-- stop -->

So at this point I had decided that it would be easier to just implement any wee one-off features I wanted into the core application, than it would be write and maintain a golang plugin for them.

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
- Used as a user plugin language for a lot of software (`MySQL Workbench`, `mpv`, `smocker`).
- It's a config language for a lot of software (`neovim`, `wezterm`, `awesomewm`).
- Some games (such as `Garry's Mod`, `Roblox`, `Hades 2`) have some-to-extensive portions of their code written in lua allowing for easy modding.
- and the [list goes on](https://en.wikipedia.org/wiki/List_of_applications_using_Lua)

<!-- stop -->

There are many packages offering to embed lua into a go program, but the one I chose is [gopher-lua](https://github.com/yuin/gopher-lua).

---

# Demo

Let's write a stupid plugin called `showcase`.

It's defined in the markdown doc using a html comment: `<!-- plugin:showcase name=wow -->`

<!-- stop -->

<!-- plugin:showcase name=wow -->

---

# That is the most impressive thing I've seen in my life!!

<!-- stop -->

Thanks!

---

# Explain!

Ok!!!

---

# Embedding lua

Embedding lua is so simple, you set your PATH, declare a state, and execute code against that state.

## Initing your lua state

```go
lua.LuaPathDefault = strings.Join(
	[]string{
		"./lua/?.lua",
		"./lua/?/init.lua",
		lua.LuaPathDefault,
    },
)

L := lua.NewState(lua.Options{
	CallStackSize:       120,
	MinimizeStackMemory: true,
})
```

---

## Calling your lua plugin

```go
// return require 'showcase'
fn, err := L.LoadString("return require'" + name + "'")
// handle err
```
<!-- stop -->
```go
err = L.CallByParam(lua.P{
	Fn:      fn,
	NRet:    1,
	Protect: true,
})
// handle err
```
<!-- stop -->
```go
plugin := L.ToTable(1)
exec, ok := plugin.RawGetString("plugin").(*lua.LFunction)
// check if ok
isMarkdown, ok := plugin.RawGetString("is_markdown").(lua.LBool)
if !ok {
    isMarkdown = lua.LTrue
}
// handle

// turn the input input params
params := // build params from yaml
```

```go
err = L.CallByParam(lua.P{
	Fn:      exec,
	NRet:    1,
	Protect: true,
}, params)
// handle err

result := L.Get(-1).String() // the string returned 'hello everyone'
```

---

# But wait there's more

Lua is a full on language, so there is more fun to be had than just returning string. Do you remember the big reveal plugin?

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
<!-- stop -->

It is defined using a code fence:

~~~plugin:not_markdown
text: |-
 ```plugin:reveal
 text: |-
   ~ _     _   _    _
   ~| |   | | | |  / \
   ~| |   | | | | / _ \
   ~| |___| |_| |/ ___ \
   ~|_____|\___//_/   \_\
 slide_length: 75
 ```
~~~

<!-- stop -->

### The reveal plugin

To get this to work, we need lua to communicate with go.

Silly algo is:

1. Start with a string full of whitespace, equal in length to the input string.
1. Sleep for `slide_length`.
1. Reveal a random new character.

<!-- stop -->

Sounds easy, but lua doesn't have a `sleep` function, so what can we do?

---

# Calling Go from Lua

Calling go from lua is so simple. We a define a function, and add it to our lua state (`L`).

```go
func sleep(L *lua.LState) int {
	millis := time.Duration(L.ToInt(1))
	time.Sleep(millis * time.Millisecond)
	return 0
}

L.SetGlobal("sleep", &lua.LFunction{
	IsG:       true,
	GFunction: sleep,
})
```
<!-- stop -->

Now, from within our plugins, we can call `sleep`:

```lua
return {
	plugin = function(args)
		while true do
			sleep(args.slide_length)
			-- reveal new character
			-- break when all revealed
		end
		return "???"
	end,
}
```

---

# Sending to go channels in Lua

We can add a `chan` to our lua state:

```go
ch := make(chan lua.LValue)
L.SetGlobal("ch", lua.LChannel(ch))

go func() {
	for ev := range ch {
		fmt.Println(ev.String())
	}
}()
```

<!-- stop -->

And we send content to this channel in lua, using `ch:send(str)`.

---

# A quick `fib` showcase

<!-- stop -->

<!-- plugin:fib -->

---

# But wait there's EVEN MORE

<!-- plugin:echo -->

<!-- stop -->

## Accepting input

In our plugin, we define an `on_key` function, and while focused, all bubbletea `tea.KeyMsg` inputs are just forwarded to this function for Lua to handle.

```lua
return {
	plugin = function() return "> " end,
	on_key = function(input)
		-- do something with input
		return input
	end,
	is_markdown = false,
}
```

---

# Let's build a plugin that actually does something

<!-- plugin:weather prompt=> -->

---

# And that's us folks

<!-- stop -->

...unless

<!-- stop -->

## What I can show you today

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [x] Supports widgets

- [x] Ships as a standalone binary

- [ ] Image rendering

- [x] Supports user written plugins

---

# And that's us folks

...unless

## What I can show you today

- [x] Presentation is written in markdown

- [x] Supports gradually revealing a slide

- [ ] Image rendering

---

# And that's us folks

...unless

## What I can show you today

- [ ] Image rendering

<!-- stop -->

```plugin:image
path: ./gopher.png
```
