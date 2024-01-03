#!/bin/bash

go build -o infogen
tar czf infogen.tgz infogen coords.json form.pdf README.md
rm infogen
