# embiggen
Generates BOGUS Java files from a directory of Smali files.

This lets you remotely set breakpoints in an APK you do not have source code for using Android Studio.

"A noble spirit embiggens the smallest man"

The project name is a reminder that this is stupid.  Stupid, but necessary.  If it was called smali2java or something, people would be expecting actual Java code...

# install
go get github.com/awgh/embiggen

# usage
Usage of embiggen:

  -i string
  
   Input Directory or File (\*.smali) (default "./")
      
  -o string
  
   Output Directory (default "java/")
      
# example usage

embiggen -i smali -o java

This will generate bogus Java sources in java/ for each .smali file in smali/. 
