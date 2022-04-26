# OrangeWire
A secure peer-to-peer file sharing application that utilizes encryption and obfuscation to ensure that files and users are not traceable.

To demo the program you will need to setup your flags correctly:
run "go run .\OrangeWire.go -first=true -serverport=2000" in the root directory.
Then you can run the test.bat script to connect 20 peers to your test environment.

There is already a default file listed within the upload folder for the root peer (the peer run with the -first flag).
Uploading a file simply requires the user to type the word "upload" and the program will prompt for file name and encryption password.
It may seem as though the program is not doing anything for a moment as it gathers the list of peers that are available for uploading and parses through them.
The file is finished uploading when the console ">" pops up once again prompting for user input.

Searching a file simply requires the user to type the word "search" and the program will prompt for file name and return all available files that match
that name. The user simply selects which file they want and the program handles the downloading automatically.
The file is finished downloading when the console ">" pops up once again prompting for user input.

The network architecure for this program is as follows.
Searching for a file:
![Untitled Diagram](https://user-images.githubusercontent.com/38273962/165395370-160c7ebe-8130-4631-9aa7-47ee7bf2b870.jpg)
Downloading a file:
![Download](https://user-images.githubusercontent.com/38273962/165395427-49283ff0-794e-4792-9bf3-ebbb125a65ab.jpg)
Uploading a file:
