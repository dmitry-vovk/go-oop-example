go-oop-example
==============

Run it by ```go run OOP.go```. Navigate with your browser to hostname where you have run it, and try to download a file.

E.g. go to http://localhost/file/big_file.iso.

While the file is downloaded, open new tab and navigate to http://localhost/sessions. You will see something like this:

```js
{
  "sessions": [
    {
      "file": "big_file.iso",
      "start": 1405712431,
      "ip": "127.0.0.1",
      "bytes": 1518698496
    }
  ]
}
```