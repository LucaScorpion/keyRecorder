# keyRecorder

A tool that records keystrokes, and outputs them as a [keyScripter](https://github.com/LucaScorpion/keyScripter) script. It uses the Win32 [GetAsyncKeyState](https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getasynckeystate) API to read key states from the OS.

## Usage

Basic usage:

```
keyRecorder output.txt
```

By default, it will record keys until the escape key is pressed. To change this to a different key use the `-s` option with a [virtual key code](https://docs.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes). For example, to use the backspace key:

```
keyRecorder output.txt -s 8
```
