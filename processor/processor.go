package processor

import (
	"fmt"
	"github.com/karrick/godirwalk"
	"github.com/monochromegane/go-gitignore"
	"path/filepath"
	"runtime"
	"strings"
)

var ExtensionToLanguage = map[string]string{"scm": "Scheme", "asmx": "ASP.NET", "handlebars": "Handlebars", "less": "LESS", "csproj": "MSBuild", "hamlet": "Hamlet", "lds": "LD Script", "tcl": "TCL", "gd": "GDScript", "lsp": "Lisp", "go": "Go", "sv": "SystemVerilog", "xml": "XML", "ads": "Ada", "clj": "Clojure", "elm": "Elm", "ihex": "Intel HEX", "ede": "Emacs Dev Env", "ts": "TypeScript", "yaml": "YAML", "sml": "Standard ML (SML)", "adb": "Ada", "tsx": "TypeScript", "vbproj": "MSBuild", "pcc": "C++", "dockerfile": "Dockerfile", "yml": "YAML", "ly": "Happy", "coffee": "CoffeeScript", "rx": "Forth", "vert": "GLSL", "psl": "PSL Assertion", "md": "Markdown", "rake": "Rakefile", "ahk": "AutoHotKey", "dockerignore": "Dockerfile", "dart": "Dart", "bat": "Batch", "d": "D", "lisp": "Lisp", "h": "C Header", "cmd": "Batch", "abap": "ABAP", "c++": "C++", "p": "Prolog", "ml": "OCaml", "x": "Alex", "sitemap": "ASP.NET", "el": "Emacs Lisp", "cfm": "ColdFusion", "htm": "HTML", "cfc": "ColdFusion CFScript", "ec": "C", "tex": "TeX", "ex": "Elixir", "btm": "Batch", "ur": "Ur/Web", "targets": "MSBuild", "webinfo": "ASP.NET", "for": "FORTRAN Legacy", "scala": "Scala", "lean": "Lean", "lua": "Lua", "ceylon": "Ceylon", "gvy": "Groovy", "geom": "GLSL", "rb": "Ruby", "props": "MSBuild", "fth": "Forth", "forth": "Forth", "ftn": "FORTRAN Legacy", "asp": "ASP", "fish": "Fish", "irunargs": "Verilog Args File", "comp": "GLSL", "vala": "Vala", "wl": "Wolfram", "asax": "ASP.NET", "jl": "Julia", "erl": "Erlang", "asa": "ASP", "cxx": "C++", "org": "Org", "js": "JavaScript", "swift": "Swift", "bash": "BASH", "c": "C", "lucius": "Lucius", "xrunargs": "Verilog Args File", "upkg": "Unreal Script", "idr": "Idris", "nix": "Nix", "cogent": "Cogent", "cob": "COBOL", "toml": "TOML", "oz": "Oz", "srt": "SRecode Template", "pde": "Processing", "hbs": "Handlebars", "mad": "Madlang", "markdown": "Markdown", "cabal": "Cabal", "pgc": "C", "cc": "C++", "hex": "HEX", "sql": "SQL", "cpp": "C++", "vim": "Vim Script", "ipp": "C++ Header", "xtend": "Xtend", "e4": "Forth", "cs": "C#", "cr": "Crystal", "txt": "Plain Text", "uci": "Unreal Script", "ckt": "Spice Netlist", "java": "Java", "fsproj": "MSBuild", "hrl": "Erlang", "py": "Python", "makefile": "Makefile", "f83": "Forth", "gtpl": "Groovy", "json": "JSON", "master": "ASP.NET", "asm": "Assembly", "f03": "FORTRAN Modern", "frt": "Forth", "rhtml": "Ruby HTML", "f08": "FORTRAN Modern", "pl": "Perl", "pm": "Perl", "mm": "Objective C++", "hx": "Haxe", "ascx": "ASP.NET", "hs": "Haskell", "xaml": "XAML", "jai": "JAI", "pfo": "FORTRAN Legacy", "aspx": "ASP.NET", "hh": "C++ Header", "vue": "Vue", "fsx": "F#", "f95": "FORTRAN Modern", "rst": "ReStructuredText", "fsscript": "F#", "dtsi": "Device Tree", "groovy": "Groovy", "polly": "Polly", "thy": "Isabelle", "f": "FORTRAN Legacy", "julius": "Julius", "ada": "Ada", "mli": "OCaml", "mk": "Makefile", "kts": "Kotlin", "r": "R", "svg": "SVG", "dts": "Device Tree", "vhd": "VHDL", "qcl": "QCL", "zsh": "Zsh", "sty": "TeX", "def": "Module-Definition", "qml": "QML", "vb": "Visual Basic", "v": "Coq", "uc": "Unreal Script", "vg": "Verilog", "vh": "Verilog", "cassius": "Cassius", "pro": "Prolog", "cbl": "COBOL", "ccp": "COBOL", "inl": "C++ Header", "as": "ActionScript", "cobol": "COBOL", "svh": "SystemVerilog", "in": "Autoconf", "purs": "PureScript", "grt": "Groovy", "pas": "Pascal", "cmake": "CMake", "csh": "C Shell", "proto": "Protocol Buffers", "fpm": "Forth", "tese": "GLSL", "nb": "Wolfram", "hpp": "C++ Header", "s": "Assembly", "tesc": "GLSL", "html": "HTML", "pad": "Ada", "fst": "F*", "hxx": "C++ Header", "text": "Plain Text", "css": "CSS", "frag": "GLSL", "fr": "Forth", "fs": "F#", "ft": "Forth", "nim": "Nim", "urs": "Ur/Web", "y": "Happy", "fb": "Forth", "4th": "Forth", "hlean": "Lean", "cshtml": "Razor", "cljs": "ClojureScript", "mak": "Makefile", "php": "PHP", "jsx": "JSX", "agda": "Agda", "lidr": "Idris", "scss": "Sass", "e": "Specman e", "sass": "Sass", "ss": "Scheme", "fsi": "F#", "rs": "Rust", "m": "Objective C", "f90": "FORTRAN Modern", "cpy": "COBOL", "sh": "Shell", "urp": "Ur/Web Project", "exs": "Elixir", "kt": "Kotlin", "sc": "Scala", "f77": "FORTRAN Legacy", "mustache": "Mustache"}
var PathBlacklist = ""
var FilesOutput = ""
var DirFilePaths = []string{}

// Get all the files that exist in the directory
func walkDirectory(root string, output *chan *FileJob) {
	gitignore, gitignoreerror := gitignore.NewGitIgnore(filepath.Join(root, ".gitignore"))

	godirwalk.Walk(root, &godirwalk.Options{
		Unsorted: true,
		Callback: func(root string, info *godirwalk.Dirent) error {
			// TODO this should be configurable via command line
			if strings.HasPrefix(root, ".git/") || strings.HasPrefix(root, ".hg/") || strings.HasPrefix(root, ".svn/") {
				return filepath.SkipDir
			}

			if !info.IsDir() {
				if gitignoreerror != nil || !gitignore.Match(filepath.Join(root, info.Name()), false) {

					extension := getExtension(info.Name())
					language, ok := ExtensionToLanguage[extension]

					if ok {
						*output <- &FileJob{Location: root, Filename: info.Name(), Extension: extension, Language: language}
					}
				}
			}

			return nil
		},
		ErrorCallback: func(osPathname string, err error) godirwalk.ErrorAction {
			return godirwalk.SkipNode
		},
	})

	close(*output)
}

func Process() {
	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	// TODO these should be configurable by command line
	// TODO need to add an error channel and spit them out
	fileListQueue := make(chan *FileJob, runtime.NumCPU()*10000)          // Files ready to be read from disk
	fileReadJobQueue := make(chan *FileJob, runtime.NumCPU()*10)          // Workers reading from disk
	fileReadContentJobQueue := make(chan *FileJob, runtime.NumCPU()*5000) // Files ready to be processed
	fileProcessJobQueue := make(chan *FileJob, runtime.NumCPU())          // Workers doing the hard work
	fileSummaryJobQueue := make(chan *FileJob, runtime.NumCPU()*1000)     // Files ready to be summerised

	go walkDirectory(DirFilePaths[0], &fileListQueue)
	go fileBufferReader(&fileListQueue, &fileReadJobQueue)
	go fileReaderWorker(&fileReadJobQueue, &fileReadContentJobQueue)
	go fileBufferReader(&fileReadContentJobQueue, &fileProcessJobQueue)
	go fileProcessorWorker(&fileProcessJobQueue, &fileSummaryJobQueue)

	if FilesOutput == "" {
		fileSummerize(&fileSummaryJobQueue) // Bring it all back to you
	} else {
		fileSummerizeFiles(&fileSummaryJobQueue)
	}

	fmt.Println("NumCPU", runtime.NumCPU())
}
