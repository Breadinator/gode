### overview
gode updates vscode workspaces saved in `%GOPATH%/src` directories to include all go modules in `%GOPATH%/src` and then, if only one workspace was found, it launches vscode using that workspace.

### raison d'être
vscode doesn't seem to like just opening the `%GOPATH%/src` directory, it messes with the imports, at least for me, and i'm not about to spend money on goland

### limitations
* i'm pretty sure any workspaces that are more complicated than just list of folders will lose all that data because it shouldn't be parsed properly
* it might only function as intended in windows?