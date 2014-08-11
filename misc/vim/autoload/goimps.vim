let s:save_cpo = &cpo
set cpo&vim

function! goimps#Importable()
  let s = system('goimps importable')
  if v:shell_error
    echoerr '[ERROR] goimps: errors occur on excuting `goimps importable`'
    return []
  endif

  return split(s, '\n')
endfunction

function! goimps#Dropable(filename)
  let s = system('goimps dropable ' . shellescape(a:filename))
  if v:shell_error
    echoerr '[ERROR] goimps: errors occur on excuting `goimps dropable`'  . shellescape(a:filename)
    return []
  endif

  return split(s, '\n')
endfunction

function! goimps#Unused(filename)
  let s = system('goimps unused ' . shellescape(a:filename))
  if v:shell_error
    echoerr '[ERROR] goimps: errors occur on excuting `goimps unused`'  . shellescape(a:filename)
    return []
  endif

  return split(s, '\n')
endfunction

let &cpo = s:save_cpo
unlet s:save_cpo
