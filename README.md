# markdown to pukiwiki converter

Convert markdown to pukiwiki notation.

## pukiWiki notaion coverage

**Roadmap**

- [x] Headers
- [ ] List
- [ ] Codeblock
- [ ] Bold
- [ ] Link

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
