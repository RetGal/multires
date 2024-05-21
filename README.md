## multires

Creates PNG images in different resolutions from a given directory containing one or more SVG file/s.

Example usage: `./multires images/`

Creates a subdirectory for each resolution, each containing the images in the corresponding resolution:

```
images/example.svg
images/exampleToo.svg
imgaes/100/example.png
imgaes/100/exampleToo.png
images/125/example.png
images/125/exampleToo.png
...
```

The desired scales and dimensions have to be adjusted in the resolution map @line 18
