# Tasker
just a tasker

## 1.install
`go install github.com/glennliao/task@latest`

## 2.create taskfile.js
```js
task("hi")
    .echo("hello")
    .cmd("ls")
```

## 3.run task
```
task hi
```