@echo off
 
rem メッセージ表示
echo MsgBox "%1",vbInformation,"Alert" > %TEMP%\msgbox.vbs & %TEMP%\msgbox.vbs
 
rem ファイル削除
del /Q %TEMP%\msgbox.vbs
