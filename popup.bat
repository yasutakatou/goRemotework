@echo off
 
rem メッセージ表示
echo MsgBox "Alert",vbInformation,%1 > %TEMP%\msgbox.vbs & %TEMP%\msgbox.vbs
 
rem ファイル削除
del /Q %TEMP%\msgbox.vbs
