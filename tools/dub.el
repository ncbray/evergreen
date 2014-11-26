(require 'generic-x)
    
(define-generic-mode 'dub-mode
  '("//") ;; comments
  '("func" "test" "struct" "implements" "star" "plus" "choose" "or" "question" "if" "else" "return" "var" "true" "false" "nil")
  '(
    ("\\[\\([^\]]\\)*\\]" . font-lock-constant-face) ;; TODO escaped brackets.
    ("\\+\\|\\*\\|/\\|\\-\\|\\$\\|!" . 'font-lock-builtin-face)
    ("\\b[0-9]+\\b" . 'font-lock-constant-face)
    ("\\<\\(string\\|int\\|uint32\\|int64\\|rune\\|bool\\|graph\\)\\>" . 'font-lock-type-face)
    ("{\\|}" . 'font-lock-builtin-face)
    ("=\\|:=?" . 'font-lock-builtin-face)
    (",\\|;" . 'font-lock-builtin-face)
    )
  '("\\.dub$") ;; filetype
   nil
  "A mode for dub files"
)
