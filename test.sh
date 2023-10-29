#!/bin/bash

ab -n 3000 -c 100 -g result.log -r http://localhost:80/
