#! /bin/bash

echo -n "Building... "
go build -o peek
echo "done!"

echo -n "Copying executable to ~/.local/bin... "
cp ./peek ~/.local/bin/
echo "done!"
