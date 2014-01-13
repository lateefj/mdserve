mdserve
=======

The idea is a single binary be able to edit markdown files and see what the html output is.


mdserve defaults to README.md and port 7070

```sh
mdserve 
```

mdserve can take a markdown file as its first parametert

```sh
mdserve NOTES.md
```

mdserve also change the port number by simply adding -port=7077

```sh
mdserve -port=7077 NOTES.md
```
