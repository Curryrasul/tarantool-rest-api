package main

func logger(logString string) {
	LogFile.WriteString(logString)
	LogFile.Sync()
}
