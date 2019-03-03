#!/bin/bash
set -ex
upx dist/darwin_386/xterrafile
upx dist/linux_386/xterrafile
upx dist/darwin_amd64/xterrafile
upx dist/linux_amd64/xterrafile
