## cf - codeforces.com console client

Features:
- possibility to submit code
- and see the result.

### Installation
1. install go compiler, as per https://golang.org/dl
2. `go get -u github.com/Komosa/cf`

If you don't like to compile it yourself, pre-built binaries are available here: https://github.com/Komosa/cf_binaries

### Usage
**Note**: to achieve best performance use `cf login` and `cf con CONTESTNUMBERHERE` before the rated contest starts.

Examples (commands):
- `cf login exampleuser` will try log in as _exampleuser_, asking for password if necessary;
- `cf login` will try log in as previously logged user, asking for password if necessary;
- `cf submit x.c -prob=555a -lang=10` submits file _x.c_ as solution for problem _555A_, using compiler number _10_;
- `cf submit x.c -prob=555a` submits file _x.c_ as solution for problem _555A_, using default compiler for _c_;
- `cf submit a.c -prob=555` submits file _a.c_ as solution for problem _555A_, using default compiler for _c_;
- `cf con 555` from now all subsequent _submit_ ops will use _555_ as **default** contest;
- `cf submit x.c -prob=a` submits file _x.c_ as solution for problem _555A_, using default compiler for _c_;
- `cf submit a.c` submits file _a.c_ as solution for problem _555A_, using default compiler for _c_;
- `cf submit a` will check if it is obvious which file do you want to submit and submit it using default (or specified) language as solution for problem _A_ in default (or specified) contest.

Configuration will be in `$HOME/config/.cf/conf` file, cookies will be in separate files in the same directory as conf. Feel free to edit/delete/share this files at own risk.

Limitations:
- `=` character in _handle_ is not supported (probably also by the site).
