# md2pw

Convert markdown to pukiwiki notation. (markdown -> pukiwiki)

**Supported Notations**

- [x] Headers
- [x] List
- [x] Codeblock
- [x] Bold
- [x] Link
<<<<<<< Updated upstream
=======

## Install

**cURL**

```bash
curl -sSfL https://raw.githubusercontent.com/moriT958/md2pw/main/install.sh | sh
```
>>>>>>> Stashed changes

## pukiWiki notaion coverage

Supported notaitions.

### Headers

pukiwiki は 3 段階ある。

**PukiWiki**

```text
* H1
** H2
*** H3
```

**Markdown**

```markdown
# H1

## H2

### H3
```

### List

インデントは 3 Level まで対応

#### Ordered

同じ

**PukiWiki / Markdown**

```text
- list1
- list2
- list3
```

#### Unordered

**PukiWiki**

```text
+ ordered1
+ ordered2
+ ordered3
```

**Mardown**

```markdown
1. ordered1
2. ordered2
3. ordered3
```

### Code Block

**PukiWiki**

```text
  this is sample code.
  need 2 spaces
```

**Markdown**

```markdown
\`\`\`(filetype)
this is sample code.
\`\`\`
```

### Bold

**PukiWiki**

```text
''text''
```

**Markdown**

```markdown
**text**
```

### Link

**PukiWiki**

```text
[[this is link>https://example.com]]
```

**Markdown**

```markdown
[this is link](https://example.com)
```

## Development

- deps
  - Go
  - golangci-lint
  - task

- dev commands
  - Run command: `task run`
    - with args: `task run -- <args>`
  - Build: `task build`
  - Test: `task test`
  - Lint: `task lint`
  - Release check: `task release`
  - Deploy: `task deploy -- --tag <v0.0.0>`
