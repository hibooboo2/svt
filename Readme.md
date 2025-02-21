# What do

# Deps
Make sure you have go setup and that go bin is in your path.

run

```
go install github.com/hibooboo2/svt@latest
```

# Add this line to your bashrc
```
. $THIS_FOLDER/aliases.sh
```

# Usage
```
gntag
```
```
gutag
```
```
gptag
```
Each of these will checkout development/ staging / production
respectively and then pull and if the commit is not tagged it will tag with the next tag in sequence.

If you have git pre push / hooks and you want to skip verification you can add:

--no-verify